package eltype

import (
	"gomirai"
)

type Controller struct {
	configReader ConfigReader
}

var handlerConstructor = [...]func(configList []Config, messageList []Message) (IHandler, error){
	NewPlainHandler, NewImageHandler, NewEventHandler, NewFaceHandler}

var doerConstructor = [...]func(configHitList []Config, recivedMessageList []Message) (IDoer, error){
	NewPlainDoer, NewImageDoer, NewEventDoer, NewFaceDoer}

func NewController(configReader ConfigReader) Controller {
	var controller Controller
	controller.configReader = configReader
	return controller
}

func (controller *Controller) Commit(bot *gomirai.Bot, goMiraiEvent gomirai.InEvent) {
	err := goMiraiEvent.SenderDetail()
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
	var configList []Config
	switch event.Type {
	case EventTypeGroupMessage:
		configList = controller.mergeList(configList, controller.configReader.GlobalConfigList, controller.configReader.GroupConfigList)
	case EventTypeFriendMessage:
		configList = controller.mergeList(configList, controller.configReader.GlobalConfigList, controller.configReader.FriendConifgList)
	}

	for i := 0; i < len(handlerConstructor); i++ {
		handler, err := (handlerConstructor[i](configList, event.MessageList))
		if err != nil {
			continue
		}
		doer, err := (doerConstructor[i](handler.GetConfigHitList(), event.MessageList))
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
	case EventTypeFriendMessage:
		bot.SendFriendMessage(event.SenderList[0].ID, 0, sendedGoMiraiMessageList)
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
