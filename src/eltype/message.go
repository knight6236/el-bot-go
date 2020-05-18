package eltype

import (
	"errors"
	"fmt"
	"gomirai"
)

type MessageType int

const (
	MessageTypePlain MessageType = iota
	MessageTypeImage
	MessageTypeFace
	MessageTypeEvent
)

type Message struct {
	Type  MessageType
	Value map[string]string
}

func NewMessage(messageType MessageType, value map[string]string) (Message, error) {
	var message Message
	message.Type = messageType
	message.Value = value
	return message, nil
}

func NewMessageFromGoMiraiMessage(goMiraiMessage gomirai.Message) (Message, error) {
	var message Message
	message.Value = make(map[string]string)
	switch goMiraiMessage.Type {
	case "Plain":
		message.Type = MessageTypePlain
		message.Value["text"] = goMiraiMessage.Text
	default:
		return message, fmt.Errorf("%s 是不受支持的消息类型", goMiraiMessage.Type)
	}
	return message, nil
}

func (message *Message) ToGoMiraiMessage() (gomirai.Message, error) {
	var goMiraiMessage gomirai.Message
	switch message.Type {
	case MessageTypePlain:
		goMiraiMessage.Type = "Plain"
		goMiraiMessage.Text = message.Value["text"]
	default:
		return goMiraiMessage, errors.New("暂不支持的消息类型")
	}
	return goMiraiMessage, nil
}
