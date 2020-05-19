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
	EventTypeMemberJoin
	EventTypeMemberLeaveByKick
	EventTypeMemberLeaveByQuit
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
	var err error
	switch CastGoMiraiEventTypeToEventType(goMiraiEvent.Type) {
	case EventTypeGroupMessage:
		event, err = getGroupMessageEventTemplate(goMiraiEvent)
	case EventTypeFriendMessage:
		event, err = getFriendMessageEventTemplate(goMiraiEvent)
	case EventTypeMemberMute:
		event, err = getMemberMuteEventTemplate(goMiraiEvent)
	case EventTypeMemberUnmute:
		event, err = getMemberUnmuteEventTemplate(goMiraiEvent)
	case EventTypeGroupMuteAll:
		event, err = getGroupMuteAllEventTemplate(goMiraiEvent)
	case EventTypeMemberJoin:
		event, err = getMemberJoinEventTemplate(goMiraiEvent)
	case EventTypeMemberLeaveByKick:
		event, err = getMemberLeaveByKickEventTemplate(goMiraiEvent)
	case EventTypeMemberLeaveByQuit:
		event, err = getMemberLeaveByQuitEventTemplate(goMiraiEvent)
	default:
		return event, fmt.Errorf("%s 是不受支持的事件类型\n", goMiraiEvent.Type)
	}

	if err != nil {
		return event, err
	}

	event.parseGoMiraiMessageListToMessageList(goMiraiEvent)

	return event, nil
}

func getGroupMessageEventTemplate(goMiraiEvent gomirai.InEvent) (Event, error) {
	var event Event
	var sender Sender
	var err error
	event.Type = EventTypeGroupMessage
	sender, err = NewSender(SenderTypeGroup, goMiraiEvent.SenderGroup.Group.ID,
		goMiraiEvent.SenderGroup.Group.Name, goMiraiEvent.SenderGroup.Group.Permission)
	if err != nil {
		return event, err
	}
	event.SenderList = append(event.SenderList, sender)

	sender, err = NewSender(SenderTypeMember, goMiraiEvent.SenderGroup.ID,
		goMiraiEvent.SenderGroup.MemberName, goMiraiEvent.SenderGroup.Permission)
	if err != nil {
		return event, err
	}
	event.SenderList = append(event.SenderList, sender)
	return event, nil
}

func getFriendMessageEventTemplate(goMiraiEvent gomirai.InEvent) (Event, error) {
	var event Event
	var sender Sender
	var err error
	sender, err = NewSender(SenderTypeMember, goMiraiEvent.SenderGroup.ID,
		goMiraiEvent.SenderGroup.MemberName, goMiraiEvent.SenderGroup.Permission)
	if err != nil {
		return event, err
	}
	event.SenderList = append(event.SenderList, sender)
	return event, nil
}

func getMemberMuteEventTemplate(goMiraiEvent gomirai.InEvent) (Event, error) {
	var event Event
	var sender Sender
	var operation Operation
	var err error
	event.Type = EventTypeMemberMute
	sender, err = NewSender(SenderTypeGroup, goMiraiEvent.OperatorGroup.Group.ID,
		goMiraiEvent.OperatorGroup.Group.Name, goMiraiEvent.OperatorGroup.Group.Permission)
	if err != nil {
		return event, err
	}
	event.SenderList = append(event.SenderList, sender)

	sender, err = NewSender(SenderTypeMember, goMiraiEvent.OperatorGroup.ID,
		goMiraiEvent.OperatorGroup.MemberName, goMiraiEvent.OperatorGroup.Permission)
	if err != nil {
		return event, err
	}
	event.SenderList = append(event.SenderList, sender)

	value := make(map[string]string)
	value["muteDurationSecond"] = strconv.FormatInt(goMiraiEvent.DurationSeconds, 10)
	operation, err = NewOperation(OperationTypeMemberMute, value)
	if err != nil {
		return event, err
	}
	event.OperationList = append(event.OperationList, operation)
	return event, nil
}

func getMemberUnmuteEventTemplate(goMiraiEvent gomirai.InEvent) (Event, error) {
	var event Event
	var sender Sender
	var operation Operation
	var err error
	event.Type = EventTypeMemberUnmute
	sender, err = NewSender(SenderTypeGroup, goMiraiEvent.OperatorGroup.Group.ID,
		goMiraiEvent.OperatorGroup.Group.Name, goMiraiEvent.OperatorGroup.Group.Permission)
	if err != nil {
		return event, err
	}
	event.SenderList = append(event.SenderList, sender)

	sender, err = NewSender(SenderTypeMember, goMiraiEvent.OperatorGroup.ID,
		goMiraiEvent.OperatorGroup.MemberName, goMiraiEvent.OperatorGroup.Permission)
	if err != nil {
		return event, err
	}
	event.SenderList = append(event.SenderList, sender)

	operation, err = NewOperation(OperationTypeMemberUnmute, make(map[string]string))
	if err != nil {
		return event, err
	}
	event.OperationList = append(event.OperationList, operation)
	return event, nil
}

func getGroupMuteAllEventTemplate(goMiraiEvent gomirai.InEvent) (Event, error) {
	var event Event
	var sender Sender
	var operation Operation
	var err error
	if goMiraiEvent.Origin.(bool) {
		event.Type = EventTypeGroupUnMuteAll
	} else {
		event.Type = EventTypeGroupMuteAll
	}
	sender, err = NewSender(SenderTypeGroup, goMiraiEvent.OperatorGroup.Group.ID,
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
	return event, nil
}

func getMemberJoinEventTemplate(goMiraiEvent gomirai.InEvent) (Event, error) {
	var event Event
	var sender Sender
	var operation Operation
	var err error
	event.Type = EventTypeMemberJoin
	sender, err = NewSender(SenderTypeGroup, goMiraiEvent.Member.Group.ID,
		goMiraiEvent.Member.Group.Name, goMiraiEvent.Member.Group.Permission)
	if err != nil {
		return event, err
	}
	event.SenderList = append(event.SenderList, sender)

	sender, err = NewSender(SenderTypeMember, goMiraiEvent.Member.ID,
		goMiraiEvent.Member.MemberName, goMiraiEvent.Member.Permission)
	if err != nil {
		return event, err
	}
	event.SenderList = append(event.SenderList, sender)

	operation, err = NewOperation(OperationTypeMemberJoin, make(map[string]string))

	if err != nil {
		return event, err
	}
	event.OperationList = append(event.OperationList, operation)
	return event, nil
}

func getMemberLeaveByKickEventTemplate(goMiraiEvent gomirai.InEvent) (Event, error) {
	var event Event
	var sender Sender
	var operation Operation
	var err error
	event.Type = EventTypeMemberLeaveByKick
	sender, err = NewSender(SenderTypeGroup, goMiraiEvent.Member.Group.ID,
		goMiraiEvent.Member.Group.Name, goMiraiEvent.Member.Group.Permission)
	if err != nil {
		return event, err
	}
	event.SenderList = append(event.SenderList, sender)

	sender, err = NewSender(SenderTypeMember, goMiraiEvent.Member.ID,
		goMiraiEvent.Member.MemberName, goMiraiEvent.Member.Permission)
	if err != nil {
		return event, err
	}
	event.SenderList = append(event.SenderList, sender)

	operation, err = NewOperation(OperationTypeMemberLeaveByKick, make(map[string]string))

	if err != nil {
		return event, err
	}
	event.OperationList = append(event.OperationList, operation)
	return event, nil
}

func getMemberLeaveByQuitEventTemplate(goMiraiEvent gomirai.InEvent) (Event, error) {
	var event Event
	var sender Sender
	var operation Operation
	var err error
	event.Type = EventTypeMemberLeaveByQuit
	sender, err = NewSender(SenderTypeGroup, goMiraiEvent.Member.Group.ID,
		goMiraiEvent.Member.Group.Name, goMiraiEvent.Member.Group.Permission)
	if err != nil {
		return event, err
	}
	event.SenderList = append(event.SenderList, sender)

	sender, err = NewSender(SenderTypeMember, goMiraiEvent.Member.ID,
		goMiraiEvent.Member.MemberName, goMiraiEvent.Member.Permission)
	if err != nil {
		return event, err
	}
	event.SenderList = append(event.SenderList, sender)

	operation, err = NewOperation(OperationTypeMemberLeaveByQuit, make(map[string]string))

	if err != nil {
		return event, err
	}
	event.OperationList = append(event.OperationList, operation)
	return event, nil
}

func (event *Event) parseGoMiraiMessageListToMessageList(goMiraiEvent gomirai.InEvent) {
	for _, goMiraiMessage := range goMiraiEvent.MessageChain {
		message, err := NewMessageFromGoMiraiMessage(goMiraiMessage)
		if err != nil {
			continue
		}
		event.MessageList = append(event.MessageList, message)
	}
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
	case "MemberJoinEvent":
		return EventTypeMemberJoin
	case "MemberLeaveEventKick":
		return EventTypeMemberLeaveByKick
	case "MemberLeaveEventQuit":
		return EventTypeMemberLeaveByQuit
	default:
		panic("")
	}
}
