@echo off
set DEFAULT_FILE=config/default.yml
set SETTING_FILE=plugins/MiraiAPIHTTP/setting.yml
set FACE_MAP_FILE=config/face-map.yml

set WIN_AMD64=bin/main-windows-amd64.exe
set WIN_386=bin/main-windows-386.exe
set DARWIN_AMD64=bin/main-darwin-amd64.bin
set DARWIN_386=bin/main-darwin-386.bin
set LINUX_AMD64=bin/main-linux-amd64.bin
set LINUX_386=bin/main-linux-386.bin
set LINUX_ARM=bin/main-linux-arm.bin

if not exist %WIN_386% (
    go run src/main/main.go %1 %DEFAULT_FILE%
) else (
    %WIN_386%.exe
)