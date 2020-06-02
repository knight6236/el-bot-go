package eltype

import (
	"el-bot-go/src/gomirai"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// 「」

// Controller 控制器类，作为整个机器人的中心调度模块。
// @property	configReader	ConfigReader	配置读取类
// @property	bot				*gomirai.Bot	机器人
type Controller struct {
	mute         sync.RWMutex
	configReader *ConfigReader
	cronChecker  *CronChecker
	rssListener  *RssListener
	freqMonitor  *FreqMonitor
	bot          *gomirai.Bot
}

var handlerConstructor = [...]func(configList []Config, messageList []Message, operationList []Operation,
	preDefVarMap *map[string]string) (IHandler, error){
	NewPlainHandler, NewImageHandler, NewOperationHandler, NewFaceHandler, NewXMLHandler}

var doerConstructor = [...]func(configHitList []Config, recivedMessageList []Message,
	preDefVarMap map[string]string) (IDoer, error){
	NewPlainDoer, NewImageDoer, NewOperationDoer, NewFaceDoer, NewXMLDoer}

// NewController 构造一个 Controller
// @param	configReader	ConfigReader	配置读取类
func NewController(configReader *ConfigReader, bot *gomirai.Bot) *Controller {
	controller := new(Controller)
	controller.configReader = configReader
	controller.bot = bot
	controller.cronChecker, _ = NewCronChecker(configReader.CronConfigList)
	controller.rssListener, _ = NewRssListener(configReader.RssConfigList)
	controller.freqMonitor, _ = NewFreqMonitor(configReader.FreqUpperLimit)
	controller.cronChecker.Start()
	controller.rssListener.Start()
	controller.freqMonitor.Start()
	go controller.monitorFolder()
	go controller.listenCron()
	go controller.listenRss()
	return controller
}

// Commit 将事件提交给 Controller
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

	controller.mute.RLock()

	configRelatedList := controller.getConfigRelatedList(event)

	configHitList := controller.getConfigHitList(event, configRelatedList)

	controller.mute.RUnlock()

	event.addPerDefVar("el-count-overall",
		strings.Replace(fmt.Sprintf("%v", controller.freqMonitor.CountMap), "map", "统计概要", 1))

	controller.sendMessageAndOperation(event, configHitList)

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
	configSet := make(map[int64]bool)
	var configHitList []Config
	for i := 0; i < len(handlerConstructor); i++ {
		handler, err := (handlerConstructor[i](configRelatedList, event.MessageList, event.OperationList, &event.PreDefVarMap))
		if err != nil {
			continue
		}

		for _, config := range handler.GetConfigHitList() {
			if !configSet[config.innerID] {
				config = config.DeepCopy()
				config.CompleteType()
				config.CompleteContent(event)
				controller.freqMonitor.Commit(config)
				if !controller.convertToUnBlockConfig(&config) {
					configHitList = append(configHitList, config)
					configSet[config.innerID] = true
				}
			}
		}
	}
	return configHitList
}

func (controller *Controller) convertToUnBlockConfig(configHit *Config) bool {
	isAllBlocked := true
	for i := 0; i < len(configHit.Do.Message.Receiver.GroupIDList); i++ {
		isBlocked := controller.freqMonitor.IsBlocked(configHit.innerID,
			ReceiverTypeGroup, CastStringToInt64(configHit.Do.Message.Receiver.GroupIDList[i]))
		isAllBlocked = isAllBlocked && isBlocked
		if isBlocked {
			configHit.Do.Message.Receiver.GroupIDList[i] = "0"
		}
	}
	for i := 0; i < len(configHit.Do.Message.Receiver.UserIDList); i++ {
		isBlocked := controller.freqMonitor.IsBlocked(configHit.innerID,
			ReceiverTypeUser, CastStringToInt64(configHit.Do.Message.Receiver.UserIDList[i]))
		isAllBlocked = isAllBlocked && isBlocked
		if isBlocked {
			configHit.Do.Message.Receiver.UserIDList[i] = "0"
		}
	}
	for i := 0; i < len(configHit.Do.OperationList); i++ {
		switch configHit.Do.OperationList[i].innerType {
		case OperationTypeAt:
			isBlocked := controller.freqMonitor.IsBlocked(configHit.innerID,
				ReceiverTypeGroup,
				CastStringToInt64(configHit.Do.OperationList[i].GroupID))
			isAllBlocked = isAllBlocked && isBlocked
			if isBlocked {
				configHit.Do.OperationList[i].GroupID = "0"
			}
		case OperationTypeAtAll:
			isBlocked := controller.freqMonitor.IsBlocked(configHit.innerID,
				ReceiverTypeGroup,
				CastStringToInt64(configHit.Do.OperationList[i].GroupID))
			isAllBlocked = isAllBlocked && isBlocked
			if isBlocked {
				configHit.Do.OperationList[i].GroupID = "0"
			}
		}
	}
	return isAllBlocked
}

func (controller *Controller) sendMessageAndOperation(event Event, configHitList []Config) {
	willBeSentGoMiraiGroupMessageMap := make(map[int64]map[int64][]gomirai.Message)
	willBeSentGoMiraiUserMessageMap := make(map[int64]map[int64][]gomirai.Message)
	for i := 0; i < len(doerConstructor); i++ {
		doer, err := (doerConstructor[i](configHitList, event.MessageList, event.PreDefVarMap))
		if err != nil {
			continue
		}

		for _, message := range doer.GetWillBeSentMessageList() {
			message.CompleteType()
			message.CompleteContent(event)
			var quoteID int64
			if message.IsQuote {
				quoteID = event.MessageID
			} else {
				quoteID = 0
			}
			goMiraiMessageList, isSuccess := message.ToGoMiraiMessageList()
			if !isSuccess {
				continue
			}
			for _, nativeGroupID := range message.Receiver.GroupIDList {
				groupID := CastStringToInt64(nativeGroupID)
				for _, goMiraiMessage := range goMiraiMessageList {
					if willBeSentGoMiraiGroupMessageMap[groupID] == nil {
						willBeSentGoMiraiGroupMessageMap[groupID] = make(map[int64][]gomirai.Message)
					}
					willBeSentGoMiraiGroupMessageMap[groupID][quoteID] =
						append(willBeSentGoMiraiGroupMessageMap[groupID][quoteID], goMiraiMessage)
				}
			}
			for _, nativeUserID := range message.Receiver.UserIDList {
				userID := CastStringToInt64(nativeUserID)
				for _, goMiraiMessage := range goMiraiMessageList {
					if goMiraiMessage.Type == "At" || goMiraiMessage.Type == "AtAll" {
						continue
					}
					if willBeSentGoMiraiUserMessageMap[userID] == nil {
						willBeSentGoMiraiUserMessageMap[userID] = make(map[int64][]gomirai.Message)
					}
					willBeSentGoMiraiUserMessageMap[userID][quoteID] =
						append(willBeSentGoMiraiUserMessageMap[userID][quoteID], goMiraiMessage)
				}
			}
		}

		for _, operation := range doer.GetWillBeSentOperationList() {
			operation.CompleteType()
			operation.CompleteContent(event)
			controller.sendOperation(operation)
		}
	}
	for receiverID, innerMap := range willBeSentGoMiraiGroupMessageMap {
		for quoteID, willBeSentMessageList := range innerMap {
			controller.sendMessage(ReceiverTypeGroup, receiverID, quoteID, willBeSentMessageList)
		}
	}
	for receiverID, innerMap := range willBeSentGoMiraiUserMessageMap {
		for quoteID, willBeSentMessageList := range innerMap {
			controller.sendMessage(ReceiverTypeUser, receiverID, quoteID, willBeSentMessageList)
		}
	}
}

func (controller *Controller) sendMessage(receiverType ReceiverType, receiverID int64, quoteID int64, willBeSentGoMiraiMessageList []gomirai.Message) {
	switch receiverType {
	case ReceiverTypeGroup:
		_, err := controller.bot.SendGroupMessage(receiverID, quoteID, willBeSentGoMiraiMessageList)
		if err != nil {
			log.Printf("Controller.sendMessage: %s", err.Error())
		}
	case ReceiverTypeUser:
		_, err := controller.bot.SendFriendMessage(receiverID, quoteID, willBeSentGoMiraiMessageList)
		if err != nil {
			log.Printf("Controller.sendMessage: %s", err.Error())
		}
	}
}

func (controller *Controller) sendOperation(operation Operation) {
	groupID := CastStringToInt64(operation.GroupID)
	userID := CastStringToInt64(operation.UserID)
	switch operation.innerType {
	case OperationTypeAt:
		goMiraiMessage, isSuccess := operation.ToGoMiraiMessage()
		if isSuccess {
			controller.bot.SendGroupMessage(CastStringToInt64(operation.GroupID), 0, []gomirai.Message{goMiraiMessage})
		}
	case OperationTypeAtAll:
		goMiraiMessage, isSuccess := operation.ToGoMiraiMessage()
		if isSuccess {
			controller.bot.SendGroupMessage(CastStringToInt64(operation.GroupID), 0, []gomirai.Message{goMiraiMessage})
		}
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
						controller.mute.Lock()
						controller.configReader.reLoad()
						controller.cronChecker, _ = NewCronChecker(controller.configReader.CronConfigList)
						controller.rssListener, _ = NewRssListener(controller.configReader.RssConfigList)
						controller.freqMonitor, _ = NewFreqMonitor(controller.configReader.FreqUpperLimit)
						controller.mute.Unlock()
						log.Println("检测到配置目录下的文件被创建，已经自动更新配置。")
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						controller.mute.Lock()
						controller.configReader.reLoad()
						controller.cronChecker, _ = NewCronChecker(controller.configReader.CronConfigList)
						controller.rssListener, _ = NewRssListener(controller.configReader.RssConfigList)
						controller.freqMonitor, _ = NewFreqMonitor(controller.configReader.FreqUpperLimit)
						controller.mute.Unlock()
						log.Println("检测到配置目录下的文件被修改，已经自动更新配置。")
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						controller.mute.Lock()
						controller.configReader.reLoad()
						controller.cronChecker, _ = NewCronChecker(controller.configReader.CronConfigList)
						controller.rssListener, _ = NewRssListener(controller.configReader.RssConfigList)
						controller.freqMonitor, _ = NewFreqMonitor(controller.configReader.FreqUpperLimit)
						controller.mute.Unlock()
						log.Println("检测到配置目录下的文件被移除，已经自动更新配置。")
					}
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						controller.mute.Lock()
						controller.configReader.reLoad()
						controller.cronChecker, _ = NewCronChecker(controller.configReader.CronConfigList)
						controller.rssListener, _ = NewRssListener(controller.configReader.RssConfigList)
						controller.freqMonitor, _ = NewFreqMonitor(controller.configReader.FreqUpperLimit)
						controller.mute.Unlock()
						log.Println("检测到配置目录下的文件被重命名，已经自动更新配置。")
					}
					if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
						controller.mute.Lock()
						controller.configReader.reLoad()
						controller.cronChecker, _ = NewCronChecker(controller.configReader.CronConfigList)
						controller.rssListener, _ = NewRssListener(controller.configReader.RssConfigList)
						controller.freqMonitor, _ = NewFreqMonitor(controller.configReader.FreqUpperLimit)
						controller.mute.Unlock()
						log.Println("检测到配置目录下的文件权限变化，已经自动更新配置。")
					}
				}
			case err := <-watch.Errors:
				{
					log.Printf("Controller.monitorFolder: %s", err.Error())
					return
				}
			}
		}
	}()

	//循环
	select {}
}

func (controller *Controller) listenCron() {
	for true {
		config := <-controller.cronChecker.WillBeSentConfig
		controller.sendMessageAndOperation(Event{}, []Config{config})
	}
}

func (controller *Controller) listenRss() {
	for true {
		config := <-controller.rssListener.WillBeSentConfig
		event := <-controller.rssListener.WillBeUsedEvent
		controller.sendMessageAndOperation(event, []Config{config})
	}
}
