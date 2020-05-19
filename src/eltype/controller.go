package eltype

import (
	"gomirai"
)

type Controller struct {
	configReader ConfigReader
}

var handlerConstructor = [...]func(configList []Config, messageList []Message, operationList []Operation) (IHandler, error){
	NewPlainHandler, NewImageHandler, NewOperationHandler, NewFaceHandler}

var doerConstructor = [...]func(configHitList []Config, recivedMessageList []Message) (IDoer, error){
	NewPlainDoer, NewImageDoer, NewOperationDoer, NewFaceDoer}

func NewController(configReader ConfigReader) Controller {
	var controller Controller
	controller.configReader = configReader
	return controller
}

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
	// err = goMiraiEvent.OperatorDetail()
	if err != nil {
		return
	}
	event, err := NewEventFromGoMiraiEvent(goMiraiEvent)
	if err != nil {
		return
	}

	var sendedGoMiraiMessageList []gomirai.Message
	var configHitList []Config
	var configList []Config
	switch event.Type {
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
	}

	for i := 0; i < len(handlerConstructor); i++ {
		handler, err := (handlerConstructor[i](configList, event.MessageList, event.OperationList))
		if err != nil {
			continue
		}

		for _, config := range handler.GetConfigHitList() {
			configHitList = append(configHitList, config)
		}
	}

	for i := 0; i < len(doerConstructor); i++ {
		doer, err := (doerConstructor[i](configHitList, event.MessageList))
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
	}

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
