package eltype

import (
	"fmt"
	"gomirai"
	"strconv"
)

type EventType int

const (
	EventTypeGroupMessage EventType = iota
	EventTypeFriendMessage
	EventTypeMemberMute
	EventTypeGroupMuteAll
	EventTypeGroupUnMuteAll
	EventTypeMemberUnmute
)

type Event struct {
	Type          EventType
	MessageID     int64
	SenderList    []Sender
	MessageList   []Message
	OperationList []Operation
	// MuteDurationSecond int64
}

func NewEventFromGoMiraiEvent(goMiraiEvent gomirai.InEvent) (Event, error) {
	var event Event
	var operation Operation
	switch CastGoMiraiEventTypeToEventType(goMiraiEvent.Type) {
	case EventTypeGroupMessage:
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

	case EventTypeFriendMessage:
		event.Type = EventTypeFriendMessage
		sender, err := NewSender(SenderTypeFriend, goMiraiEvent.SenderFriend.ID,
			goMiraiEvent.SenderFriend.NickName, goMiraiEvent.SenderFriend.Remark)
		if err != nil {

		}
		event.SenderList = append(event.SenderList, sender)
	// case "GroupMuteAllEvent":
	case EventTypeMemberMute:
		event.Type = EventTypeMemberMute
		sender, err := NewSender(SenderTypeGroup, goMiraiEvent.OperatorGroup.Group.ID,
			goMiraiEvent.OperatorGroup.Group.Name, goMiraiEvent.OperatorGroup.Group.Permission)
		if err != nil {

		}
		event.SenderList = append(event.SenderList, sender)

		sender, err = NewSender(SenderTypeMember, goMiraiEvent.OperatorGroup.ID,
			goMiraiEvent.OperatorGroup.MemberName, goMiraiEvent.OperatorGroup.Permission)
		if err != nil {

		}
		event.SenderList = append(event.SenderList, sender)

		value := make(map[string]string)
		value["muteDurationSecond"] = strconv.FormatInt(goMiraiEvent.DurationSeconds, 10)
		operation, err = NewOperation(OperationTypeMemberMute, value)
		if err != nil {

		}
		event.OperationList = append(event.OperationList, operation)
	case EventTypeMemberUnmute:
		event.Type = EventTypeMemberUnmute
		sender, err := NewSender(SenderTypeGroup, goMiraiEvent.OperatorGroup.Group.ID,
			goMiraiEvent.OperatorGroup.Group.Name, goMiraiEvent.OperatorGroup.Group.Permission)
		if err != nil {

		}
		event.SenderList = append(event.SenderList, sender)

		sender, err = NewSender(SenderTypeMember, goMiraiEvent.OperatorGroup.ID,
			goMiraiEvent.OperatorGroup.MemberName, goMiraiEvent.OperatorGroup.Permission)
		if err != nil {

		}
		event.SenderList = append(event.SenderList, sender)

		operation, err = NewOperation(OperationTypeMemberUnmute, make(map[string]string))
		if err != nil {

		}
		event.OperationList = append(event.OperationList, operation)
	case EventTypeGroupMuteAll:
		if goMiraiEvent.Origin.(bool) {
			event.Type = EventTypeGroupUnMuteAll
		} else {
			event.Type = EventTypeGroupMuteAll
		}
		sender, err := NewSender(SenderTypeGroup, goMiraiEvent.OperatorGroup.Group.ID,
			goMiraiEvent.OperatorGroup.Group.Name, goMiraiEvent.OperatorGroup.Group.Permission)
		if err != nil {

		}
		event.SenderList = append(event.SenderList, sender)

		sender, err = NewSender(SenderTypeMember, goMiraiEvent.OperatorGroup.ID,
			goMiraiEvent.OperatorGroup.MemberName, goMiraiEvent.OperatorGroup.Permission)
		if err != nil {

		}
		event.SenderList = append(event.SenderList, sender)

		if event.Type == EventTypeGroupMuteAll {
			operation, err = NewOperation(OperationTypeGroupMuteAll, make(map[string]string))
		} else {
			operation, err = NewOperation(OperationTypeGroupUnMuteAll, make(map[string]string))
		}
		if err != nil {

		}
		event.OperationList = append(event.OperationList, operation)

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

func CastGoMiraiEventTypeToEventType(goMiaraiEventType string) EventType {
	switch goMiaraiEventType {
	case "GroupMessage":
		return EventTypeGroupMessage
	case "FriendMessage":
		return EventTypeFriendMessage
	case "GroupMuteAllEvent":
		return EventTypeGroupMuteAll
	case "MemberMuteEvent":
		return EventTypeMemberMute
	case "MemberUnmuteEvent":
		return EventTypeMemberUnmute
	default:
		panic("")
	}
}
