package eltype

import (
	"fmt"
	"gomirai"
)

type EventType int

const (
	EventTypeGroupMessage EventType = iota
	EventTypeFriendMessage
)

type Event struct {
	Type        EventType
	MessageID   int64
	SenderList  []Sender
	MessageList []Message
}

func NewEventFromGoMiraiEvent(goMiraiEvent gomirai.InEvent) (Event, error) {
	var event Event
	switch goMiraiEvent.Type {
	case "GroupMessage":
		event.Type = EventTypeGroupMessage
		sender, err := NewSender(SenderTypeGroup, goMiraiEvent.SenderGroup.Group.ID,
			goMiraiEvent.SenderGroup.Group.Name, goMiraiEvent.SenderGroup.Group.Permission)
		if err != nil {

		}
		event.SenderList = append(event.SenderList, sender)

		sender, err = NewSender(SenderTypeMember, goMiraiEvent.SenderGroup.ID,
			goMiraiEvent.SenderGroup.MemberName, goMiraiEvent.SenderGroup.Permission)
		if err != nil {

		}
		event.SenderList = append(event.SenderList, sender)

	case "FriendMessage":
		event.Type = EventTypeFriendMessage
		sender, err := NewSender(SenderTypeFriend, goMiraiEvent.SenderFriend.ID,
			goMiraiEvent.SenderFriend.NickName, goMiraiEvent.SenderFriend.Remark)
		if err != nil {

		}
		event.SenderList = append(event.SenderList, sender)
	default:
		return event, fmt.Errorf("%s 是不受支持的事件类型\n", goMiraiEvent.Type)
	}

	for _, goMiraiMessage := range goMiraiEvent.MessageChain {
		message, err := NewMessageFromGoMiraiMessage(goMiraiMessage)
		if err != nil {
			continue
		}
		event.MessageList = append(event.MessageList, message)
	}

	return event, nil
}
