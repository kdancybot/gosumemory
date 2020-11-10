package pp

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
	"unsafe"

	"github.com/k0kubun/pp"
	"github.com/l3lackShark/gosumemory/db"
	"github.com/l3lackShark/gosumemory/memory"
	"github.com/spf13/cast"
)

//#cgo LDFLAGS: -lm
//#cgo CPPFLAGS: -DOPPAI_STATIC_HEADER
//#include <stdlib.h>
//#include "oppai.c"
import "C"

var ez C.ezpp_t

type PP struct {
	Total         C.float
	FC            C.float
	Strain        []float64
	StarRating    C.float
	AimStars      C.float
	SpeedStars    C.float
	AimPP         C.float
	SpeedPP       C.float
	Accuracy      C.float
	N300          C.int
	N100          C.int
	N50           C.int
	NMiss         C.int
	AR            C.float
	CS            C.float
	OD            C.float
	HP            C.float
	Artist        string
	ArtistUnicode string
	Title         string
	TitleUnicode  string
	Version       string
	Creator       string
	NCircles      C.int
	NSliders      C.int
	NSpinners     C.int
	ODMS          C.float
	Mode          C.int
	Combo         C.int
	MaxCombo      C.int
	Mods          C.int
	ScoreVersion  C.int
}

var strainArray []float64
var tempBeatmapFile string
var currMaxCombo C.int

func readData(data *PP, ez C.ezpp_t, needStrain bool, path string) error {

	if strings.HasSuffix(path, ".osu") {
		cpath := C.CString(path)

		defer C.free(unsafe.Pointer(cpath))
		if rc := C.ezpp(ez, cpath); rc < 0 {
			memory.MenuData.PP.PpStrains = []float64{0}
			return errors.New(C.GoString(C.errstr(rc)))
		}
		C.ezpp_set_base_ar(ez, C.float(memory.MenuData.Bm.Stats.MemoryAR))
		C.ezpp_set_base_od(ez, C.float(memory.MenuData.Bm.Stats.MemoryOD))
		C.ezpp_set_base_cs(ez, C.float(memory.MenuData.Bm.Stats.MemoryCS))
		C.ezpp_set_base_hp(ez, C.float(memory.MenuData.Bm.Stats.MemoryHP))
		C.ezpp_set_accuracy_percent(ez, C.float(memory.GameplayData.Accuracy))
		C.ezpp_set_mods(ez, C.int(memory.MenuData.Mods.AppliedMods))
		*data = PP{
			Artist:     C.GoString(C.ezpp_artist(ez)),
			Title:      C.GoString(C.ezpp_title(ez)),
			Version:    C.GoString(C.ezpp_version(ez)),
			Creator:    C.GoString(C.ezpp_creator(ez)),
			AR:         C.ezpp_ar(ez),
			CS:         C.ezpp_cs(ez),
			OD:         C.ezpp_od(ez),
			HP:         C.ezpp_hp(ez),
			StarRating: C.ezpp_stars(ez),
		}
		memory.MenuData.Bm.Stats.BeatmapSR = cast.ToFloat32(fmt.Sprintf("%.2f", float32(data.StarRating)))
		memory.MenuData.Bm.Stats.BeatmapAR = cast.ToFloat32(fmt.Sprintf("%.2f", float32(data.AR)))
		memory.MenuData.Bm.Stats.BeatmapCS = cast.ToFloat32(fmt.Sprintf("%.2f", float32(data.CS)))
		memory.MenuData.Bm.Stats.BeatmapOD = cast.ToFloat32(fmt.Sprintf("%.2f", float32(data.OD)))
		memory.MenuData.Bm.Stats.BeatmapHP = cast.ToFloat32(fmt.Sprintf("%.2f", float32(data.HP)))

		if needStrain == true {
			C.ezpp_set_end_time(ez, 0)
			C.ezpp_set_combo(ez, 0)
			C.ezpp_set_nmiss(ez, 0)
			memory.MenuData.Bm.Stats.BeatmapMaxCombo = int32(C.ezpp_max_combo(ez))
			memory.MenuData.Bm.Stats.FullSR = cast.ToFloat32(fmt.Sprintf("%.2f", float32(C.ezpp_stars(ez))))
			var bpmChanges []int
			for i := 0; i < int(C.ezpp_ntiming_points(ez)); i++ {
				msPerBeat := float64(C.ezpp_timing_ms_per_beat(ez, C.int(i)))
				timingChanges := int(C.ezpp_timing_change(ez, C.int(i)))
				if timingChanges == 1 {
					bpmFormula := int(math.Round(1 / msPerBeat * 1000 * 60 * 1)) //1 = bmpMultiplier
					if bpmFormula > 0 {
						bpmChanges = append(bpmChanges, bpmFormula)
					}
				}
			}
			memory.MenuData.Bm.Stats.BeatmapBPM.Minimal, memory.MenuData.Bm.Stats.BeatmapBPM.Maximal = minMax(bpmChanges)
			strainArray = nil
			seek := 0
			var window []float64
			var total []float64
			// for seek < int(C.ezpp_time_at(ez, C.ezpp_nobjects(ez)-1)) { //len-1
			for int32(seek) < memory.MenuData.Bm.Time.Mp3Time {
				for obj := 0; obj <= int(C.ezpp_nobjects(ez)-1); obj++ {
					if tempBeatmapFile != memory.MenuData.Bm.Path.BeatmapOsuFileString {
						return nil //Interrupt calcualtion if user has changed the map.
					}
					if int(C.ezpp_time_at(ez, C.int(obj))) >= seek && int(C.ezpp_time_at(ez, C.int(obj))) <= seek+3000 {
						window = append(window, float64(C.ezpp_strain_at(ez, C.int(obj), 0))+float64(C.ezpp_strain_at(ez, C.int(obj), 1)))
					}
				}
				sum := 0.0
				for _, num := range window {
					sum += num
				}
				total = append(total, sum/math.Max(float64(len(window)), 1))
				window = nil
				seek += 500
			}
			strainArray = total
			memory.MenuData.Bm.Time.FirstObj = int32(C.ezpp_time_at(ez, 0))
			memory.MenuData.Bm.Time.FullTime = int32(C.ezpp_time_at(ez, C.ezpp_nobjects(ez)-1))
		} else {
			C.ezpp_set_end_time(ez, C.float(memory.MenuData.Bm.Time.PlayTime))
			currMaxCombo = C.ezpp_max_combo(ez) //for RestSS
			C.ezpp_set_combo(ez, C.int(memory.GameplayData.Combo.Max))
			C.ezpp_set_nmiss(ez, C.int(memory.GameplayData.Hits.H0))
		}

		*data = PP{
			Total:  C.ezpp_pp(ez),
			Strain: strainArray,

			AimStars:   C.ezpp_aim_stars(ez),
			SpeedStars: C.ezpp_speed_stars(ez),
			AimPP:      C.ezpp_aim_pp(ez),
			SpeedPP:    C.ezpp_speed_pp(ez),
			Accuracy:   C.ezpp_accuracy_percent(ez),
			N300:       C.ezpp_n300(ez),
			N100:       C.ezpp_n100(ez),
			N50:        C.ezpp_n50(ez),
			NMiss:      C.ezpp_nmiss(ez),
			//ArtistUnicode: C.GoString(C.ezpp_artist_unicode(ez)),
			//	TitleUnicode:  C.GoString(C.ezpp_title_unicode(ez)),
			NCircles:     C.ezpp_ncircles(ez),
			NSliders:     C.ezpp_nsliders(ez),
			NSpinners:    C.ezpp_nspinners(ez),
			ODMS:         C.ezpp_odms(ez),
			Mode:         C.ezpp_mode(ez),
			Combo:        C.ezpp_combo(ez),
			MaxCombo:     C.ezpp_max_combo(ez),
			Mods:         C.ezpp_mods(ez),
			ScoreVersion: C.ezpp_score_version(ez),
		}
		memory.MenuData.PP.PpStrains = data.Strain
	}
	return nil
}

var maniaSR float64
var maniaMods int32
var maniaHitObjects float64
var tempMods string

func GetData() {

	ez = C.ezpp_new()
	C.ezpp_set_autocalc(ez, 1)
	//defer C.ezpp_free(ez)

	for {

		if memory.DynamicAddresses.IsReady == true {
			switch memory.MenuData.GameMode {
			case 0, 1:
				var data PP
				if tempBeatmapFile != memory.MenuData.Bm.Path.BeatmapOsuFileString || memory.MenuData.Mods.PpMods != tempMods { //On map/mods change
					tempBadJudgments = 0
					path := memory.MenuData.Bm.Path.FullDotOsu
					tempBeatmapFile = memory.MenuData.Bm.Path.BeatmapOsuFileString
					tempMods = memory.MenuData.Mods.PpMods
					mp3Time, err := calculateMP3Time()
					if err == nil {
						memory.MenuData.Bm.Time.Mp3Time = mp3Time
					}
					//Get Strains
					readData(&data, ez, true, path)

					//pp.Println(memory.MenuData.Bm.Metadata)
				}

				switch memory.MenuData.OsuStatus {
				case 2, 7:
					path := memory.MenuData.Bm.Path.FullDotOsu
					readData(&data, ez, false, path)
					if memory.GameplayData.Combo.Max > 1 && float64(data.Total) > 0 {
						memory.GameplayData.PP.Pp = cast.ToInt32(float64(data.Total))
					}
				case 5:
					memory.GameplayData.PP.Pp = 0
				}

			case 3:

				if tempBeatmapFile != memory.MenuData.Bm.Path.BeatmapOsuFileString || memory.MenuData.Mods.PpMods != tempMods { //On map/mods change
					memory.MenuData.Bm.Time.FullTime = 0        //Not implemented for mania yet
					memory.MenuData.Bm.Stats.BeatmapAR = 0      //Not implemented for mania yet
					memory.MenuData.Bm.Stats.BeatmapCS = 0      //Not implemented for mania yet
					memory.MenuData.Bm.Stats.BeatmapOD = 0      //Not implemented for mania yet
					memory.MenuData.Bm.Stats.BeatmapHP = 0      //Not implemented for mania yet
					memory.MenuData.PP.PpStrains = []float64{0} //Not implemented for mania yet

					tempBeatmapFile = memory.MenuData.Bm.Path.BeatmapOsuFileString
					tempMods = memory.MenuData.Mods.PpMods
					maniaSR = 0.0
					maniaMods = 0
					maniaHitObjects = 0.0
					for i := 0; i < len(db.OsuDB.BmInfo); i++ {
						if tempBeatmapFile == db.OsuDB.BmInfo[i].Filename {
							if strings.Contains(memory.MenuData.Mods.PpMods, "DT") {
								maniaMods = 64
							} else if strings.Contains(memory.MenuData.Mods.PpMods, "HT") {
								maniaMods = 256
							} else {
								maniaMods = 0 //assuming NM
							}
							for j := 0; j < len(db.OsuDB.BmInfo[i].StarRatingMania); j++ {
								if maniaMods == db.OsuDB.BmInfo[i].StarRatingMania[j].BitMods {
									maniaSR = db.OsuDB.BmInfo[i].StarRatingMania[j].StarRating
									maniaHitObjects = float64(db.OsuDB.BmInfo[i].NumHitCircles) + float64(db.OsuDB.BmInfo[i].NumSliders) + float64(db.OsuDB.BmInfo[i].NumSpinners)
									memory.MenuData.Bm.Stats.BeatmapSR = cast.ToFloat32(fmt.Sprintf("%.2f", float32(maniaSR)))
									// memory.MenuData.Bm.Metadata.Artist = db.OsuDB.BmInfo[i].Artist
									// memory.MenuData.Bm.Metadata.Title = db.OsuDB.BmInfo[i].Title
									// memory.MenuData.Bm.Metadata.Mapper = db.OsuDB.BmInfo[i].Creator
									// memory.MenuData.Bm.Metadata.Version = db.OsuDB.BmInfo[i].Difficulty //Now sets through memory.
									memory.GameplayData.PP.PPifFC = int32(calculateManiaPP(float64(memory.MenuData.Bm.Stats.MemoryOD), maniaSR, maniaHitObjects, 1000000.0)) //PP if SS
									break
								}
							}
							if maniaSR == 0.0 {
								pp.Println("Could not find mania star rating in the database. PP output will be unavailable for this beatmap!")
							}
							break
						}
					}
				}
				if maniaSR >= 0.01 {
					if memory.GameplayData.Score >= 500000 {
						memory.GameplayData.PP.Pp = int32(calculateManiaPP(float64(memory.MenuData.Bm.Stats.MemoryOD), maniaSR, maniaHitObjects, float64(memory.GameplayData.Score)))
					} else {
						memory.GameplayData.PP.Pp = 0
					}

				}

			}

		}

		time.Sleep(time.Duration(memory.UpdateTime) * time.Millisecond)
	}
}
