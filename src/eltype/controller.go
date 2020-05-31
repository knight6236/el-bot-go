package eltype

import (
	"el-bot-go/src/gomirai"
	"fmt"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// 「」

// Controller 控制器类，作为整个机器人的中心调度模块。
// @property	configReader	ConfigReader	配置读取类
// @property	bot				*gomirai.Bot	机器人
type Controller struct {
	configReader ConfigReader
	cronChecker  CronChecker
	rssListener  RssListener
	bot          *gomirai.Bot
	countMap     map[string]int
}

var handlerConstructor = [...]func(configList []Config, messageList []Message, operationList []Operation,
	preDefVarMap *map[string]string) (IHandler, error){
	NewPlainHandler, NewImageHandler, NewOperationHandler, NewFaceHandler, NewXMLHandler}

var doerConstructor = [...]func(configHitList []Config, recivedMessageList []Message,
	preDefVarMap map[string]string) (IDoer, error){
	NewPlainDoer, NewImageDoer, NewOperationDoer, NewFaceDoer, NewXMLDoer}

// NewController 构造一个 Controller
// @param	configReader	ConfigReader	配置读取类
func NewController(configReader ConfigReader, bot *gomirai.Bot) Controller {
	var controller Controller
	controller.configReader = configReader
	controller.bot = bot
	controller.cronChecker, _ = NewCronChecker(configReader.CronConfigList)
	controller.rssListener, _ = NewRssListener(configReader.RssConfigList)
	controller.countMap = make(map[string]int)
	controller.cronChecker.Start()
	controller.rssListener.Start()
	go controller.monitorFolder()
	go controller.listenCron()
	go controller.listenRss()
	return controller
}

// Commit 将事件提交给 Controller
// @param	bot				*gomirai.Bot		机器人
// @param	goMiraiEvent	gomirai.InEvent		事件
func (controller *Controller) Commit(goMiraiEvent gomirai.InEvent) {
	var err error
	switch CastGoMiraiEventTypeToEventType(goMiraiEvent.Type) {
	case EventTypeFriendMessage:
		err = goMiraiEvent.SenderDetail()
	case EventTypeGroupMessage:
		err = goMiraiEvent.SenderDetail()
	case EventTypeGroupMuteAll:
		err = goMiraiEvent.OperatorDetail()
	case EventTypeMemberMute:
		err = goMiraiEvent.OperatorDetail()
	case EventTypeMemberUnmute:
		err = goMiraiEvent.OperatorDetail()
	}
	if err != nil {
		return
	}

	event, err := NewEventFromGoMiraiEvent(goMiraiEvent)
	if err != nil {
		return
	}

	configRelatedList := controller.getConfigRelatedList(event)

	configHitList := controller.getConfigHitList(event, configRelatedList)

	// fmt.Printf("%v\n", configHitList)

	controller.doCount(configHitList)
	event.addPerDefVar("el-count-overall",
		strings.Replace(fmt.Sprintf("%v", controller.countMap), "map", "统计概要", 1))

	controller.sendMessageAndOperation(event, configHitList)

}

func (controller *Controller) doCount(configHitList []Config) {
	for _, config := range configHitList {
		if config.IsCount {
			controller.countMap[config.CountID]++
		}
	}
}

func (controller *Controller) listenCron() {
	for true {
		config := <-controller.cronChecker.WillBeSentConfig
		controller.sendMessageAndOperation(Event{}, []Config{config})
	}
}

func (controller *Controller) getConfigRelatedList(event Event) []Config {
	var configList []Config
	switch event.Type {
	case EventTypeGroupMessage:
		mergeConfigList(&configList,
			controller.getConfigRelatedConfigList(event.Type, controller.configReader.GlobalConfigList, event.Sender),
			controller.getConfigRelatedConfigList(event.Type, controller.configReader.GroupConfigList, event.Sender))
	case EventTypeFriendMessage:
		mergeConfigList(&configList,
			controller.getConfigRelatedConfigList(event.Type, controller.configReader.GlobalConfigList, event.Sender),
			controller.getConfigRelatedConfigList(event.Type, controller.configReader.FriendConfigList, event.Sender))
	default:
		mergeConfigList(&configList,
			controller.getConfigRelatedConfigList(event.Type, controller.configReader.GlobalConfigList, event.Sender),
			controller.getConfigRelatedConfigList(event.Type, controller.configReader.GroupConfigList, event.Sender))
	}
	return configList
}

func (controller *Controller) getConfigRelatedConfigList(eventType EventType, configList []Config, sender Sender) []Config {
	var ret []Config
	for _, config := range configList {
		if (config.When.Message.Sender.UserIDList == nil || len(config.When.Message.Sender.UserIDList) == 0) &&
			(config.When.Message.Sender.GroupIDList == nil || len(config.When.Message.Sender.GroupIDList) == 0) {
			ret = append(ret, config)
			continue
		}

		switch eventType {
		case EventTypeFriendMessage:
			for _, friendID := range config.When.Message.Sender.UserIDList {
				if friendID == sender.UserIDList[0] {
					ret = append(ret, config)
					goto LOOP_END
				}
			}

		default:
			for _, groupID := range config.When.Message.Sender.GroupIDList {
				if groupID == sender.GroupIDList[0] {
					ret = append(ret, config)
					goto LOOP_END
				}
			}
		}

	LOOP_END:
	}
	return ret
}

func (controller *Controller) getConfigHitList(event Event, configRelatedList []Config) []Config {
	configSet := make(map[int]bool)
	var configHitList []Config
	for i := 0; i < len(handlerConstructor); i++ {
		handler, err := (handlerConstructor[i](configRelatedList, event.MessageList, event.OperationList, &event.PreDefVarMap))
		if err != nil {
			continue
		}

		for _, config := range handler.GetConfigHitList() {
			if !configSet[config.innerID] {
				configHitList = append(configHitList, config)
				configSet[config.innerID] = true
			}
		}
	}
	return configHitList
}

func (controller *Controller) sendMessageAndOperation(event Event, configHitList []Config) {
	willBeSentGoMiraiGroupMessageMap := make(map[int64][]gomirai.Message)
	willBeSentGoMiraiUserMessageMap := make(map[int64][]gomirai.Message)
	for i := 0; i < len(doerConstructor); i++ {
		doer, err := (doerConstructor[i](configHitList, event.MessageList, event.PreDefVarMap))
		if err != nil {
			continue
		}

		for _, message := range doer.GetWillBeSentMessageList() {
			message.CompleteType()
			message.CompleteContent(event.PreDefVarMap)
			goMiraiMessageList, isSuccess := message.ToGoMiraiMessageList()
			if !isSuccess {
				continue
			}
			if (message.Receiver.GroupIDList == nil || len(message.Receiver.GroupIDList) == 0) &&
				(message.Receiver.UserIDList == nil || len(message.Receiver.UserIDList) == 0) {
				switch event.Type {
				case EventTypeFriendMessage:
					message.Receiver.UserIDList = append(message.Receiver.UserIDList, event.Sender.UserIDList[0])
				default:
					message.Receiver.GroupIDList = append(message.Receiver.GroupIDList, event.Sender.GroupIDList[0])
				}
			}
			for _, nativeGroupID := range message.Receiver.GroupIDList {
				groupID := CastStringToInt64(nativeGroupID)
				for _, goMiraiMessage := range goMiraiMessageList {
					willBeSentGoMiraiGroupMessageMap[groupID] = append(willBeSentGoMiraiGroupMessageMap[groupID], goMiraiMessage)
				}
			}
			for _, nativeUserID := range message.Receiver.UserIDList {
				userID := CastStringToInt64(nativeUserID)
				for _, goMiraiMessage := range goMiraiMessageList {
					if goMiraiMessage.Type == "At" || goMiraiMessage.Type == "AtAll" {
						continue
					}
					willBeSentGoMiraiUserMessageMap[userID] = append(willBeSentGoMiraiUserMessageMap[userID], goMiraiMessage)
				}
			}
		}

		for _, operation := range doer.GetWillBeSentOperationList() {
			operation.CompleteContent(event.PreDefVarMap)
			controller.sendOperation(operation)
		}
	}
	for receiverID, willBeSentMessageList := range willBeSentGoMiraiGroupMessageMap {
		controller.sendMessage(ReceiverTypeGroup, receiverID, willBeSentMessageList)
	}
	for receiverID, willBeSentMessageList := range willBeSentGoMiraiUserMessageMap {
		controller.sendMessage(ReceiverTypeUser, receiverID, willBeSentMessageList)
	}
}

func (controller *Controller) sendMessage(receiverType ReceiverType, receiverID int64, willBeSentGoMiraiMessageList []gomirai.Message) {
	switch receiverType {
	case ReceiverTypeGroup:
		_, err := controller.bot.SendGroupMessage(receiverID, 0, willBeSentGoMiraiMessageList)
		if err != nil {
			log.Println(err)
		}
	case ReceiverTypeUser:
		_, err := controller.bot.SendFriendMessage(receiverID, 0, willBeSentGoMiraiMessageList)
		if err != nil {
			log.Println(err)
		}
	}
}

func (controller *Controller) sendOperation(operation Operation) {
	groupID := CastStringToInt64(operation.GroupID)
	userID := CastStringToInt64(operation.UserID)

	switch operation.innerType {
	case OperationTypeMemberMute:
		controller.bot.Mute(groupID, userID, CastStringToInt64(operation.Second))
	case OperationTypeMemberUnMute:
		controller.bot.Mute(groupID, userID, CastStringToInt64(operation.Second))
	case OperationTypeGroupMuteAll:
		controller.bot.MuteAll(groupID)
	case OperationTypeGroupUnMuteAll:
		controller.bot.UnmuteAll(groupID)
	}
}

func (controller *Controller) monitorFolder() {
	if controller.configReader.folder == "default" {
		return
	}
	//创建一个监控对象
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watch.Close()
	//添加要监控的对象，文件或文件夹
	err = watch.Add(ConfigRoot + "/" + controller.configReader.folder)
	if err != nil {
		log.Fatal(err)
	}
	//我们另启一个goroutine来处理监控对象的事件
	go func() {
		for {
			select {
			case ev := <-watch.Events:
				{
					if ev.Op&fsnotify.Create == fsnotify.Create {
						controller.configReader.reLoad()
						log.Println("检测到配置目录下的文件被创建，已经自动更新配置。")
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						controller.configReader.reLoad()
						log.Println("检测到配置目录下的文件被修改，已经自动更新配置。")
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						controller.configReader.reLoad()
						log.Println("检测到配置目录下的文件被移除，已经自动更新配置。")
					}
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						controller.configReader.reLoad()
						log.Println("检测到配置目录下的文件被重命名，已经自动更新配置。")
					}
					if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
						controller.configReader.reLoad()
						log.Println("检测到配置目录下的文件权限变化，已经自动更新配置。")
					}
				}
			case err := <-watch.Errors:
				{
					log.Println("error : ", err)
					return
				}
			}
		}
	}()

	//循环
	select {}
}

func (controller *Controller) listenRss() {
	for true {
		config := <-controller.rssListener.willBeSentConfig
		event := <-controller.rssListener.willBeUsedEvent
		controller.sendMessageAndOperation(event, []Config{config})
	}
}
