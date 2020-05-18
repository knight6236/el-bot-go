package eltype

import (
	"gomirai"
)

type Controller struct {
}

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

	var plainHandler PlainHandler
	plainHandler, err = NewPlainHandler(configList, event.MessageList)
	if err != nil {

	}
	var plainDoer PlainDoer
	plainDoer, err = NewPlainDoer(plainHandler.ConfigHitList, event.MessageList)
	if err != nil {

	}
	for _, message := range plainDoer.SendedMessageList {
		goMiraiMessage, err := message.ToGoMiraiMessage()
		if err != nil {
			continue
		}
		sendedGoMiraiMessageList = append(sendedGoMiraiMessageList, goMiraiMessage)
	}

	var imageHandler ImageHandler
	imageHandler, err = NewImageHandler(configList, event.MessageList)
	if err != nil {

	}
	var imageDoer ImageDoer
	imageDoer, err = NewImageDoer(imageHandler.ConfigHitList, event.MessageList)
	if err != nil {

	}
	for _, message := range imageDoer.SendedMessageList {
		goMiraiMessage, err := message.ToGoMiraiMessage()
		if err != nil {
			continue
		}
		sendedGoMiraiMessageList = append(sendedGoMiraiMessageList, goMiraiMessage)
	}

	var eventHandler EventHandler
	eventHandler, err = NewEventHandler(configList, event.MessageList)
	if err != nil {

	}
	var eventDoer EventDoer
	eventDoer, err = NewEventDoer(eventHandler.ConfigHitList, event.MessageList)
	if err != nil {

	}
	for _, message := range eventDoer.SendedMessageList {
		goMiraiMessage, err := message.ToGoMiraiMessage()
		if err != nil {
			continue
		}
		sendedGoMiraiMessageList = append(sendedGoMiraiMessageList, goMiraiMessage)
	}

	var faceHandler FaceHandler
	faceHandler, err = NewFaceHandler(configList, event.MessageList)
	if err != nil {

	}
	var faceDoer FaceDoer
	faceDoer, err = NewFaceDoer(faceHandler.ConfigHitList, event.MessageList)
	if err != nil {

	}
	for _, message := range faceDoer.SendedMessageList {
		goMiraiMessage, err := message.ToGoMiraiMessage()
		if err != nil {
			continue
		}
		sendedGoMiraiMessageList = append(sendedGoMiraiMessageList, goMiraiMessage)
	}

	switch event.Type {
	case EventTypeGroupMessage:
		bot.SendGroupMessage(event.SenderList[0].ID, 0, sendedGoMiraiMessageList)
	}

}
