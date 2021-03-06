package main

import (
	"el-bot-go/src/eltype"
	"el-bot-go/src/gomirai"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	eltype.ConfigRoot = os.Getenv("CONFIG_ROOT")
	eltype.SettingFullPath = os.Getenv("SETTING_FILE")
	eltype.FaceMapFullPath = os.Getenv("FACE_MAP_FILE")
	eltype.ImageFolder = os.Getenv("IMAGE_FOLDER")
	eltype.DataRoot = os.Getenv("DATA_ROOT")
	eltype.DefaultConfigFileName = os.Getenv("DEFAULT_CONFIG_FILE_NAME")
	eltype.RssDataFileName = os.Getenv("RSS_DATA_FILE_NAME")

	switch len(os.Args) {
	case 0:
		log.Println("缺少启动参数「QQ号」和「自定义配置目录（相对于 config 目录）」")
	case 1:
		log.Println("缺少启动参数「QQ号」")
	case 2:
		log.Println("缺少启动参数和「自定义配置目录（相对于 config 目录）」")
	}

	reader := eltype.NewConfigReader(os.Args[2])
	reader.Load(true)

	address := "http://127.0.0.1:" + reader.Port
	authKey := reader.AuthKey
	// 用于进行网络操作的Client
	client := gomirai.NewMiraiClient(address, authKey)

	// 可对Client做出自定义修改，该修改会应用于所有使用该client的网络请求
	// 如使用Proxy
	client.HTTPClient.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}

	// 获取Bot，Session信息保存在Bot中
	// 也可通过Client.Bots[]获取
	qq, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		log.Printf("获取 QQ 号失败，可能是启动参数有误 %s。\n", os.Args[1])
	}
	bot, err := client.Verify(qq)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 释放资源
	defer bot.Release()

	// 初始化消息通道
	// FetchMessage时间间隔 1s,每次获取的数量20,channel缓存容量20
	bot.InitChannel(20, time.Second)

	// 在协程中开始获取消息，消息传输至Bot.MessageChan
	// 忽略错误
	go bot.FetchMessage()
	// 检查错误
	go func() {
		err = bot.FetchMessage()
		if err != nil {
			//handle Error
		}
	}()

	controller := eltype.NewController(&reader, bot)
	fmt.Println("启动成功")

	// 从bot.MessageChan获取收到事件并处理
	for {
		e := <-bot.MessageChan
		switch e.Type {
		case "GroupMessage": // do something
			go controller.Commit(e)
		case "FriendMessage": // do something
			go controller.Commit(e)
		case "GroupMuteAllEvent": // do something
			go controller.Commit(e)
		case "MemberMuteEvent":
			go controller.Commit(e)
		case "MemberUnmuteEvent":
			go controller.Commit(e)
		case "MemberJoinEvent":
			go controller.Commit(e)
		case "MemberLeaveEventKick":
			go controller.Commit(e)
		case "MemberLeaveEventQuit":
			go controller.Commit(e)
		default:
			// do something
		}
	}
}
