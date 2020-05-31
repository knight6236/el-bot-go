package eltype

import (
	"fmt"
	"strings"

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
	MessageTypeAt
	MessageTypeAtAll
)

type Message struct {
	Sender     Sender          `yaml:"sender"`
	Receiver   Receiver        `yaml:"receiver"`
	DetailList []MessageDetail `yaml:"detail"`
}

type MessageDetail struct {
	innerType MessageType
	UserID    string `yaml:"userID"`
	GroupID   string `yaml:"groupID"`
	Type      string `yaml:"type"`
	Text      string `yaml:"text"`
	Regex     string `yaml:"regex"`
	URL       string `yaml:"url"`
	JSON      bool   `yaml:"json"`
	Path      string `yaml:"path"`
	ReDirect  bool   `yaml:"reDirect"`
	FaceID    int64
	FaceName  string `yaml:"faceName"`
}

// NewMessageFromGoMiraiMessage 从 gomirai.Message 中 构造一个 Message
// @param	goMiraiMessage	gomirai.Message
func NewMessageFromGoMiraiMessage(goMiraiEvent gomirai.InEvent, goMiraiMessage gomirai.Message) (Message, error) {
	var message Message
	var messageDetail MessageDetail
	switch goMiraiMessage.Type {
	case "Image":
		messageDetail.innerType = MessageTypeImage
		messageDetail.URL = goMiraiMessage.URL
	case "Plain":
		messageDetail.innerType = MessageTypePlain
		messageDetail.Text = goMiraiMessage.Text
	case "Face":
		messageDetail.innerType = MessageTypeFace
		messageDetail.FaceID = goMiraiMessage.FaceID
		messageDetail.FaceName = goMiraiMessage.Name
	case "Xml":
		messageDetail.innerType = MessageTypeXML
		messageDetail.Text = goMiraiMessage.XML
	case "At":
		messageDetail.innerType = MessageTypeAt
		messageDetail.UserID = CastInt64ToString(goMiraiMessage.Target)
	case "AtAll":
		messageDetail.innerType = MessageTypeAtAll
		messageDetail.GroupID = CastInt64ToString(goMiraiEvent.SenderGroup.Group.ID)
	default:
		return message, fmt.Errorf("%s 是不受支持的消息类型", goMiraiMessage.Type)
	}
	message.AddDetail(messageDetail)
	return message, nil
}

func (message *Message) DeepCopy() Message {
	newMessage := Message{
		Sender:   message.Sender.DeepCopy(),
		Receiver: message.Receiver.DeepCopy(),
	}
	for _, detail := range message.DetailList {
		newMessage.DetailList = append(newMessage.DetailList, detail.DeepCopy())
	}
	return newMessage
}

func (detail *MessageDetail) DeepCopy() MessageDetail {
	return MessageDetail{
		innerType: detail.innerType,
		Type:      detail.Type,
		Text:      detail.Text,
		Regex:     detail.Regex,
		URL:       detail.URL,
		JSON:      detail.JSON,
		Path:      detail.Path,
		ReDirect:  detail.ReDirect,
		FaceID:    detail.FaceID,
		FaceName:  detail.FaceName,
	}
}

func (message *Message) AddDetail(detail MessageDetail) {
	message.DetailList = append(message.DetailList, detail)
}

// ToGoMiraiMessage 将 Message 转换为 gomirai.Message
func (message *Message) ToGoMiraiMessageList() ([]gomirai.Message, bool) {
	message.CompleteType()
	var goMiraiMessageList []gomirai.Message
	for _, detail := range message.DetailList {
		goMiaraiMessage, isSuccess := detail.ToGoMiraiMessage()
		if isSuccess {
			goMiraiMessageList = append(goMiraiMessageList, goMiaraiMessage)
		}
	}
	return goMiraiMessageList, true
}

func (detail *MessageDetail) ToGoMiraiMessage() (gomirai.Message, bool) {
	detail.CompleteType()
	var goMiraiMessage gomirai.Message
	switch detail.innerType {
	case MessageTypePlain:
		goMiraiMessage.Type = "Plain"
		if detail.Text == "" {
			return goMiraiMessage, false
		}
		goMiraiMessage.Text = detail.Text
	case MessageTypeFace:
		goMiraiMessage.Type = "Face"
		goMiraiMessage.FaceID = detail.FaceID
		goMiraiMessage.Name = detail.FaceName
	case MessageTypeImage:
		goMiraiMessage.Type = "Image"
		if detail.Path == "" && detail.URL == "" {
			return goMiraiMessage, false
		}
		goMiraiMessage.Path = detail.Path
		goMiraiMessage.URL = detail.URL
	case MessageTypeXML:
		goMiraiMessage.Type = "Xml"
		if detail.Text == "" {
			return goMiraiMessage, false
		}
		goMiraiMessage.XML = detail.Text
	case MessageTypeAt:
		goMiraiMessage.Type = "At"
		goMiraiMessage.Target = CastStringToInt64(detail.UserID)
	case MessageTypeAtAll:
		goMiraiMessage.Type = "AtAll"
	default:
		return goMiraiMessage, false
	}
	return goMiraiMessage, true
}

func (message *Message) CompleteType() {
	for i := 0; i < len(message.DetailList); i++ {
		temp := message.DetailList[i]
		temp.CompleteType()
		message.DetailList[i] = temp
	}
}

func (message *Message) CompleteContent(event Event) {
	message.Sender.CompleteContent(event.PreDefVarMap)
	message.Receiver.CompleteContent(event)
	for i := 0; i < len(message.DetailList); i++ {
		temp := message.DetailList[i]
		temp.CompleteContent(event.PreDefVarMap)
		message.DetailList[i] = temp
	}
}

func (detail *MessageDetail) CompleteType() {
	switch detail.Type {
	case "Plain":
		detail.innerType = MessageTypePlain
	case "Image":
		detail.innerType = MessageTypeImage
	case "Face":
		detail.innerType = MessageTypeFace
	case "Xml":
		detail.innerType = MessageTypeXML
	case "At":
		detail.innerType = MessageTypeAt
	case "AtAll":
		detail.innerType = MessageTypeAtAll
	}
	switch detail.innerType {
	case MessageTypePlain:
		detail.Type = "Plain"
	case MessageTypeImage:
		detail.Type = "Image"
	case MessageTypeFace:
		detail.Type = "Face"
	case MessageTypeXML:
		detail.Type = "Xml"
	case MessageTypeAt:
		detail.Type = "At"
	case MessageTypeAtAll:
		detail.Type = "AtAll"
	}
}

func (detail *MessageDetail) CompleteContent(preDefVarMap map[string]string) {
	for key, value := range preDefVarMap {
		varName := fmt.Sprintf("{%s}", key)
		detail.UserID = strings.ReplaceAll(detail.UserID, varName, value)
		detail.GroupID = strings.ReplaceAll(detail.GroupID, varName, value)
		detail.Text = strings.ReplaceAll(detail.Text, varName, value)
		detail.URL = strings.ReplaceAll(detail.URL, varName, value)
		detail.Path = strings.ReplaceAll(detail.Path, varName, value)
		detail.FaceName = strings.ReplaceAll(detail.FaceName, varName, value)
	}
}

// ToString ...
func (message Message) ToString() string {
	// switch message.Type {
	// case MessageTypePlain:
	// 	return fmt.Sprintf("Message: {Type: %s, Text: %s}", "Plain", message.Value["text"])
	// case MessageTypeFace:
	// 	return fmt.Sprintf("Message: {Type: %s, Name: %s}", "Face", message.Value["name"])
	// default:
	// 	// TODO
	// 	return ""
	// }
	return ""
}
