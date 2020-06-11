package eltype

import (
	"fmt"
	"strconv"

	"github.com/ADD-SP/gomirai"
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
	InnerType     EventType   `json:"-"`
	Type          string      `json:"type"`
	MessageID     int64       `json:"messageID"`
	Sender        Sender      `json:"sender"`
	Message       Message     `json:"message"`
	OperationList []Operation `json:"operation"`
	PreDefVarMap  map[string]string
}

// NewEventFromGoMiraiEvent 从 gomirai.InEvent 中构造一个 Event
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

	if len(goMiraiEvent.MessageChain) != 0 {
		event.MessageID = goMiraiEvent.MessageChain[0].ID
	}

	event.parseGoMiraiMessageListToMessageList(goMiraiEvent)

	for i := 0; i < len(event.OperationList); i++ {
		temp := event.OperationList[i]
		temp.CompleteType()
		event.OperationList[i] = temp
	}

	event.addSomePreDefVar()

	return event, nil
}

func (event *Event) makeGroupMessageEventTemplate(goMiraiEvent gomirai.InEvent) error {
	event.InnerType = EventTypeGroupMessage
	event.Sender.AddGroupID(CastInt64ToString(goMiraiEvent.SenderGroup.Group.ID))
	event.Sender.AddUserID(CastInt64ToString(goMiraiEvent.SenderGroup.ID))
	event.AddPerDefVar("el-sender-group-id", goMiraiEvent.SenderGroup.Group.ID)
	event.AddPerDefVar("el-sender-group-name", goMiraiEvent.OperatorGroup.Group.Name)
	event.AddPerDefVar("el-sender-user-id", goMiraiEvent.SenderGroup.ID)
	event.AddPerDefVar("el-sender-user-name", goMiraiEvent.SenderGroup.MemberName)
	return nil
}

func (event *Event) makeFriendMessageEventTemplate(goMiraiEvent gomirai.InEvent) error {
	event.InnerType = EventTypeFriendMessage
	event.Sender.AddUserID(CastInt64ToString(goMiraiEvent.SenderFriend.ID))
	event.AddPerDefVar("el-sender-user-id", goMiraiEvent.SenderFriend.ID)
	event.AddPerDefVar("el-sender-user-name", goMiraiEvent.SenderFriend.NickName)
	return nil
}

func (event *Event) makeMemberMuteEventTemplate(goMiraiEvent gomirai.InEvent) error {
	event.InnerType = EventTypeMemberMute
	event.Sender.AddGroupID(CastInt64ToString(goMiraiEvent.OperatorGroup.Group.ID))
	event.Sender.AddUserID(CastInt64ToString(goMiraiEvent.OperatorGroup.ID))
	event.Sender.AddUserID(CastInt64ToString(goMiraiEvent.Member.ID))

	operation := Operation{
		InnerType:    OperationTypeMemberMute,
		GroupID:      CastInt64ToString(goMiraiEvent.OperatorGroup.Group.ID),
		GroupName:    goMiraiEvent.OperatorGroup.Group.Name,
		OperatorID:   goMiraiEvent.OperatorGroup.ID,
		OperatorName: goMiraiEvent.OperatorGroup.MemberName,
		UserID:       CastInt64ToString(goMiraiEvent.Member.ID),
		UserName:     goMiraiEvent.Member.MemberName,
		Second:       CastInt64ToString(goMiraiEvent.DurationSeconds),
	}

	event.AddPerDefVar("el-sender-group-id", operation.GroupID)
	event.AddPerDefVar("el-sender-group-name", operation.GroupName)

	event.AddPerDefVar("el-operator-group-id", operation.GroupID)
	event.AddPerDefVar("el-operator-group-name", operation.GroupName)
	event.AddPerDefVar("el-operator-user-id", operation.OperatorID)
	event.AddPerDefVar("el-operator-user-name", operation.OperatorName)

	event.AddPerDefVar("el-target-group-id", operation.GroupID)
	event.AddPerDefVar("el-target-group-name", operation.GroupName)
	event.AddPerDefVar("el-target-user-id", operation.UserID)
	event.AddPerDefVar("el-target-user-name", operation.UserName)

	event.AddPerDefVar("el-mute-second-", operation.Second)

	event.OperationList = append(event.OperationList, operation)
	return nil
}

func (event *Event) makeMemberUnmuteEventTemplate(goMiraiEvent gomirai.InEvent) error {
	event.InnerType = EventTypeMemberUnmute
	event.Sender.AddGroupID(CastInt64ToString(goMiraiEvent.OperatorGroup.Group.ID))
	event.Sender.AddUserID(CastInt64ToString(goMiraiEvent.OperatorGroup.ID))
	event.Sender.AddUserID(CastInt64ToString(goMiraiEvent.Member.ID))

	operation := Operation{
		InnerType:    OperationTypeMemberUnMute,
		GroupID:      CastInt64ToString(goMiraiEvent.OperatorGroup.Group.ID),
		GroupName:    goMiraiEvent.OperatorGroup.Group.Name,
		OperatorID:   goMiraiEvent.OperatorGroup.ID,
		OperatorName: goMiraiEvent.OperatorGroup.MemberName,
		UserID:       CastInt64ToString(goMiraiEvent.Member.ID),
		UserName:     goMiraiEvent.Member.MemberName,
	}

	event.AddPerDefVar("el-sender-group-id", operation.GroupID)
	event.AddPerDefVar("el-sender-group-name", operation.GroupName)

	event.AddPerDefVar("el-operator-group-id", operation.GroupID)
	event.AddPerDefVar("el-operator-group-name", operation.GroupName)
	event.AddPerDefVar("el-operator-user-id", operation.OperatorID)
	event.AddPerDefVar("el-operator-user-name", operation.OperatorName)

	event.AddPerDefVar("el-target-group-id", operation.GroupID)
	event.AddPerDefVar("el-target-group-name", operation.GroupName)
	event.AddPerDefVar("el-target-user-id", operation.UserID)
	event.AddPerDefVar("el-target-user-name", operation.UserName)

	event.OperationList = append(event.OperationList, operation)
	return nil
}

func (event *Event) makeGroupMuteAllEventTemplate(goMiraiEvent gomirai.InEvent) error {
	var operationType OperationType
	if goMiraiEvent.Origin.(bool) {
		event.InnerType = EventTypeGroupUnMuteAll
		operationType = OperationTypeGroupUnMuteAll
	} else {
		event.InnerType = EventTypeGroupMuteAll
		operationType = OperationTypeGroupMuteAll
	}
	event.Sender.AddGroupID(CastInt64ToString(goMiraiEvent.OperatorGroup.Group.ID))
	// event.Sender.AddUserID(CastInt64ToString(goMiraiEvent.OperatorGroup.ID))
	event.Sender.AddUserID(CastInt64ToString(goMiraiEvent.Member.ID))

	operation := Operation{
		InnerType:    operationType,
		GroupID:      CastInt64ToString(goMiraiEvent.OperatorGroup.Group.ID),
		GroupName:    goMiraiEvent.OperatorGroup.Group.Name,
		OperatorID:   goMiraiEvent.OperatorGroup.ID,
		OperatorName: goMiraiEvent.OperatorGroup.MemberName,
	}

	event.AddPerDefVar("el-sender-group-id", operation.GroupID)
	event.AddPerDefVar("el-sender-group-name", operation.GroupName)

	event.AddPerDefVar("el-operator-group-id", operation.GroupID)
	event.AddPerDefVar("el-operator-group-name", operation.GroupName)
	event.AddPerDefVar("el-operator-user-id", operation.OperatorID)
	event.AddPerDefVar("el-operator-user-name", operation.OperatorName)

	event.AddPerDefVar("el-target-group-id", operation.GroupID)
	event.AddPerDefVar("el-target-group-name", operation.GroupName)

	event.OperationList = append(event.OperationList, operation)
	return nil
}

func (event *Event) makeMemberJoinEventTemplate(goMiraiEvent gomirai.InEvent) error {
	event.InnerType = EventTypeMemberJoin
	event.Sender.AddGroupID(CastInt64ToString(goMiraiEvent.Member.Group.ID))
	event.Sender.AddUserID(CastInt64ToString(goMiraiEvent.Member.ID))

	operation := Operation{
		InnerType: OperationTypeMemberJoin,
		GroupID:   CastInt64ToString(goMiraiEvent.OperatorGroup.Group.ID),
		GroupName: goMiraiEvent.OperatorGroup.Group.Name,
		UserID:    CastInt64ToString(goMiraiEvent.Member.ID),
		UserName:  goMiraiEvent.Member.MemberName,
	}

	event.AddPerDefVar("el-sender-group-id", operation.GroupID)
	event.AddPerDefVar("el-sender-group-name", operation.GroupName)

	event.AddPerDefVar("el-target-group-id", operation.GroupID)
	event.AddPerDefVar("el-target-group-name", operation.GroupName)
	event.AddPerDefVar("el-target-user-id", operation.UserID)
	event.AddPerDefVar("el-target-user-name", operation.UserName)

	event.OperationList = append(event.OperationList, operation)
	return nil
}

func (event *Event) makeMemberLeaveByKickEventTemplate(goMiraiEvent gomirai.InEvent) error {
	event.InnerType = EventTypeMemberLeaveByKick
	event.Sender.AddGroupID(CastInt64ToString(goMiraiEvent.Member.Group.ID))
	event.Sender.AddUserID(CastInt64ToString(goMiraiEvent.Member.ID))

	operation := Operation{
		InnerType:    OperationTypeMemberLeaveByKick,
		GroupID:      CastInt64ToString(goMiraiEvent.OperatorGroup.Group.ID),
		GroupName:    goMiraiEvent.OperatorGroup.Group.Name,
		OperatorID:   goMiraiEvent.OperatorGroup.ID,
		OperatorName: goMiraiEvent.OperatorGroup.MemberName,
		UserID:       CastInt64ToString(goMiraiEvent.Member.ID),
		UserName:     goMiraiEvent.Member.MemberName,
	}

	event.AddPerDefVar("el-sender-group-id", operation.GroupID)
	event.AddPerDefVar("el-sender-group-name", operation.GroupName)

	event.AddPerDefVar("el-operator-group-id", operation.GroupID)
	event.AddPerDefVar("el-operator-group-name", operation.GroupName)
	event.AddPerDefVar("el-operator-user-id", operation.OperatorID)
	event.AddPerDefVar("el-operator-user-name", operation.OperatorName)

	event.AddPerDefVar("el-target-group-id", operation.GroupID)
	event.AddPerDefVar("el-target-group-name", operation.GroupName)
	event.AddPerDefVar("el-target-user-id", operation.UserID)
	event.AddPerDefVar("el-target-user-name", operation.UserName)

	event.OperationList = append(event.OperationList, operation)
	return nil
}

func (event *Event) makeMemberLeaveByQuitEventTemplate(goMiraiEvent gomirai.InEvent) error {
	event.InnerType = EventTypeMemberLeaveByQuit
	event.Sender.AddGroupID(CastInt64ToString(goMiraiEvent.Member.Group.ID))
	event.Sender.AddUserID(CastInt64ToString(goMiraiEvent.Member.ID))

	operation := Operation{
		InnerType: OperationTypeMemberLeaveByQuit,
		GroupID:   CastInt64ToString(goMiraiEvent.OperatorGroup.Group.ID),
		GroupName: goMiraiEvent.OperatorGroup.Group.Name,
		UserID:    CastInt64ToString(goMiraiEvent.Member.ID),
		UserName:  goMiraiEvent.Member.MemberName,
	}

	event.AddPerDefVar("el-sender-group-id", operation.GroupID)
	event.AddPerDefVar("el-sender-group-name", operation.GroupName)

	event.AddPerDefVar("el-target-group-id", operation.GroupID)
	event.AddPerDefVar("el-target-group-name", operation.GroupName)
	event.AddPerDefVar("el-target-user-id", operation.UserID)
	event.AddPerDefVar("el-target-user-name", operation.UserName)

	event.OperationList = append(event.OperationList, operation)
	return nil
}

func (event *Event) parseGoMiraiMessageListToMessageList(goMiraiEvent gomirai.InEvent) {
	message, err := NewMessageFromGoMiraiMessage(goMiraiEvent)
	if err != nil {
		return
	}
	event.Message = message
}

func (event *Event) addSomePreDefVar() {
	text := ""
	xml := ""
	imageIndex := 0
	for _, detail := range event.Message.DetailList {
		if detail.InnerType == MessageTypePlain {
			text = text + detail.Text
		} else if detail.InnerType == MessageTypeXML {
			xml = xml + detail.Text
		} else if detail.InnerType == MessageTypeImage {
			event.AddPerDefVar(fmt.Sprintf("el-message-image-url-%d", imageIndex), detail.URL)
			imageIndex++
		}
	}
	event.AddPerDefVar("\\n", "\n")
	event.AddPerDefVar("el-message-text", text)
	event.AddPerDefVar("el-message-xml", xml)
}

func (event *Event) AddPerDefVar(varName string, value interface{}) {
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
func (event *Event) CompleteType() {
	event.Message.CompleteType()
	if event.Type != "" {
		switch event.Type {
		case "GroupMessage":
			event.InnerType = EventTypeGroupMessage
		case "FriendMessage":
			event.InnerType = EventTypeFriendMessage
		case "GroupMuteAll":
			event.InnerType = EventTypeGroupMuteAll
		case "GroupUnMuteAll":
			event.InnerType = EventTypeGroupUnMuteAll
		case "MemberMute":
			event.InnerType = EventTypeMemberMute
		case "MemberUnmute":
			event.InnerType = EventTypeMemberUnmute
		case "MemberJoin":
			event.InnerType = EventTypeMemberJoin
		case "MemberLeaveByKick":
			event.InnerType = EventTypeMemberLeaveByKick
		case "MemberLeaveByQuit":
			event.InnerType = EventTypeMemberLeaveByQuit
		}
	}

	switch event.InnerType {
	case EventTypeGroupMessage:
		event.Type = "GroupMessage"
	case EventTypeFriendMessage:
		event.Type = "FriendMessage"
	case EventTypeGroupMuteAll:
		event.Type = "GroupMuteAll"
	case EventTypeGroupUnMuteAll:
		event.Type = "GroupUnMuteAll"
	case EventTypeMemberMute:
		event.Type = "MemberMute"
	case EventTypeMemberUnmute:
		event.Type = "MemberUnmute"
	case EventTypeMemberJoin:
		event.Type = "MemberJoin"
	case EventTypeMemberLeaveByKick:
		event.Type = "MemberLeaveByKick"
	case EventTypeMemberLeaveByQuit:
		event.Type = "MemberLeaveByQuit"
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
