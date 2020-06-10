package eltype

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ADD-SP/gomirai"
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
	// MessageTypeAt At
	MessageTypeAt
	// MessageTypeAtAll AtAll
	MessageTypeAtAll
)

type Message struct {
	At         bool            `yaml:"at" json:"at"`
	IsQuote    bool            `yaml:"quote" json:"quote"`
	QuoteID    int64           `yaml:"-" json:"quoteID"`
	Sender     Sender          `yaml:"sender" json:"sender"`
	Receiver   Receiver        `yaml:"receiver" json:"receiver"`
	DetailList []MessageDetail `yaml:"detail" json:"detail"`
}

type MessageDetail struct {
	InnerType MessageType `json:"-"`
	UserID    string      `yaml:"userID" json:"userID"`
	GroupID   string      `yaml:"groupID" json:"groupID"`
	Type      string      `yaml:"type" json:"type"`
	Text      string      `yaml:"text" json:"text"`
	Regex     string      `yaml:"regex" json:"regex"`
	URL       string      `yaml:"url" json:"url"`
	JSON      bool        `yaml:"json" json:"json"`
	Path      string      `yaml:"path" json:"path"`
	ReDirect  bool        `yaml:"reDirect" json:"reDirect"`
	FaceID    int64
	FaceName  string `yaml:"faceName" json:"faceName"`
}

// NewMessageFromGoMiraiMessage 从 gomirai.Message 中 构造一个 Message
// @param	goMiraiMessage	gomirai.Message
func NewMessageFromGoMiraiMessage(goMiraiEvent gomirai.InEvent) (Message, error) {
	var message Message
	for _, goMiraiMessage := range goMiraiEvent.MessageChain {
		var messageDetail MessageDetail
		switch goMiraiMessage.Type {
		case "Image":
			messageDetail.InnerType = MessageTypeImage
			messageDetail.URL = goMiraiMessage.URL
		case "Plain":
			messageDetail.InnerType = MessageTypePlain
			messageDetail.Text = goMiraiMessage.Text
		case "Face":
			messageDetail.InnerType = MessageTypeFace
			messageDetail.FaceID = goMiraiMessage.FaceID
			messageDetail.FaceName = goMiraiMessage.Name
		case "Xml":
			messageDetail.InnerType = MessageTypeXML
			messageDetail.Text = goMiraiMessage.XML
		case "At":
			messageDetail.InnerType = MessageTypeAt
			messageDetail.UserID = CastInt64ToString(goMiraiMessage.Target)
			if QQ == goMiraiMessage.Target {
				message.At = true
			}
		case "AtAll":
			messageDetail.InnerType = MessageTypeAtAll
			messageDetail.GroupID = CastInt64ToString(goMiraiEvent.SenderGroup.Group.ID)
		default:
			continue
			// return message, fmt.Errorf("%s 是不受支持的消息类型", goMiraiMessage.Type)
		}
		message.AddDetail(messageDetail)
	}

	return message, nil
}

func (message *Message) DeepCopy() Message {
	newMessage := Message{
		At:       message.At,
		Sender:   message.Sender.DeepCopy(),
		Receiver: message.Receiver.DeepCopy(),
		IsQuote:  message.IsQuote,
	}
	for _, detail := range message.DetailList {
		newMessage.DetailList = append(newMessage.DetailList, detail.DeepCopy())
	}
	return newMessage
}

func (detail *MessageDetail) DeepCopy() MessageDetail {
	return MessageDetail{
		InnerType: detail.InnerType,
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
	switch detail.InnerType {
	case MessageTypePlain:
		goMiraiMessage.Type = "Plain"
		if detail.Text == "" {
			return goMiraiMessage, false
		}
		goMiraiMessage.Text = detail.Text
	case MessageTypeFace:
		goMiraiMessage.Type = "Face"
		goMiraiMessage.FaceID = detail.FaceID
		if detail.FaceName == "" {
			return goMiraiMessage, false
		}
		goMiraiMessage.Name = detail.FaceName
	case MessageTypeImage:
		goMiraiMessage.Type = "Image"
		isMatch, err := regexp.Match("[a-zA-z]+://[^\\s]*", []byte(detail.URL))
		if detail.Path == "" && (detail.URL == "" || !isMatch || err != nil) {

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
		userID := CastStringToInt64(detail.UserID)
		if userID == 0 {
			return goMiraiMessage, false
		}
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
		detail.InnerType = MessageTypePlain
	case "Image":
		detail.InnerType = MessageTypeImage
	case "Face":
		detail.InnerType = MessageTypeFace
	case "Xml":
		detail.InnerType = MessageTypeXML
	case "At":
		detail.InnerType = MessageTypeAt
	case "AtAll":
		detail.InnerType = MessageTypeAtAll
	}
	switch detail.InnerType {
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
