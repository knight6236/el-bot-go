package eltype

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/ADD-SP/gomirai"

	"github.com/fsnotify/fsnotify"
)

// 「」

// Controller 控制器类，作为整个机器人的中心调度模块。
// @property	configReader	ConfigReader	配置读取类
// @property	bot				*gomirai.Bot	机器人
type Controller struct {
	isSuspend         bool
	blockedGroupIDSet map[string]bool
	blockedUserIDSet  map[string]bool
	controlMute       sync.RWMutex
	configMute        sync.RWMutex
	configReader      *ConfigReader
	cronChecker       *CronChecker
	rssListener       *RssListener
	freqMonitor       *FreqMonitor
	bot               *gomirai.Bot
}

var handlerConstructor = [...]func(configList []Config, messageList []Message, operationList []Operation,
	preDefVarMap *map[string]string) (IHandler, error){
	NewPlainHandler, NewImageHandler, NewOperationHandler, NewFaceHandler, NewXMLHandler}

var doerConstructor = [...]func(configHitList []Config, recivedMessageList []Message,
	preDefVarMap map[string]string) (IDoer, error){
	NewPlainDoer, NewImageDoer, NewOperationDoer, NewFaceDoer, NewXMLDoer, NewControlDoer}

// NewController 构造一个 Controller
// @param	configReader	ConfigReader	配置读取类
func NewController(configReader *ConfigReader, bot *gomirai.Bot) *Controller {
	controller := new(Controller)
	controller.configReader = configReader
	controller.bot = bot
	controller.blockedGroupIDSet = make(map[string]bool)
	controller.blockedUserIDSet = make(map[string]bool)
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

	configRelatedList := controller.getConfigRelatedList(event)

	configHitList := controller.getConfigHitList(event, configRelatedList)

	event.AddPerDefVar("el-count-overall",
		strings.Replace(fmt.Sprintf("%v", controller.freqMonitor.CountMap), "map", "统计概要", 1))

	controller.sendMessageAndOperation(event, configHitList)

}

func (controller *Controller) getConfigRelatedList(event Event) []Config {
	controller.configMute.RLock()
	var configList []Config
	switch event.Type {
	case EventTypeGroupMessage:
		MergeConfigList(&configList,
			controller.getConfigRelatedConfigList(event.Type, controller.configReader.GlobalConfigList, event.Sender),
			controller.getConfigRelatedConfigList(event.Type, controller.configReader.GroupConfigList, event.Sender))
	case EventTypeFriendMessage:
		MergeConfigList(&configList,
			controller.getConfigRelatedConfigList(event.Type, controller.configReader.GlobalConfigList, event.Sender),
			controller.getConfigRelatedConfigList(event.Type, controller.configReader.FriendConfigList, event.Sender))
	default:
		MergeConfigList(&configList,
			controller.getConfigRelatedConfigList(event.Type, controller.configReader.GlobalConfigList, event.Sender),
			controller.getConfigRelatedConfigList(event.Type, controller.configReader.GroupConfigList, event.Sender))
	}
	controller.configMute.RUnlock()
	return configList
}

func (controller *Controller) getConfigRelatedConfigList(eventType EventType, configList []Config, sender Sender) []Config {
	controller.configMute.RLock()
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
	controller.configMute.RUnlock()
	return ret
}

func (controller *Controller) getConfigHitList(event Event, configRelatedList []Config) []Config {
	controller.configMute.RLock()
	configSet := make(map[int64]bool)
	var configHitList []Config
	for i := 0; i < len(handlerConstructor); i++ {
		handler, err := (handlerConstructor[i](configRelatedList, event.MessageList, event.OperationList, &event.PreDefVarMap))
		if err != nil {
			continue
		}

		for _, config := range handler.GetConfigHitList() {
			if !configSet[config.InnerID] {
				config = config.DeepCopy()
				config.CompleteType()
				config.CompleteContent(event)
				controller.freqMonitor.Commit(config)
				if !controller.convertToUnBlockConfig(&config) {
					configHitList = append(configHitList, config)
					configSet[config.InnerID] = true
				}
			}
		}
	}
	controller.configMute.RUnlock()
	return configHitList
}

func (controller *Controller) convertToUnBlockConfig(configHit *Config) bool {
	isAllBlocked := true
	for i := 0; i < len(configHit.Do.Message.Receiver.GroupIDList); i++ {
		isBlocked := controller.freqMonitor.IsBlocked(configHit.InnerID,
			ReceiverTypeGroup, CastStringToInt64(configHit.Do.Message.Receiver.GroupIDList[i]))
		isAllBlocked = isAllBlocked && isBlocked
		if isBlocked {
			configHit.Do.Message.Receiver.GroupIDList[i] = "0"
		}
	}
	for i := 0; i < len(configHit.Do.Message.Receiver.UserIDList); i++ {
		isBlocked := controller.freqMonitor.IsBlocked(configHit.InnerID,
			ReceiverTypeUser, CastStringToInt64(configHit.Do.Message.Receiver.UserIDList[i]))
		isAllBlocked = isAllBlocked && isBlocked
		if isBlocked {
			configHit.Do.Message.Receiver.UserIDList[i] = "0"
		}
	}
	for i := 0; i < len(configHit.Do.OperationList); i++ {
		switch configHit.Do.OperationList[i].InnerType {
		case OperationTypeAt:
			isBlocked := controller.freqMonitor.IsBlocked(configHit.InnerID,
				ReceiverTypeGroup,
				CastStringToInt64(configHit.Do.OperationList[i].GroupID))
			isAllBlocked = isAllBlocked && isBlocked
			if isBlocked {
				configHit.Do.OperationList[i].GroupID = "0"
			}
		case OperationTypeAtAll:
			isBlocked := controller.freqMonitor.IsBlocked(configHit.InnerID,
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
	controller.configMute.RLock()
	willBeSentGoMiraiGroupMessageMap := make(map[int64]map[int64][]gomirai.Message)
	willBeSentGoMiraiUserMessageMap := make(map[int64]map[int64][]gomirai.Message)
	var willBeSentOperaitonList []Operation
	var willBeSentControlList []Control
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
			willBeSentOperaitonList = append(willBeSentOperaitonList, operation)
		}

		for _, control := range doer.GetwillBeSentControlList() {
			willBeSentControlList = append(willBeSentControlList, control)
		}
	}
	controller.configMute.RUnlock()

	for _, operation := range willBeSentOperaitonList {
		operation.CompleteType()
		operation.CompleteContent(event)
		controller.sendOperation(operation)
	}
	for _, control := range willBeSentControlList {
		control.CompleteType()
		control.CompleteContent(event)
		controller.sendControl(control)
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
	controller.configMute.RLock()
	defer controller.configMute.RUnlock()
	if controller.isSuspend {
		return
	}
	if receiverType == ReceiverTypeGroup && controller.blockedGroupIDSet[CastInt64ToString(receiverID)] {
		return
	}
	if receiverType == ReceiverTypeUser && controller.blockedUserIDSet[CastInt64ToString(receiverID)] {
		return
	}
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
	controller.configMute.RLock()
	defer controller.configMute.RUnlock()
	if controller.isSuspend {
		return
	}
	if operation.GroupID != "" && controller.blockedGroupIDSet[operation.GroupID] {
		return
	}
	if operation.UserID != "" && controller.blockedUserIDSet[operation.UserID] {
		return
	}
	groupID := CastStringToInt64(operation.GroupID)
	userID := CastStringToInt64(operation.UserID)
	switch operation.InnerType {
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

func (controller *Controller) sendControl(control Control) {
	controller.configMute.Lock()
	defer controller.configMute.Unlock()

	switch control.InnerType {
	case ControlTypeSuspend:
		controller.isSuspend = true
	case ControlTypeActive:
		controller.isSuspend = false
	case ControlTypeDestory:
		log.Println("接收到终止指令，程序自动终止。")
		os.Exit(0)
	case ControlTypeRestart:
		// TODO
	case ControlTypeBlock:
		for _, groupID := range control.GroupIDList {
			if groupID == "" {
				continue
			}
			controller.blockedGroupIDSet[groupID] = true
		}
		for _, userID := range control.UserIDList {
			if userID == "" {
				continue
			}
			controller.blockedUserIDSet[userID] = true
		}
	case ControlTypeUnblock:
		for _, groupID := range control.GroupIDList {
			if groupID == "" {
				continue
			}
			controller.blockedGroupIDSet[groupID] = false
		}
		for _, userID := range control.UserIDList {
			if userID == "" {
				continue
			}
			controller.blockedUserIDSet[userID] = false
		}
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
			isChange := false
			select {
			case ev := <-watch.Events:
				{
					if ev.Op&fsnotify.Create == fsnotify.Create {
						isChange = true
						log.Println("检测到配置目录下的文件被创建，已经自动更新配置。")
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						isChange = true
						log.Println("检测到配置目录下的文件被修改，已经自动更新配置。")
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						isChange = true
						log.Println("检测到配置目录下的文件被移除，已经自动更新配置。")
					}
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						isChange = true
						log.Println("检测到配置目录下的文件被重命名，已经自动更新配置。")
					}
					if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
						isChange = true
						log.Println("检测到配置目录下的文件权限变化，已经自动更新配置。")
					}
					if isChange {
						controller.configMute.Lock()
						controller.configReader.reLoad()
						controller.cronChecker.Stop()
						controller.rssListener.Stop()
						controller.freqMonitor.Stop()
						controller.cronChecker, _ = NewCronChecker(controller.configReader.CronConfigList)
						controller.rssListener, _ = NewRssListener(controller.configReader.RssConfigList)
						controller.freqMonitor, _ = NewFreqMonitor(controller.configReader.FreqUpperLimit)
						controller.cronChecker.Start()
						controller.rssListener.Start()
						controller.freqMonitor.Start()
						controller.configMute.Unlock()
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
		select {
		case config := <-controller.cronChecker.WillBeSentConfig:
			controller.sendMessageAndOperation(Event{PreDefVarMap: map[string]string{"\\n": "\n"}}, []Config{config})
		case signalType := <-controller.cronChecker.Signal:
			if signalType == SingalTypeStop {
				return
			}
		}
	}
}

func (controller *Controller) listenRss() {
	for true {
		select {
		case config := <-controller.rssListener.WillBeSentConfig:
			event := <-controller.rssListener.WillBeUsedEvent
			controller.sendMessageAndOperation(event, []Config{config})
		case signalType := <-controller.rssListener.Signal:
			if signalType == SingalTypeStop {
				return
			}
		}
	}
}
