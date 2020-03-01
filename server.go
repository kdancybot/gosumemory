package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/Andoryuuta/kiwi"
	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
)

var proc, procerr = kiwi.GetProcessByFileName("osu!.exe")

var osuBase uintptr
var osuStatus uint16
var currentBeatmapData uintptr
var playContainer uintptr
var playContainerBase uintptr
var serverBeatmapString string
var outStrLoop string
var baseDir string
var playTimeBase uintptr
var playTime uintptr
var currentBeatmapDataBase uint32
var currentBeatmapDataFirtLevel uint32
var playContainerBaseAddr uint32
var playContainerFirstlevel uint32
var playContainer38 uint32
var fullPathToOsu string
var osuFileStdIN string
var ourTime []int
var lastObjectInt int
var lastObject string
var ppAcc string
var ppCombo string
var pp100 string
var pp50 string
var ppMiss string
var ppMods string = ""
var pp string = ""
var ppSS string = ""
var pp99 string = ""
var pp98 string = ""
var pp97 string = ""
var pp96 string = ""
var pp95 string = ""
var ppifFC string = ""
var innerBGPath string = ""
var updateTime int
var isRunning = 0
var workingDirectory string
var operatingSystem int8
var uintptrOsuStatus uintptr
var jsonByte []byte
var isPizdec int8 = 0
var counter string

func Cmd(cmd string, shell bool) []byte {

	if shell {
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			println("some error found", err)
		}
		return out
	} else {
		out, err := exec.Command(cmd).Output()
		if err != nil {
			println("some error found2", err)
		}
		return out

	}
}

func OsuStatusAddr() uintptr { //in hopes to deprecate this
	if operatingSystem == 1 {
		cmd, err := exec.Command("OsuStatusAddr.exe").Output()
		if err != nil {
			fmt.Println(err)
		}
		outStr := cast.ToString(cmd)
		outStr = strings.Replace(outStr, "\n", "", -1)
		outStr = strings.Replace(outStr, "\r", "", -1)
		outInt := cast.ToUint32(outStr)

		osuBase = uintptr(outInt)

	} else {
		x := Cmd("scanmem -p `pgrep osu\\!.exe` -e -c 'option scan_data_type bytearray;48 83 F8 04 73 1E;list;exit'", true)
		outStr := cast.ToString(x)
		outStr = strings.Replace(outStr, " ", "", -1)

		input := outStr
		if input == "" {
			log.Fatalln("osu! is probably not fully loaded, please load the game up and try again!")
		}
		output := (input[3:])
		yosuBase := firstN(output, 8)
		check := strings.Contains(yosuBase, ",")
		if check == true {
			yosuBase = strings.Replace(yosuBase, ",", "", -1)
		}
		osuBaseString := "0x" + yosuBase
		osuBaseUINT32 := cast.ToUint32(osuBaseString)
		osuBase = uintptr(osuBaseUINT32)
	}
	if osuBase == 0 {
		log.Fatalln("could not find osuStatusAddr, is osu! running?")
	}

	//println(CurrentBeatmapFolderString())
	return osuBase

}
func OsuBaseAddr() uintptr { //in hopes to deprecate this
	if operatingSystem == 1 {
		cmd, err := exec.Command("OsuBaseAddr.exe").Output()
		if err != nil {
			fmt.Println(err)
		}
		outStr := cast.ToString(cmd)
		outStr = strings.Replace(outStr, "\n", "", -1)
		outStr = strings.Replace(outStr, "\r", "", -1)
		outInt := cast.ToUint32(outStr)
		osuBase = uintptr(outInt)
	} else {
		x := Cmd("scanmem -p `pgrep osu\\!.exe` -e -c 'option scan_data_type bytearray;F8 01 74 04 83;list;exit'", true)
		outStr := cast.ToString(x)
		outStr = strings.Replace(outStr, " ", "", -1)

		input := outStr
		if input == "" {
			log.Fatalln("OsuBase addr fail")
		}
		output := (input[3:])
		yosuBase := firstN(output, 8)
		check := strings.Contains(yosuBase, ",")
		if check == true {
			yosuBase = strings.Replace(yosuBase, ",", "", -1)
		}
		osuBaseString := "0x" + yosuBase
		osuBaseUINT32 := cast.ToUint32(osuBaseString)
		osuBase = uintptr(osuBaseUINT32)
	}

	if osuBase == 0 {
		log.Fatalln("Could not find OsuBaseAddr, is osu! running?")
	}

	//println(CurrentBeatmapFolderString())
	return osuBase
}

func OsuPlayTimeAddr() uintptr { //in hopes to deprecate this
	if operatingSystem == 1 {
		cmd, err := exec.Command("OsuPlayTimeAddr.exe").Output()
		if err != nil {
			fmt.Println(err)
		}
		outStr := cast.ToString(cmd)
		outStr = strings.Replace(outStr, "\n", "", -1)
		outStr = strings.Replace(outStr, "\r", "", -1)
		outInt := cast.ToUint32(outStr)

		osuBase = uintptr(outInt)
	} else {
		x := Cmd("scanmem -p `pgrep osu\\!.exe` -e -c 'option scan_data_type bytearray;5E 5F 5D C3 A1 ?? ?? ?? ?? 89 ?? 04;list;exit'", true)
		outStr := cast.ToString(x)
		outStr = strings.Replace(outStr, " ", "", -1)

		input := outStr
		if input == "" {
			log.Fatalln("OsuBase addr fail")
		}
		output := (input[3:])
		yosuBase := firstN(output, 8)
		check := strings.Contains(yosuBase, ",")
		if check == true {
			yosuBase = strings.Replace(yosuBase, ",", "", -1)
		}
		osuBaseString := "0x" + yosuBase
		osuBaseUINT32 := cast.ToUint32(osuBaseString)
		osuBase = uintptr(osuBaseUINT32)
	}
	if osuBase == 0 {
		log.Fatalln("OsuPlayTimeAddr is not found")
	}

	//println(CurrentBeatmapFolderString())
	return osuBase
}

func OsuplayContainer() uintptr { //in hopes to deprecate this
	if operatingSystem == 1 {

		cmd, err := exec.Command("OsuPlayContainer.exe").Output()
		if err != nil {
			fmt.Println(err)
		}
		outStr := cast.ToString(cmd)
		outStr = strings.Replace(outStr, "\n", "", -1)
		outStr = strings.Replace(outStr, "\r", "", -1)
		outInt := cast.ToUint32(outStr)

		osuBase = uintptr(outInt)
	} else {
		x := Cmd("scanmem -p `pgrep osu\\!.exe` -e -c 'option scan_data_type bytearray;85 C9 74 1F 8D 55 F0 8B 01;list;exit'", true)
		outStr := cast.ToString(x)
		outStr = strings.Replace(outStr, " ", "", -1)

		input := outStr
		if input == "" {
			log.Fatalln("osuplayContainer addr fail")
		}
		output := (input[3:])
		yosuBase := firstN(output, 8)
		check := strings.Contains(yosuBase, ",")
		if check == true {
			yosuBase = strings.Replace(yosuBase, ",", "", -1)
		}
		osuBaseString := "0x" + yosuBase

		osuBaseUINT32 := cast.ToUint32(osuBaseString)
		osuBase = uintptr(osuBaseUINT32)
	}
	if osuBase == 0 {
		log.Fatalln("Could not find osuplayContainer address, is osu! running?")
	}

	//println(CurrentBeatmapFolderString())
	return osuBase
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
	for procerr != nil { //TODO: refactor
		ws.WriteMessage(1, []byte("osu!.exe not found"))
		if operatingSystem == 1 {
			log.Println("is osu! running? (osu! process was not found, waiting...)")
		} else {
			log.Println("is osu! running? (We don't support client restarts on linux, assuming that we just lost the process for a second, retrying... (if you closed the game, pleae restart the program.))")
		}

		proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
	}
	if isRunning == 0 {
		fmt.Println("Client Connected, please go to the SongSelect and check this console back.")
		//	time.Sleep(5 * time.Second)            //hack to wait for the game
		StaticOsuStatusAddr := OsuStatusAddr() //we should only check for this address once.
		osuStatusOffset, err := proc.ReadUint32(StaticOsuStatusAddr - 0x4)
		if err != nil {
			ws.WriteMessage(1, []byte("osu!status offset was not found"))
			log.Fatalln("osu!status offset was not found, are you sure that osu!stable is running? If so, please report this to GitHub!")
		}
		uintptrOsuStatus = uintptr(osuStatusOffset)
		osuStatusValue, err := proc.ReadUint16(uintptrOsuStatus)
		if err != nil {
			ws.WriteMessage(1, []byte("osu!status value was not found"))
			log.Fatalln("osu!status value was not found, are you sure that osu!stable is running? If so, please report this to GitHub!")
		}

		for osuStatusValue != 5 {
			log.Println("please go to songselect in order to proceed!")
			osuStatusValue, err = proc.ReadUint16(uintptrOsuStatus)
			if err != nil {
				log.Fatalln("is osu! running? (osu! status was not found)")
			}
			ws.WriteMessage(1, []byte("osu! is not in SongSelect!"))

			time.Sleep(500 * time.Millisecond)

		}

		//time.Sleep(5 * time.Second)
		osuBase = OsuBaseAddr()
		currentBeatmapData = (osuBase - 0xC)
		playTimeBase = OsuPlayTimeAddr()
		playContainer = OsuplayContainer()
		playContainerBase = (playContainer - 0x4)
		playTime = (playTimeBase + 0x5)

		if CurrentPlayTime() == -1 {
			fmt.Println("Failed to get the correct offsets, retrying...")
			restart()

		}

	}
	isRunning = 1
	fmt.Println("it seems that we got the correct offsets, you are good to go!")
	log.Println("Client Connected")

	var tempCurrentBeatmapOsu string

	for {

		osuStatusValue, err := proc.ReadUint16(uintptrOsuStatus)
		osuStatus = osuStatusValue
		var proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
		for procerr != nil {
			isPizdec = 1
			fmt.Println("isPizdec = 1")
			log.Println("is osu! running? (osu! process was not found, waiting...)")
			proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
			time.Sleep(1 * time.Second)
		}
		if isPizdec == 1 {
			if operatingSystem == 1 {
				fmt.Println("It looks like we have a client restart!")
				isPizdec = 0
				fmt.Println("isPizdec = 0")
				time.Sleep(10 * time.Second) // hack to wait for a client restart
				restart()
			} else {
				fmt.Println("We don't support client restart on linux yet!")
				isPizdec = 0 // Assuming that it was just a matter of losing the process
			}

		}

		//init base stuff
		currentBeatmapDataBase, err = proc.ReadUint32(currentBeatmapData)
		if err != nil {
			log.Println("currentBeatmapDataBase Base level failure")
		}
		currentBeatmapDataFirtLevel, err = proc.ReadUint32(uintptr(currentBeatmapDataBase))
		if err != nil {
			log.Println("currentBeatmapDataFirtLevel First level pointer failure")
		}
		if osuStatus == 2 {
			playContainerBaseAddr, err = proc.ReadUint32(playContainerBase)
			if err != nil {
				log.Println("playContainerBaseAddr Base level failure")
			}
			playContainerFirstlevel, err = proc.ReadUint32(uintptr(playContainerBaseAddr))
			if err != nil {
				log.Println("playContainerFirstlevel pointer failure")
			}
			playContainer38, err = proc.ReadUint32(uintptr(playContainerFirstlevel) + 0x38)
			if err != nil {
				//	log.Println("playContainer38 pointer failure")

			}
		}

		type PlayContainer struct {
			CurrentHit300c          int16   `json:"300"`
			CurrentHit100c          int16   `json:"100"`
			CurrentHit50c           int16   `json:"50"`
			CurrentHitMiss          int16   `json:"miss"`
			CurrentAccuracy         float64 `json:"accuracy"`
			CurrentScore            int32   `json:"score"`
			CurrentCombo            int32   `json:"combo"`
			CurrentGameMode         int32   `json:"gameMode"`
			CurrentAppliedMods      int32   `json:"appliedMods"`
			PpMods                  string  `json:"appliedModsString"`
			CurrentMaxCombo         int32   `json:"maxCombo"`
			CurrentPlayerHP         int8    `json:"playerHP"`
			CurrentPlayerHPSmoothed int8    `json:"playerHPSmoothed"`
			Pp                      string  `json:"pp"`
			PPifFC                  string  `json:"ppIfFC"`
		}
		type EverythingInMenu struct {
			CurrentState                uint16  `json:"osuState"`
			CurrentBeatmapID            uint32  `json:"bmID"`
			CurrentBeatmapSetID         uint32  `json:"bmSetID"`
			CurrentBeatmapCS            float32 `json:"CS"`
			CurrentBeatmapAR            float32 `json:"AR"`
			CurrentBeatmapOD            float32 `json:"OD"`
			CurrentBeatmapHP            float32 `json:"HP"`
			CurrentBeatmapString        string  `json:"bmInfo"`
			CurrentBeatmapFolderString  string  `json:"bmFolder"`
			CurrentBeatmapOsuFileString string  `json:"pathToBM"`
			CurrentPlayTime             int32   `json:"bmCurrentTime"`
			InnerBGPath                 string  `json:"innerBG"`
			PpSS                        string  `json:"ppSS"`
			Pp99                        string  `json:"pp99"`
			Pp98                        string  `json:"pp98"`
			Pp97                        string  `json:"pp97"`
			Pp96                        string  `json:"pp96"`
			Pp95                        string  `json:"pp95"`
		}

		type EverythingInMenu2 struct { //order sets here
			D EverythingInMenu `json:"menuContainer"`
			P PlayContainer    `json:"gameplayContainer"`
		}

		PlayContainerStruct := PlayContainer{
			CurrentHit300c:          CurrentHit300c(),
			CurrentHit100c:          CurrentHit100c(),
			CurrentHit50c:           CurrentHit50c(),
			CurrentHitMiss:          CurrentHitMiss(),
			CurrentScore:            CurrentScore(),
			CurrentAccuracy:         CurrentAccuracy(),
			CurrentCombo:            CurrentCombo(),
			CurrentGameMode:         CurrentGameMode(),
			CurrentAppliedMods:      CurrentAppliedMods(),
			CurrentMaxCombo:         CurrentMaxCombo(),
			CurrentPlayerHP:         CurrentPlayerHP(),
			CurrentPlayerHPSmoothed: CurrentPlayerHPSmoothed(),
			Pp:                      pp,
			PPifFC:                  ppifFC,
			PpMods:                  ppMods,
		}

		//println(ValidCurrentBeatmapFolderString())
		// if strings.HasSuffix(CurrentBeatmapOsuFileString(), ".osu") == false {
		// 	println(".osu ends with ???")
		// }
		// if strings.HasSuffix(CurrentBeatmapString(), "]") == false {
		// 	println("beatmapstring ends with ???")
		// }
		MenuContainerStruct := EverythingInMenu{CurrentState: osuStatus,
			CurrentBeatmapID:            CurrentBeatmapID(),
			CurrentBeatmapSetID:         CurrentBeatmapSetID(),
			CurrentBeatmapString:        CurrentBeatmapString(),
			CurrentBeatmapFolderString:  CurrentBeatmapFolderString(),
			CurrentBeatmapOsuFileString: CurrentBeatmapOsuFileString(),
			CurrentBeatmapAR:            CurrentBeatmapAR(),
			CurrentBeatmapOD:            CurrentBeatmapOD(),
			CurrentBeatmapCS:            CurrentBeatmapCS(),
			CurrentBeatmapHP:            CurrentBeatmapHP(),
			CurrentPlayTime:             CurrentPlayTime(),
			InnerBGPath:                 innerBGPath,
			PpSS:                        ppSS,
			Pp99:                        pp99,
			Pp98:                        pp98,
			Pp97:                        pp97,
			Pp96:                        pp96,
			Pp95:                        pp95,
		}

		for _, hitObjectTime := range ourTime {

			if int32(hitObjectTime) >= MenuContainerStruct.CurrentPlayTime { //TODO: Fix inaccuracy

				lastObjectInt = SliceIndex(len(ourTime), func(i int) bool { return ourTime[i] == hitObjectTime })
				lastObject = cast.ToString(lastObjectInt)
				ppAcc = cast.ToString(PlayContainerStruct.CurrentAccuracy)
				pp100 = cast.ToString(PlayContainerStruct.CurrentHit100c)
				pp50 = cast.ToString(PlayContainerStruct.CurrentHit50c)
				ppCombo = cast.ToString(PlayContainerStruct.CurrentMaxCombo)
				ppMiss = cast.ToString(PlayContainerStruct.CurrentHitMiss)
				ppMods = ModsResolver(cast.ToUint32(PlayContainerStruct.CurrentAppliedMods)) //TODO: Should only be called once)
				pp = PP()                                                                    //current pp
				ppifFC = PPifFC()

				break

			}
		}

		// if osuStatus == 5 || osuStatus == 4 { //TODO: Refactor
		// 	ppSS = PPSS()
		// 	pp99 = PP99()
		// 	pp98 = PP98()
		// 	pp97 = PP97()
		// 	pp96 = PP96()
		// 	pp95 = PP95()
		// }

		if MenuContainerStruct.CurrentBeatmapOsuFileString != tempCurrentBeatmapOsu {
			ourTime = nil
			pp = ""
			ppifFC = ""

			tempCurrentBeatmapOsu = MenuContainerStruct.CurrentBeatmapOsuFileString
			fullPathToOsu = fmt.Sprintf(baseDir + "/" + MenuContainerStruct.CurrentBeatmapFolderString + "/" + MenuContainerStruct.CurrentBeatmapOsuFileString)
			ppSS = PPSS()
			pp99 = PP99()
			pp98 = PP98()
			pp97 = PP97()
			pp96 = PP96()
			pp95 = PP95()

			j, err := ioutil.ReadFile(fullPathToOsu) // possibe file open exc
			if err != nil {
				//fmt.Println("osu file was not found2")
			}
			osuFileStdIN = string(j)
			if strings.Contains(osuFileStdIN, "[HitObjects]") == true {
				splitted := strings.Split(osuFileStdIN, "[HitObjects]")[1]
				newline := strings.Split(splitted, "\n")

				for i := 0; i < len(newline); i++ { //TODO: Add proper exception handler
					if len(newline[i]) > 1 {
						elements := strings.Split(newline[i], ",")[2]
						elementsInt := cast.ToInt(elements)
						ourTime = append(ourTime, elementsInt)
					}
				}
			}

			if strings.HasSuffix(fullPathToOsu, ".osu") == true {
				//fmt.Println(fullPathToOsu)
				file, err := os.Open(fullPathToOsu)
				if err != nil {
					log.Println(err, "in error")
					defer file.Close()
				}
				defer file.Close()
				scanner := bufio.NewScanner(file)
				var bgString string
				for scanner.Scan() {

					//fmt.Println(scanner.Text())
					if strings.Contains(scanner.Text(), ".jpg") == true {
						bg := strings.Split(scanner.Text(), "\"")
						bgString = (bg[1])
						break
						//log.Fatalln(scanner.Text())
					}
					if strings.Contains(scanner.Text(), ".png") == true {
						bg := strings.Split(scanner.Text(), "\"")
						bgString = (bg[1])
						//log.Fatalln(scanner.Text())
						break
					}
					if strings.Contains(scanner.Text(), ".JPG") == true {
						bg := strings.Split(scanner.Text(), "\"")
						bgString = (bg[1])
						//log.Fatalln(scanner.Text())
						break
					}
					if strings.Contains(scanner.Text(), ".PNG") == true {
						bg := strings.Split(scanner.Text(), "\"")
						bgString = (bg[1])
						break
						//log.Fatalln(scanner.Text())
					} else {
						bgString = ""
					}
				}
				if err := scanner.Err(); err != nil {
					log.Println(err)
				}
				if bgString != "" {
					//var fullPathToBG string = fmt.Sprintf(baseDir + "/" + MenuContainerStruct.CurrentBeatmapFolderString + "/" + bgString)
					innerBGPath = MenuContainerStruct.CurrentBeatmapFolderString + "/" + bgString
					//var fullBGCommand string = fmt.Sprintf("ln -nsf " + "\"" + fullPathToBG + "\"" + " " + "$PWD" + "/bg.png")
					//Cmd((fullBGCommand), true)
				}

				//fmt.Println(OsuHitobjects())
			} else {
				//	fmt.Println("osu file was not found")
			}

		}
		group := EverythingInMenu2{
			P: PlayContainerStruct,
			D: MenuContainerStruct,
		}
		jsonByte, err = json.Marshal(group)
		if err != nil {
			fmt.Println("error:", err)
		}
		ws.WriteMessage(1, []byte(jsonByte)) //sending data to the client

		//if err != nil {
		//	log.Println(err)
		//}
		time.Sleep(time.Duration(updateTime) * time.Millisecond)

	}

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {

	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Enter text: ")
	// counter, _ = reader.ReadString('\n')

	if runtime.GOOS == "windows" {
		fmt.Println("Hello from Windows, Please add a browser source in obs to http://127.0.0.1:24050 or refresh the page if you already did that.")
		operatingSystem = 1
	}
	if runtime.GOOS == "linux" {
		fmt.Println("Hello from Linux, Please add a browser source in obs to http://127.0.0.1:24050 or refresh the page if you already did that.")
		operatingSystem = 2

	}
	if operatingSystem == 2 { // hack to fix "Too many open files"
		var rLimit syscall.Rlimit
		err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err != nil {
			fmt.Println("Error Getting Rlimit ", err)
		}
		fmt.Println(rLimit)
		rLimit.Max = 999999
		rLimit.Cur = 999999
		err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err != nil {
			fmt.Println("Error Setting Rlimit ", err)
		}
		err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err != nil {
			fmt.Println("Error Getting Rlimit ", err)
		}
		fmt.Println("Rlimit Final", rLimit)
	}

	path := flag.String("path", "null", "Path to osu! Songs directory ex: C:\\Users\\BlackShark\\AppData\\Local\\osu!\\Songs")
	updateTimeAs := flag.Int("update", 100, "How fast should we update the values? (in milliseconds)")
	flag.Parse()
	updateTime = *updateTimeAs
	workingDirectory = *path

	if workingDirectory == "null" {
		log.Fatalln("Please set up your osu! Songs directory. (see --help)")
	}
	baseDir = workingDirectory
	go HTTPServer()
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8085", nil))
}
func restart() {
	isRunning = 0
	proc, procerr = kiwi.GetProcessByFileName("osu!.exe")
	if procerr != nil { //TODO: refactor
		log.Fatalln("is osu! running? (osu! process was not found)")
	}
	if isRunning == 0 {
		fmt.Println("Client Connected, please go to song select and check this console back.")
		StaticOsuStatusAddr := OsuStatusAddr() //we should only check for this address once.
		osuStatusOffset, err := proc.ReadUint32(StaticOsuStatusAddr - 0x4)
		if err != nil {
			log.Fatalln("osu!status offset was not found, are you sure that osu!stable is running? If so, please report this to GitHub!")
		}
		uintptrOsuStatus = uintptr(osuStatusOffset)
		osuStatusValue, err := proc.ReadUint16(uintptrOsuStatus)
		if err != nil {
			log.Fatalln("osu!status value was not found, are you sure that osu!stable is running? If so, please report this to GitHub!")
		}

		for osuStatusValue != 5 {
			log.Println("please go to songselect in order to proceed!")
			osuStatusValue, err = proc.ReadUint16(uintptrOsuStatus)
			if err != nil {
				log.Fatalln("is osu! running? (osu! status was not found)")
			}

			time.Sleep(500 * time.Millisecond)

			time.Sleep(1 * time.Second)

		}
		osuBase = OsuBaseAddr()
		currentBeatmapData = (osuBase - 0xC)
		playTimeBase = OsuPlayTimeAddr()
		playContainer = OsuplayContainer()
		playContainerBase = (playContainer - 0x4)
		playTime = (playTimeBase + 0x5)
		isRunning = 1
		if CurrentPlayTime() == -1 {
			fmt.Println("Failed to get the correct offsets, retrying...")
			restart()

		}

		fmt.Println("it seems that we got the correct offsets, you are good to go!")
	}

}

func firstN(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

func CurrentBeatmapID() uint32 { //currentbeatmapdata

	currentBeatmapID, err := proc.ReadUint32(uintptr(currentBeatmapDataFirtLevel + 0xC4))
	if err != nil {
		//log.Println("CurrentBeatmapID result pointer failure")
		return 0
	}
	return currentBeatmapID
}
func CurrentBeatmapSetID() uint32 {
	currentSetBeatmapID, err := proc.ReadUint32(uintptr(currentBeatmapDataFirtLevel + 0xC8))
	if err != nil {
		//log.Println("CurrentBeatmapSetID result pointer failure")
		return 0
	}
	return currentSetBeatmapID
}
func CurrentBeatmapString() string {
	beatmapStringSecondLevel, err := proc.ReadUint32(uintptr(currentBeatmapDataFirtLevel + 0x7C))
	if err != nil {
		//log.Println("BeatMapString Second level pointer failure")
		return "-4"
	}

	beatmapStringResult, err := proc.ReadNullTerminatedUTF16String((uintptr(beatmapStringSecondLevel + 0x8)))
	if err != nil {
		//	log.Println("BeatMapString Third level pointer failure")
		return "-5"
	}
	return beatmapStringResult
}

func CurrentBeatmapFolderString() string {

	beatmapFolderStringSecondLevel, err := proc.ReadUint32(uintptr(currentBeatmapDataFirtLevel + 0x74))
	if err != nil {
		//	log.Println("BeatMapFolderString Second level pointer failure")
		return "-4"
	}
	beatmapStringResult, err := proc.ReadNullTerminatedUTF16String((uintptr(beatmapFolderStringSecondLevel + 0x8)))
	if err != nil {
		//	log.Println("BeatMapFolderString Third level pointer failure")
		return "-5"
	}
	return beatmapStringResult
}
func CurrentBeatmapOsuFileString() string {

	beatmapFolderStringSecondLevel, err := proc.ReadUint32(uintptr(currentBeatmapDataFirtLevel + 0x8C))
	if err != nil {
		//log.Println("BeatMapOsuFileString Second level pointer failure")
		return "-4"
	}
	beatmapStringResult, err := proc.ReadNullTerminatedUTF16String((uintptr(beatmapFolderStringSecondLevel + 0x8)))
	if err != nil {
		//	log.Println("BeatMapString Third level pointer failure")
		return "-5"
	}
	// beatmapString := string(beatmapStringResult)
	// beatmapValidString := strings.ToValidUTF8(beatmapString, "")
	return beatmapStringResult
}
func CurrentBeatmapAR() float32 {
	currentSetBeatmapID, err := proc.ReadFloat32(uintptr(currentBeatmapDataFirtLevel + 0x2C))
	if err != nil {
		//	log.Println("AR result level pointer failure")
		return -5
	}
	return currentSetBeatmapID
}
func CurrentBeatmapCS() float32 {
	currentSetBeatmapID, err := proc.ReadFloat32(uintptr(currentBeatmapDataFirtLevel + 0x30))
	if err != nil {
		//	log.Println("CS result level pointer failure")
		return -4
	}
	return currentSetBeatmapID
}
func CurrentBeatmapHP() float32 {
	currentSetBeatmapID, err := proc.ReadFloat32(uintptr(currentBeatmapDataFirtLevel + 0x34))
	if err != nil {
		//	log.Println("HP result level pointer failure")
		return -5
	}
	return currentSetBeatmapID
}
func CurrentBeatmapOD() float32 {
	currentSetBeatmapID, err := proc.ReadFloat32(uintptr(currentBeatmapDataFirtLevel + 0x38))
	if err != nil {
		//	log.Println("OD result level pointer failure")
		return -5
	}
	return currentSetBeatmapID
}

// ------------------- PlayContainer
func CurrentAppliedMods() int32 {
	if osuStatus != 2 {
		return -1
	}
	currentCombo, err := proc.ReadInt32(uintptr(playContainer38 + 0x1C))
	if err != nil {
		//		log.Println("CurrentCombo result pointer failure")
		return -5
	}
	xorVal1, err := proc.ReadInt32(uintptr(currentCombo + 0xC))
	if err != nil {
		//	log.Println("CurrentCombo result pointer failure")
		return -6
	}
	xorVal2, err := proc.ReadInt32(uintptr(currentCombo + 0x8))
	if err != nil {
		//	log.Println("CurrentCombo result pointer failure")
		return -7
	}
	val := xorVal2 ^ xorVal1
	return val
}
func CurrentCombo() int32 {
	if osuStatus != 2 {
		return -1
	}

	currentCombo, err := proc.ReadInt32(uintptr(playContainer38 + 0x90))
	if err != nil {
		//	log.Println("CurrentCombo result pointer failure")
		return -5
	}
	return currentCombo
}

func CurrentMaxCombo() int32 {
	if osuStatus != 2 {
		return -1
	}

	currentCombo, err := proc.ReadInt32(uintptr(playContainer38 + 0x68))
	if err != nil {
		//		log.Println("CurrentCombo result pointer failure")
		return -5
	}
	return currentCombo
}
func CurrentHit100c() int16 {
	if osuStatus != 2 {
		return -1
	}

	currentCombo, err := proc.ReadInt16(uintptr(playContainer38 + 0x84)) //2 bytes
	if err != nil {
		//	log.Println("CurrentHit100c result pointer failure")
		return -5
	}
	return currentCombo
}
func CurrentHit300c() int16 {
	if osuStatus != 2 {
		return -1
	}
	current300, err := proc.ReadInt16(uintptr(playContainer38 + 0x86)) //2 bytes
	if err != nil {
		//		log.Println("CurrentHit300c result pointer failure")
		return -5
	}
	//currentgeki, err := proc.ReadInt16(uintptr(comboSecondLevel + 0x8A)) //2 bytes
	//current300Result := current300 + currentgeki // thats not how the game works /shrug
	return current300
}
func CurrentHit50c() int16 {
	if osuStatus != 2 {
		return -1
	}

	currentCombo, err := proc.ReadInt16(uintptr(playContainer38 + 0x88)) //2 bytes
	if err != nil {
		//	log.Println("CurrentHitMiss result pointer failure")
		return -5
	}
	return currentCombo
}
func CurrentHitMiss() int16 {
	if osuStatus != 2 {
		return -1
	}

	currentCombo, err := proc.ReadInt16(uintptr(playContainer38 + 0x8E)) //2 bytes
	if err != nil {
		//		log.Println("CurrentHitMiss result pointer failure")
		return -5
	}
	return currentCombo
}
func CurrentScore() int32 {
	if osuStatus != 2 {
		return -1
	}

	currentCombo, err := proc.ReadInt32(uintptr(playContainer38 + 0x74))
	if err != nil {
		//	log.Println("CurrentScore result pointer failure")
		return -5
	}
	return currentCombo
}
func CurrentGameMode() int32 {
	if osuStatus != 2 {
		return -1
	}

	currentCombo, err := proc.ReadInt32(uintptr(playContainer38 + 0x64))
	if err != nil {
		//		log.Println("GameMode result pointer failure")
		return -5
	}
	return currentCombo
}
func CurrentAccuracy() float64 {
	if osuStatus != 2 {
		return -1
	}

	comboSecondLevel, err := proc.ReadUint32(uintptr(playContainerFirstlevel) + 0x48)
	if err != nil {
		//	log.Println("Accuracy Second level pointer failure")
		return -4
	}
	currentCombo, err := proc.ReadFloat64(uintptr(comboSecondLevel + 0x14))
	if err != nil {
		//		log.Println("Accuracy result pointer failure")
		return -5
	}
	return currentCombo
}
func CurrentPlayerHP() int8 {
	if osuStatus != 2 {
		return -1
	}

	comboSecondLevel, err := proc.ReadUint32(uintptr(playContainerFirstlevel) + 0x40)
	if err != nil {
		//		log.Println("CurrentPlayerHP Second level pointer failure")
		return -4
	}
	currentCombo, err := proc.ReadInt8(uintptr(comboSecondLevel + 0x1C))
	if err != nil {
		//	log.Println("CurrentPlayerHP result pointer failure")
		return -5
	}
	return currentCombo
}
func CurrentPlayerHPSmoothed() int8 {
	if osuStatus != 2 {
		return -1
	}

	comboSecondLevel, err := proc.ReadUint32(uintptr(playContainerFirstlevel) + 0x40)
	if err != nil {
		//	log.Println("CurrentPlayerHPSmoothed Second level pointer failure")
		return -4
	}
	currentCombo, err := proc.ReadInt8(uintptr(comboSecondLevel + 0x14))
	if err != nil {
		//		log.Println("CurrentPlayerHPSmoothed result pointer failure")
		return -5
	}
	return currentCombo
}
func CurrentPlayTime() int32 {
	playTimeFirstLevel, err := proc.ReadUint32(playTime)
	if err != nil {
		//	log.Println("playTime Base level failure")
		return -1
	}
	playTimeValue, err := proc.ReadUint32(uintptr(playTimeFirstLevel))
	if err != nil {
		//	log.Println("playTime Result level failure")
		return -1
	}

	return cast.ToInt32(playTimeValue)
}
func PP() string {
	if operatingSystem == 1 {
		calc, err := exec.Command("oppai.exe", fullPathToOsu, "-end"+lastObject, ppAcc+"%", ppCombo+"x", ppMiss+"m", pp100+"x100", pp50+"x50", "+"+ppMods, "-ojson").Output()
		if err != nil {
			fmt.Println(err)
		}

		return strings.ToValidUTF8(cast.ToString(calc), "")
	} else {
		calc := Cmd("oppai"+" "+"\""+fullPathToOsu+"\""+" "+"-end"+lastObject+" "+ppAcc+"%"+" "+ppCombo+"x"+" "+ppMiss+"m"+" "+pp100+"x100"+" "+pp50+"x50"+" "+"+"+ppMods+" "+"-ojson", true)
		//calc := Cmd("oppai"+" "+"\""+fullPathToOsu+"\""+" "+"-end"+lastObject+" "+ppAcc+"%"+" "+ppCombo+"x"+" "+ppMiss+"m"+" "+pp100+"x100"+" "+pp50+"x50"+" "+" "+"-ojson", true)
		return strings.ToValidUTF8(cast.ToString(calc), "")
	}

}
func PPifFC() string {
	if operatingSystem == 1 {
		calc, err := exec.Command("oppai.exe", fullPathToOsu, ppAcc+"%", pp100+"x100", pp50+"x50", "+"+ppMods, "-ojson").Output()
		if err != nil {
			fmt.Println(err)
		}

		return strings.ToValidUTF8(cast.ToString(calc), "")
	} else {
		calc := Cmd("oppai"+" "+"\""+fullPathToOsu+"\""+" "+" "+ppAcc+"%"+" "+pp100+"x100"+" "+pp50+"x50"+" "+"+"+ppMods+" "+"-ojson", true)

		return strings.ToValidUTF8(cast.ToString(calc), "")
	}

}
func PPSS() string {
	if operatingSystem == 1 {
		calc, err := exec.Command("oppai.exe", fullPathToOsu, "100%", "+"+ppMods, "-ojson").Output()
		if err != nil {
			fmt.Println(err)
		}

		return strings.ToValidUTF8(cast.ToString(calc), "")
	} else {
		calc := Cmd("oppai"+" "+"\""+fullPathToOsu+"\""+" "+"100%"+" "+"+"+ppMods+" "+"-ojson", true)

		return strings.ToValidUTF8(cast.ToString(calc), "")
	}

}
func PP99() string {
	if operatingSystem == 1 {
		calc, err := exec.Command("oppai.exe", fullPathToOsu, "99%", "+"+ppMods, "-ojson").Output()
		if err != nil {
			fmt.Println(err)
		}

		return strings.ToValidUTF8(cast.ToString(calc), "")
	} else {
		calc := Cmd("oppai"+" "+"\""+fullPathToOsu+"\""+" "+"99%"+" "+"+"+ppMods+" "+"-ojson", true)

		return strings.ToValidUTF8(cast.ToString(calc), "")
	}

}
func PP98() string {
	if operatingSystem == 1 {
		calc, err := exec.Command("oppai.exe", fullPathToOsu, "98%", "+"+ppMods, "-ojson").Output()
		if err != nil {
			fmt.Println(err)
		}

		return strings.ToValidUTF8(cast.ToString(calc), "")
	} else {
		calc := Cmd("oppai"+" "+"\""+fullPathToOsu+"\""+" "+"98%"+" "+"+"+ppMods+" "+"-ojson", true)

		return strings.ToValidUTF8(cast.ToString(calc), "")
	}

}
func PP97() string {
	if operatingSystem == 1 {
		calc, err := exec.Command("oppai.exe", fullPathToOsu, "97%", "+"+ppMods, "-ojson").Output()
		if err != nil {
			fmt.Println(err)
		}

		return strings.ToValidUTF8(cast.ToString(calc), "")
	} else {
		calc := Cmd("oppai"+" "+"\""+fullPathToOsu+"\""+" "+"97%"+" "+"+"+ppMods+" "+"-ojson", true)

		return strings.ToValidUTF8(cast.ToString(calc), "")
	}

}
func PP96() string {
	if operatingSystem == 1 {
		calc, err := exec.Command("oppai.exe", fullPathToOsu, "96%", "+"+ppMods, "-ojson").Output()
		if err != nil {
			fmt.Println(err)
		}

		return strings.ToValidUTF8(cast.ToString(calc), "")
	} else {
		calc := Cmd("oppai"+" "+"\""+fullPathToOsu+"\""+" "+"96%"+" "+"+"+ppMods+" "+"-ojson", true)

		return strings.ToValidUTF8(cast.ToString(calc), "")
	}

}
func PP95() string {
	if operatingSystem == 1 {
		calc, err := exec.Command("oppai.exe", fullPathToOsu, "95%", "+"+ppMods, "-ojson").Output()
		if err != nil {
			fmt.Println(err)
		}

		return strings.ToValidUTF8(cast.ToString(calc), "")
	} else {
		calc := Cmd("oppai"+" "+"\""+fullPathToOsu+"\""+" "+"95%"+" "+"+"+ppMods+" "+"-ojson", true)

		return strings.ToValidUTF8(cast.ToString(calc), "")
	}

}
func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

// ModsResolver is just a placeholder for now, needs proper logic or at least switch statements
func ModsResolver(xor uint32) string {
	if xor >= 2048 {
		xor = xor - 2048 //autoplay hack, TODO: Refactor Only works with STD
	}
	NoMod := uint32(0)
	NoFail := uint32(1) << 0
	Easy := uint32(1) << 1
	//TouchDevice := uint32(1) << 2
	Hidden := uint32(1) << 3
	HardRock := uint32(1) << 4
	SuddenDeath := uint32(1) << 5
	DoubleTime := uint32(1) << 6
	Relax := uint32(1) << 7
	HalfTime := uint32(1) << 8
	Nightcore := uint32(1) << 9
	Flashlight := uint32(1) << 10
	Autoplay := uint32(1) << 11
	SpunOut := uint32(1) << 12
	AutoPilot := uint32(1) << 13
	Perfect := uint32(1) << 14
	ScoreV2 := uint32(1) << 29

	if xor == NoMod {
		return "NM"
	}
	if xor == NoFail {
		return "NF"
	}
	if xor == Easy {
		return "EZ"
	}
	if xor == Hidden {
		return "HD"
	}
	if xor == HardRock {
		return "HR"
	}
	if xor == SuddenDeath {
		return "SD"
	}
	if xor == DoubleTime {
		return "DT"
	}
	if xor == Relax {
		return ""
	}
	if xor == HalfTime {
		return "HT"
	}
	if xor == Nightcore {
		return "NC"
	}
	if xor == Flashlight {
		return "FL"
	}
	if xor == Autoplay {
		return ""
	}
	if xor == SpunOut {
		return ""
	}
	if xor == AutoPilot {
		return ""
	}
	if xor == Perfect {
		return ""
	}
	if xor == ScoreV2 {
		return "" // we actually support that
	}
	if xor == Hidden+DoubleTime {
		return "HDDT" // we actually support that
	}
	if xor == Hidden+HardRock {
		return "HDHR" // we actually support that
	}
	if xor == Hidden+DoubleTime+HardRock {
		return "HDDTHR" // we actually support that
	}
	if xor == Hidden+DoubleTime+HardRock+Flashlight {
		return "HDHRDTFL" // we actually support that
	}
	if xor == Hidden+DoubleTime+Easy {
		return "EZHDDT" // we actually support that
	}
	if xor == DoubleTime+Easy {
		return "EZDT" // we actually support that
	}
	if xor == Hidden+Easy {
		return "EZHD" // we actually support that
	}
	return "NM"

}
func HTTPServer() {

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	//box := packr.NewBox("./index")
	http.Handle("/Songs/", http.StripPrefix("/Songs/", http.FileServer(http.Dir(workingDirectory))))
	//	http.Handle("/", http.FileServer(box))
	http.HandleFunc("/json", handler)
	http.ListenAndServe(":24050", nil)
}
func handler(w http.ResponseWriter, r *http.Request) {
	if isRunning == 1 {
		fmt.Fprintf(w, cast.ToString(jsonByte))

	} else {
		fmt.Fprintf(w, "osu! is not fully loaded!")
	}

}
