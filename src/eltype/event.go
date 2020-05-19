package eltype

import (
	"fmt"
	"gomirai"
	"strconv"
)

// EventType 事件类型
type EventType int

const (
	// EventTypeGroupMessage 群消息事件类型
	EventTypeGroupMessage EventType = iota
	// EventTypeFriendMessage 好友消息事件类型
	EventTypeFriendMessage
	// EventTypeMemberMute 群成员禁言事件类型
	EventTypeMemberMute
	// EventTypeGroupMuteAll 全员禁言事件类型
	EventTypeGroupMuteAll
	// EventTypeGroupUnMuteAll 解除全员禁言事件类型
	EventTypeGroupUnMuteAll
	// EventTypeMemberUnmute 解除群成员禁言事件类型
	EventTypeMemberUnmute
	// EventTypeMemberJoin 新成员入群事件类型
	EventTypeMemberJoin
	// EventTypeMemberLeaveByKick 踢人事件类型
	EventTypeMemberLeaveByKick
	// EventTypeMemberLeaveByQuit 群成员自行退群事件类型
	EventTypeMemberLeaveByQuit
)

// 「」

// Event 事件
// @property	Type			EventType			事件类型
// @property	MessageID		int64				接收到的消息ID
// @property	SenderList		[]Sender			消息发送者列表。如果为群消息则 index=0 为群信息， index=1 为发送消息的成员信息
// @property	MessageList		[]Message			接收到的消息列表
// @property	OperationList	[]Operation			接收到的事件/操作列表
// @property	PreDefVarMap	map[string]string	预定义变量 Map
type Event struct {
	Type          EventType
	MessageID     int64
	SenderList    []Sender
	MessageList   []Message
	OperationList []Operation
	PreDefVarMap  map[string]string
}

// NewEventFromGoMiraiEvent 从 gomirai.InEvent 中构造一个 Event
// @param	goMiraiEvent	gomirai.InEvent		事件
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
		return event, fmt.Errorf("%s 是不受支持的事件类型", goMiraiEvent.Type)
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

// CastGoMiraiEventTypeToEventType 将 GoMiaraiEventType 转化为 EventType
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
