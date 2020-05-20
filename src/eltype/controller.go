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
	NewPlainHandler, NewImageHandler, NewOperationHandler, NewFaceHandler, NewXMLHandler}

var doerConstructor = [...]func(configHitList []Config, recivedMessageList []Message,
	preDefVarMap map[string]string) (IDoer, error){
	NewPlainDoer, NewImageDoer, NewOperationDoer, NewFaceDoer, NewXMLDoer}

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

	configRelatedList := controller.getConfigRelatedList(event)

	configHitList := controller.getConfigHitList(event, configRelatedList)

	sendedGoMiraiMessageList := controller.getSendedGoMiraiMessageList(event, configHitList)

	controller.sendMessage(bot, event, configHitList, sendedGoMiraiMessageList)

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

func (controller *Controller) getConfigRelatedList(event Event) []Config {
	var configList []Config
	switch event.Type {
	case EventTypeGroupMessage:
		configList = controller.mergeList(configList,
			controller.getConfigRelatedListByWhenSenderList(controller.configReader.GlobalConfigList, event.SenderList),
			controller.getConfigRelatedListByWhenSenderList(controller.configReader.GroupConfigList, event.SenderList))
	case EventTypeFriendMessage:
		configList = controller.mergeList(configList,
			controller.getConfigRelatedListByWhenSenderList(controller.configReader.GlobalConfigList, event.SenderList),
			controller.getConfigRelatedListByWhenSenderList(controller.configReader.FriendConifgList, event.SenderList))
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

func (controller *Controller) getConfigRelatedListByWhenSenderList(configList []Config, senderList []Sender) []Config {
	var ret []Config
	for _, config := range configList {
		for _, sender := range config.SenderList {
			if (sender.Type == SenderTypeGroup && sender.ID == senderList[0].ID) ||
				(sender.Type == SenderTypeUser && sender.ID == senderList[1].ID) {
				ret = append(ret, config)
				break
			}
		}
	}
	return ret
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

func (controller *Controller) sendMessage(bot *gomirai.Bot, event Event, configHitList []Config, sendedGoMiraiMessageList []gomirai.Message) {
	hasReceiver := false
	groupIDSet := make(map[int64]string)
	friendIDSet := make(map[int64]string)
	for _, config := range configHitList {
		for _, receiver := range config.Receiver {
			hasReceiver = true
			switch receiver.Type {
			case SenderTypeGroup:
				if groupIDSet[receiver.ID] == "" {
					bot.SendGroupMessage(receiver.ID, 0, sendedGoMiraiMessageList)
					groupIDSet[receiver.ID] = "sent"
				}
			case SenderTypeUser:
				if friendIDSet[receiver.ID] == "" {
					bot.SendGroupMessage(receiver.ID, 0, sendedGoMiraiMessageList)
					friendIDSet[receiver.ID] = "sent"
				}
			}
		}
	}
	if hasReceiver {
		return
	}

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
