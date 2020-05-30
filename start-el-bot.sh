export SETTING_FILE=plugins/MiraiAPIHTTP/setting.yml
export FACE_MAP_FILE=config/face-map.yml
export IMAGE_FOLDER=plugins/MiraiAPIHTTP/images
export CONFIG_ROOT=config
export DATA_ROOT=data
export DEFAULT_CONFIG_FILE_NAME=default.yml
export RSS_DATA_FILE_NAME=rss.yml

export WIN_AMD64=bin/main-windows-amd64.exe
export WIN_386=bin/main-windows-386.exe
export DARWIN_AMD64=bin/main-darwin-amd64.bin
export DARWIN_386=bin/main-darwin-386.bin
export LINUX_AMD64=bin/main-linux-amd64.bin
export LINUX_386=bin/main-linux-386.bin
export LINUX_ARM=bin/main-linux-arm.bin

if [[ $(uname) == "Darwin" ]] && [[ $(arch) == "i386" ]] && [[ -d $DARWIN_386 ]]; then
  sh $DARWIN_386 $1 $2
elif [[ $(uname) == "Darwin" ]] && [[ $(arch) == "amd64" ]] && [[ -d $DARWIN_AMD64 ]]; then
  sh $DARWIN_AMD64 $1 $2
elif [[ $(uname) == "Linux" ]] && [[ $(arch) == "x86_64" || "amd64" ]] && [[ -d $LINUX_AMD64 ]]; then
  sh $LINUX_AMD64 $1 $2
elif [[ $(uname) == "linux" ]] && [[ $(arch) == "x86" || "i386" ]] && [[ -d $LINUX_386 ]]; then
  sh $LINUX_386 $1 $2
elif [[ $(uname) == "Linux" ]] && [[ $(arch) == "armv7" || "aarch64" ]] && [[ -d $LINUX_ARM ]]; then
  sh $LINUX_ARM $1 $2
else
  go run src/main/main.go $1 $2
fi
