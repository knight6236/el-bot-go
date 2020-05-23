package main

import (
	"el-bot-go/src/eltype"
	"el-bot-go/src/gomirai"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	eltype.DefaultConfigFullPath = os.Getenv("DEFAULT_FILE")
	eltype.SettingFullPath = os.Getenv("SETTING_FILE")
	eltype.FaceMapFullPath = os.Getenv("FACE_MAP_FILE")
	reader := eltype.NewConfigReader(os.Args[2])

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
		fmt.Println("获取 QQ 号失败")
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

	controller := eltype.NewController(reader, bot)
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
