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
	PreDefVarMap  map[string]string
}

func NewEventFromGoMiraiEvent(goMiraiEvent gomirai.InEvent) (Event, error) {
	var event Event
	event.PreDefVarMap = make(map[string]string)
	var err error
	switch CastGoMiraiEventTypeToEventType(goMiraiEvent.Type) {
	case EventTypeGroupMessage:
		err = event.makeGroupMessageEventTemplate(goMiraiEvent)
	case EventTypeFriendMessage:
		err = event.makeFriendMessageEventTemplate(goMiraiEvent)
	case EventTypeMemberMute:
		err = event.makeMemberMuteEventTemplate(goMiraiEvent)
	case EventTypeMemberUnmute:
		err = event.makeMemberUnmuteEventTemplate(goMiraiEvent)
	case EventTypeGroupMuteAll:
		err = event.makeGroupMuteAllEventTemplate(goMiraiEvent)
	case EventTypeMemberJoin:
		err = event.makeMemberJoinEventTemplate(goMiraiEvent)
	case EventTypeMemberLeaveByKick:
		err = event.makeMemberLeaveByKickEventTemplate(goMiraiEvent)
	case EventTypeMemberLeaveByQuit:
		err = event.makeMemberLeaveByQuitEventTemplate(goMiraiEvent)
	default:
		return event, fmt.Errorf("%s 是不受支持的事件类型\n", goMiraiEvent.Type)
	}

	if err != nil {
		return event, err
	}

	event.parseGoMiraiMessageListToMessageList(goMiraiEvent)

	event.addSomePreDefVar()

	return event, nil
}

func (event *Event) makeGroupMessageEventTemplate(goMiraiEvent gomirai.InEvent) error {
	var sender Sender
	var err error
	event.Type = EventTypeGroupMessage
	sender, err = NewSender(SenderTypeGroup, goMiraiEvent.SenderGroup.Group.ID,
		goMiraiEvent.SenderGroup.Group.Name, goMiraiEvent.SenderGroup.Group.Permission)
	if err != nil {
		return err
	}
	event.SenderList = append(event.SenderList, sender)

	sender, err = NewSender(SenderTypeMember, goMiraiEvent.SenderGroup.ID,
		goMiraiEvent.SenderGroup.MemberName, goMiraiEvent.SenderGroup.Permission)
	if err != nil {
		return err
	}
	event.SenderList = append(event.SenderList, sender)
	event.addPerDefVar("el-sender-id", sender.ID)
	event.addPerDefVar("el-sender-name", sender.Name)
	return nil
}

func (event *Event) makeFriendMessageEventTemplate(goMiraiEvent gomirai.InEvent) error {
	var sender Sender
	var err error
	sender, err = NewSender(SenderTypeMember, goMiraiEvent.SenderGroup.ID,
		goMiraiEvent.SenderGroup.MemberName, goMiraiEvent.SenderGroup.Permission)
	if err != nil {
		return err
	}
	event.SenderList = append(event.SenderList, sender)
	event.addPerDefVar("el-sender-id", sender.ID)
	event.addPerDefVar("el-sender-name", sender.Name)
	return nil
}

func (event *Event) makeMemberMuteEventTemplate(goMiraiEvent gomirai.InEvent) error {
	var sender Sender
	var operation Operation
	var err error
	event.Type = EventTypeMemberMute
	sender, err = NewSender(SenderTypeGroup, goMiraiEvent.OperatorGroup.Group.ID,
		goMiraiEvent.OperatorGroup.Group.Name, goMiraiEvent.OperatorGroup.Group.Permission)
	if err != nil {
		return err
	}
	event.SenderList = append(event.SenderList, sender)

	sender, err = NewSender(SenderTypeMember, goMiraiEvent.OperatorGroup.ID,
		goMiraiEvent.OperatorGroup.MemberName, goMiraiEvent.OperatorGroup.Permission)
	if err != nil {
		return err
	}
	event.SenderList = append(event.SenderList, sender)

	sender, err = NewSender(SenderTypeMember, goMiraiEvent.Member.ID,
		goMiraiEvent.Member.MemberName, goMiraiEvent.Member.Permission)
	if err != nil {
		return err
	}
	event.SenderList = append(event.SenderList, sender)
	event.addPerDefVar("el-target-id", sender.ID)
	event.addPerDefVar("el-target-name", sender.Name)

	value := make(map[string]string)
	value["muteDurationSecond"] = strconv.FormatInt(goMiraiEvent.DurationSeconds, 10)
	operation, err = NewOperation(OperationTypeMemberMute, value)
	if err != nil {
		return err
	}
	event.OperationList = append(event.OperationList, operation)
	return nil
}

func (event *Event) makeMemberUnmuteEventTemplate(goMiraiEvent gomirai.InEvent) error {
	var sender Sender
	var operation Operation
	var err error
	event.Type = EventTypeMemberUnmute
	sender, err = NewSender(SenderTypeGroup, goMiraiEvent.OperatorGroup.Group.ID,
		goMiraiEvent.OperatorGroup.Group.Name, goMiraiEvent.OperatorGroup.Group.Permission)
	if err != nil {
		return err
	}
	event.SenderList = append(event.SenderList, sender)

	sender, err = NewSender(SenderTypeMember, goMiraiEvent.Member.ID,
		goMiraiEvent.Member.MemberName, goMiraiEvent.Member.Permission)
	if err != nil {
		return err
	}
	event.SenderList = append(event.SenderList, sender)
	event.addPerDefVar("el-target-id", sender.ID)
	event.addPerDefVar("el-target-name", sender.Name)

	operation, err = NewOperation(OperationTypeMemberUnmute, make(map[string]string))
	if err != nil {
		return err
	}
	event.OperationList = append(event.OperationList, operation)
	return nil
}

func (event *Event) makeGroupMuteAllEventTemplate(goMiraiEvent gomirai.InEvent) error {
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
		return err
	}
	event.SenderList = append(event.SenderList, sender)

	sender, err = NewSender(SenderTypeMember, goMiraiEvent.OperatorGroup.ID,
		goMiraiEvent.OperatorGroup.MemberName, goMiraiEvent.OperatorGroup.Permission)
	if err != nil {
		return err
	}
	event.SenderList = append(event.SenderList, sender)

	if event.Type == EventTypeGroupMuteAll {
		operation, err = NewOperation(OperationTypeGroupMuteAll, make(map[string]string))
	} else {
		operation, err = NewOperation(OperationTypeGroupUnMuteAll, make(map[string]string))
	}
	if err != nil {
		return err
	}
	event.OperationList = append(event.OperationList, operation)
	return nil
}

func (event *Event) makeMemberJoinEventTemplate(goMiraiEvent gomirai.InEvent) error {
	var sender Sender
	var operation Operation
	var err error
	event.Type = EventTypeMemberJoin
	sender, err = NewSender(SenderTypeGroup, goMiraiEvent.Member.Group.ID,
		goMiraiEvent.Member.Group.Name, goMiraiEvent.Member.Group.Permission)
	if err != nil {
		return err
	}
	event.SenderList = append(event.SenderList, sender)

	sender, err = NewSender(SenderTypeMember, goMiraiEvent.Member.ID,
		goMiraiEvent.Member.MemberName, goMiraiEvent.Member.Permission)
	if err != nil {
		return err
	}
	event.SenderList = append(event.SenderList, sender)
	event.addPerDefVar("el-target-id", sender.ID)
	event.addPerDefVar("el-target-name", sender.Name)

	operation, err = NewOperation(OperationTypeMemberJoin, make(map[string]string))

	if err != nil {
		return err
	}
	event.OperationList = append(event.OperationList, operation)
	return nil
}

func (event *Event) makeMemberLeaveByKickEventTemplate(goMiraiEvent gomirai.InEvent) error {
	var sender Sender
	var operation Operation
	var err error
	event.Type = EventTypeMemberLeaveByKick
	sender, err = NewSender(SenderTypeGroup, goMiraiEvent.Member.Group.ID,
		goMiraiEvent.Member.Group.Name, goMiraiEvent.Member.Group.Permission)
	if err != nil {
		return err
	}
	event.SenderList = append(event.SenderList, sender)

	sender, err = NewSender(SenderTypeMember, goMiraiEvent.Member.ID,
		goMiraiEvent.Member.MemberName, goMiraiEvent.Member.Permission)
	if err != nil {
		return err
	}
	event.SenderList = append(event.SenderList, sender)
	event.addPerDefVar("el-target-id", sender.ID)
	event.addPerDefVar("el-target-name", sender.Name)

	operation, err = NewOperation(OperationTypeMemberLeaveByKick, make(map[string]string))

	if err != nil {
		return err
	}
	event.OperationList = append(event.OperationList, operation)
	return nil
}

func (event *Event) makeMemberLeaveByQuitEventTemplate(goMiraiEvent gomirai.InEvent) error {
	var sender Sender
	var operation Operation
	var err error
	event.Type = EventTypeMemberLeaveByQuit
	sender, err = NewSender(SenderTypeGroup, goMiraiEvent.Member.Group.ID,
		goMiraiEvent.Member.Group.Name, goMiraiEvent.Member.Group.Permission)
	if err != nil {
		return err
	}
	event.SenderList = append(event.SenderList, sender)

	sender, err = NewSender(SenderTypeMember, goMiraiEvent.Member.ID,
		goMiraiEvent.Member.MemberName, goMiraiEvent.Member.Permission)
	if err != nil {
		return err
	}
	event.SenderList = append(event.SenderList, sender)
	event.addPerDefVar("el-target-id", sender.ID)
	event.addPerDefVar("el-target-name", sender.Name)

	operation, err = NewOperation(OperationTypeMemberLeaveByQuit, make(map[string]string))

	if err != nil {
		return err
	}
	event.OperationList = append(event.OperationList, operation)
	return nil
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

func (event *Event) addSomePreDefVar() {
	text := ""
	for _, message := range event.MessageList {
		if message.Type == MessageTypePlain {
			text = text + message.Value["text"]
		}
	}
	event.addPerDefVar("el-message", text)
}

func (event *Event) addPerDefVar(varName string, value interface{}) {
	switch value.(type) {
	case string:
		event.PreDefVarMap[varName] = value.(string)
	case int:
		event.PreDefVarMap[varName] = strconv.Itoa(value.(int))
	case int64:
		event.PreDefVarMap[varName] = strconv.FormatInt(value.(int64), 10)
	case float64:
		event.PreDefVarMap[varName] = fmt.Sprintf("%.2f", value.(float64))
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
