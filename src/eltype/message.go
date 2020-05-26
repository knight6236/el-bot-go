package eltype

import (
	"fmt"
	"strconv"

	"el-bot-go/src/gomirai"
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
	// MessageTypeEvent 事件消息类型
	MessageTypeEvent
	// MessageTypeXML XML消息类型
	MessageTypeXML
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
	case "Image":
		message.Type = MessageTypeImage
		message.Value["url"] = goMiraiMessage.URL
	case "Plain":
		message.Type = MessageTypePlain
		message.Value["text"] = goMiraiMessage.Text
	case "Face":
		message.Type = MessageTypeFace
		message.Value["id"] = strconv.FormatInt(goMiraiMessage.FaceID, 10)
		message.Value["name"] = goMiraiMessage.Name
	case "Xml":
		message.Type = MessageTypeXML
		message.Value["xml"] = goMiraiMessage.XML
	default:
		return message, fmt.Errorf("%s 是不受支持的消息类型", goMiraiMessage.Type)
	}
	return message, nil
}

// ToGoMiraiMessage 将 Message 转换为 gomirai.Message
func (message *Message) ToGoMiraiMessage() (gomirai.Message, bool) {
	var goMiraiMessage gomirai.Message
	var err error
	switch message.Type {
	case MessageTypePlain:
		goMiraiMessage.Type = "Plain"
		if message.Value["text"] == "" {
			return goMiraiMessage, false
		}
		goMiraiMessage.Text = message.Value["text"]
	case MessageTypeFace:
		goMiraiMessage.Type = "Face"
		goMiraiMessage.FaceID, err = strconv.ParseInt(message.Value["id"], 10, 64)
		if err != nil {
			return goMiraiMessage, false
		}
		goMiraiMessage.Name = message.Value["name"]
	case MessageTypeImage:
		goMiraiMessage.Type = "Image"
		if message.Value["path"] == "" && message.Value["url"] == "" {
			return goMiraiMessage, false
		}
		goMiraiMessage.Path = message.Value["path"]
		goMiraiMessage.URL = message.Value["url"]
	case MessageTypeXML:
		goMiraiMessage.Type = "Xml"
		goMiraiMessage.XML = message.Value["xml"]
	default:
		return goMiraiMessage, false
	}
	return goMiraiMessage, true
}

// ToString ...
func (message Message) ToString() string {
	switch message.Type {
	case MessageTypePlain:
		return fmt.Sprintf("Message: {Type: %s, Text: %s}", "Plain", message.Value["text"])
	case MessageTypeFace:
		return fmt.Sprintf("Message: {Type: %s, Name: %s}", "Face", message.Value["name"])
	default:
		// TODO
		return ""
	}
}
