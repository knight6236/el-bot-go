export SETTING_FILE=plugins/MiraiAPIHTTP/setting.yml
export FACE_MAP_FILE=config/face-map.yml
export IMAGE_FOLDER=plugins/MiraiAPIHTTP/images
export CONFIG_ROOT=config

export WIN_AMD64=bin/main-windows-amd64.exe
export WIN_386=bin/main-windows-386.exe
export DARWIN_AMD64=bin/main-darwin-amd64.bin
export DARWIN_386=bin/main-darwin-386.bin
export LINUX_AMD64=bin/main-linux-amd64.bin
export LINUX_386=bin/main-linux-386.bin
export LINUX_ARM=bin/main-linux-arm.bin

if [ -d $LINUX_ARM ];then
./$LINUX_ARM $1 $2
else
go run src/main/main.go $1 $2
fi