package eltype

import (
	"gomirai"
)

// 「」

// Controller 控制器类，作为整个机器人的中心调度模块。
// @property	configReader	ConfigReader	配置读取类
type Controller struct {
	configReader ConfigReader
}

var handlerConstructor = [...]func(configList []Config, messageList []Message, operationList []Operation,
	preDefVarMap map[string]string) (IHandler, error){
	NewPlainHandler, NewImageHandler, NewOperationHandler, NewFaceHandler}

var doerConstructor = [...]func(configHitList []Config, recivedMessageList []Message,
	preDefVarMap map[string]string) (IDoer, error){
	NewPlainDoer, NewImageDoer, NewOperationDoer, NewFaceDoer}

// NewController 构造一个 Controller
// @param	configReader	ConfigReader	配置读取类
func NewController(configReader ConfigReader) Controller {
	var controller Controller
	controller.configReader = configReader
	return controller
}

// Commit 将事件提交给 Controller
// @param	bot				*gomirai.Bot		机器人
// @param	goMiraiEvent	gomirai.InEvent		事件
func (controller *Controller) Commit(bot *gomirai.Bot, goMiraiEvent gomirai.InEvent) {
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

	configRelatedList := controller.getConfigRelatedList(event.Type)

	configHitList := controller.getConfigHitList(event, configRelatedList)

	sendedGoMiraiMessageList := controller.getSendedGoMiraiMessageList(event, configHitList)

	controller.sendMessage(bot, event, sendedGoMiraiMessageList)

}

func (controller *Controller) mergeList(args ...[]Config) []Config {
	targetList := args[0]
	for i := 1; i < len(args); i++ {
		for _, item := range args[i] {
			targetList = append(targetList, item)
		}
	}
	return targetList
}

func (controller *Controller) getConfigRelatedList(eventType EventType) []Config {
	var configList []Config
	switch eventType {
	case EventTypeGroupMessage:
		configList = controller.mergeList(configList, controller.configReader.GlobalConfigList,
			controller.configReader.GroupConfigList)
	case EventTypeFriendMessage:
		configList = controller.mergeList(configList, controller.configReader.GlobalConfigList,
			controller.configReader.FriendConifgList)
	case EventTypeMemberMute:
		configList = controller.mergeList(configList, controller.configReader.GlobalConfigList,
			controller.configReader.GroupConfigList)
	case EventTypeMemberUnmute:
		configList = controller.mergeList(configList, controller.configReader.GlobalConfigList,
			controller.configReader.GroupConfigList)
	case EventTypeGroupMuteAll:
		configList = controller.mergeList(configList, controller.configReader.GlobalConfigList,
			controller.configReader.GroupConfigList)
	case EventTypeGroupUnMuteAll:
		configList = controller.mergeList(configList, controller.configReader.GlobalConfigList,
			controller.configReader.GroupConfigList)
	case EventTypeMemberJoin:
		configList = controller.mergeList(configList, controller.configReader.GlobalConfigList,
			controller.configReader.GroupConfigList)
	case EventTypeMemberLeaveByKick:
		configList = controller.mergeList(configList, controller.configReader.GlobalConfigList,
			controller.configReader.GroupConfigList)
	case EventTypeMemberLeaveByQuit:
		configList = controller.mergeList(configList, controller.configReader.GlobalConfigList,
			controller.configReader.GroupConfigList)
	}
	return configList
}

func (controller *Controller) getConfigHitList(event Event, configRelatedList []Config) []Config {
	var configHitList []Config
	for i := 0; i < len(handlerConstructor); i++ {
		handler, err := (handlerConstructor[i](configRelatedList, event.MessageList, event.OperationList, event.PreDefVarMap))
		if err != nil {
			continue
		}

		for _, config := range handler.GetConfigHitList() {
			configHitList = append(configHitList, config)
		}
	}
	return configHitList

}

func (controller *Controller) getSendedGoMiraiMessageList(event Event, configHitList []Config) []gomirai.Message {
	var sendedGoMiraiMessageList []gomirai.Message
	for i := 0; i < len(doerConstructor); i++ {
		doer, err := (doerConstructor[i](configHitList, event.MessageList, event.PreDefVarMap))
		if err != nil {
			continue
		}

		for _, message := range doer.GetSendedMessageList() {
			goMiraiMessage, err := message.ToGoMiraiMessage()
			if err != nil {
				continue
			}
			sendedGoMiraiMessageList = append(sendedGoMiraiMessageList, goMiraiMessage)
		}
	}
	return sendedGoMiraiMessageList
}

func (controller *Controller) sendMessage(bot *gomirai.Bot, event Event, sendedGoMiraiMessageList []gomirai.Message) {
	switch event.Type {
	case EventTypeGroupMessage:
		bot.SendGroupMessage(event.SenderList[0].ID, 0, sendedGoMiraiMessageList)
	case EventTypeMemberMute:
		bot.SendGroupMessage(event.SenderList[0].ID, 0, sendedGoMiraiMessageList)
	case EventTypeFriendMessage:
		bot.SendFriendMessage(event.SenderList[0].ID, 0, sendedGoMiraiMessageList)
	case EventTypeMemberUnmute:
		bot.SendGroupMessage(event.SenderList[0].ID, 0, sendedGoMiraiMessageList)
	case EventTypeGroupMuteAll:
		bot.SendGroupMessage(event.SenderList[0].ID, 0, sendedGoMiraiMessageList)
	case EventTypeGroupUnMuteAll:
		bot.SendGroupMessage(event.SenderList[0].ID, 0, sendedGoMiraiMessageList)
	case EventTypeMemberJoin:
		bot.SendGroupMessage(event.SenderList[0].ID, 0, sendedGoMiraiMessageList)
	case EventTypeMemberLeaveByKick:
		bot.SendGroupMessage(event.SenderList[0].ID, 0, sendedGoMiraiMessageList)
	case EventTypeMemberLeaveByQuit:
		bot.SendGroupMessage(event.SenderList[0].ID, 0, sendedGoMiraiMessageList)
	}
}
