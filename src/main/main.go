package main

import (
	"el-bot-go/src/eltype"
	"el-bot-go/src/gomirai"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

func main() {
	// 链接地址
	buf, err := ioutil.ReadFile("../../plugins/MiraiAPIHTTP/setting.yml")
	if err != nil {
		fmt.Println("读取 plugins/MiraiAPIHTTP/setting.yml 失败")
		return
	}
	var settingMap map[string]interface{}
	yaml.Unmarshal(buf, &settingMap)
	if settingMap["enableWebsocket"] != false {
		fmt.Println("enableWebsocket 应设置为 false")
		return
	}
	address := "http://127.0.0.1:" + strconv.Itoa(settingMap["port"].(int))
	authKey := settingMap["authKey"].(string)
	// 用于进行网络操作的Client
	client := gomirai.NewMiraiClient(address, authKey)

	// 可对Client做出自定义修改，该修改会应用于所有使用该client的网络请求
	// 如使用Proxy
	client.HTTPClient.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}

	// 获取Bot，Session信息保存在Bot中
	// 也可通过Client.Bots[]获取
	qq, errqq := strconv.ParseInt(os.Getenv("QQ"), 10, 64)
	if errqq != nil {
		fmt.Println("获取QQ号失败，请检查环境变量设置是否是否正确。")
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

	reader := eltype.NewConfigReader("../../config/default.yml")
	controller := eltype.NewController(reader, bot)
	fmt.Println("启动成功")

	// 从bot.MessageChan获取收到事件并处理
	for {
		e := <-bot.MessageChan
		switch e.Type {
		case "GroupMessage": // do something
			controller.Commit(e)
		case "FriendMessage": // do something
			controller.Commit(e)
		case "GroupMuteAllEvent": // do something
			controller.Commit(e)
		case "MemberMuteEvent":
			controller.Commit(e)
		case "MemberUnmuteEvent":
			controller.Commit(e)
		case "MemberJoinEvent":
			controller.Commit(e)
		case "MemberLeaveEventKick":
			controller.Commit(e)
		case "MemberLeaveEventQuit":
			controller.Commit(e)
		default:
			// do something
		}
	}
}
