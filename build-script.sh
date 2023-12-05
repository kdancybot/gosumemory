#!/usr/bin/env bash

target=gosumemory-no-window

# # Since go doesn't support cross-compiling of C code (for obvious reasons)
# # you can just comment out all unsupported platforms for yourself
# # (I may or may not write better script later)	
# ^ With removal of Cgo parts of project this comment is obsolete 
platforms=(
	"linux/amd64" "linux/386"
	"windows/amd64" "windows/386"
)

mkdir -p build

for platform in "${platforms[@]}"
do
	platform_split=(${platform//\// })
	GOOS=${platform_split[0]}
	GOARCH=${platform_split[1]}
	output_name=$target'-'$GOOS'-'$GOARCH
	flags=""
	if [ $GOOS = "windows" ]; then
		output_name+='.exe'
		flags+="-ldflags -H=windowsgui"
	fi	
	
	GOOS=$GOOS GOARCH=$GOARCH go build $flags -o build/$output_name
	if [ $? -ne 0 ]; then
   		echo 'An error has occurred! Aborting the script execution...'
		exit 1
	fi
done
