@echo off
set CONFIG_ROOT=config
set SETTING_FILE=mirai/plugins/MiraiAPIHTTP/setting.yml
set FACE_MAP_FILE=config/face-map.yml
set IMAGE_FOLDER=mirai/plugins/MiraiAPIHTTP/images
set DATA_ROOT=data
set DEFAULT_CONFIG_FILE_NAME=default.yml
set RSS_DATA_FILE_NAME=rss.yml
set PLUGIN_FOLDER=plugins
set PYCMD=python

set WIN_AMD64=bin/main-windows-amd64.exe
set WIN_386=bin/main-windows-386.exe
set DARWIN_AMD64=bin/main-darwin-amd64.bin
set DARWIN_386=bin/main-darwin-386.bin
set LINUX_AMD64=bin/main-linux-amd64.bin
set LINUX_386=bin/main-linux-386.bin
set LINUX_ARM=bin/main-linux-arm.bin

if not exist %WIN_AMD64% (
    %WIN_AMD64%.exe %1 %2
) else (
    go run src/main/main.go %1 %2
)
