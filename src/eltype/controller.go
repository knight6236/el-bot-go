package eltype

import (
	"gomirai"
)

type Controller struct {
}

var handlerConstructor = [...]func(configList []Config, messageList []Message) (IHandler, error){
	NewPlainHandler, NewImageHandler, NewEventHandler, NewFaceHandler}

var doerConstructor = [...]func(configHitList []Config, recivedMessageList []Message) (IDoer, error){
	NewPlainDoer, NewImageDoer, NewEventDoer, NewFaceDoer}

func NewController() Controller {

	var controller Controller
	// controller.configList = configList
	return controller
}

func (controller *Controller) Commit(bot *gomirai.Bot, goMiraiEvent gomirai.InEvent, configList []Config) {
	goMiraiEvent.SenderDetail()
	event, err := NewEventFromGoMiraiEvent(goMiraiEvent)
	var sendedGoMiraiMessageList []gomirai.Message
	if err != nil {
		return
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
	}

}
