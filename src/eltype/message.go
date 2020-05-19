package eltype

import (
	"errors"
	"fmt"
	"gomirai"
)

// MessageType 消息类型
type MessageType int

const (
	// MessageTypePlain 文本消息类型
	MessageTypePlain MessageType = iota
	// MessageTypeImage 图片消息类型
	MessageTypeImage
	// MessageTypeFace 表情消息类型
	MessageTypeFace
	// MessageTypeEvent 事件小类型
	MessageTypeEvent
)

// Message 消息
// @property	Type	MessageType 		消息类型
// @property	Value	map[string]string 	与消息相关的属性
type Message struct {
	Type  MessageType
	Value map[string]string
}

// NewMessage 构造一个 Message
// @param	messageType	MessageType 		消息类型
// @param	value		map[string]string 	与消息相关的属性
func NewMessage(messageType MessageType, value map[string]string) (Message, error) {
	var message Message
	message.Type = messageType
	message.Value = value
	return message, nil
}

// NewMessageFromGoMiraiMessage 从 gomirai.Message 中 构造一个 Message
// @param	goMiraiMessage	gomirai.Message
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

// ToGoMiraiMessage 将 Message 转换为 gomirai.Message
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
