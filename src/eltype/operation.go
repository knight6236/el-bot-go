package eltype

import (
	"fmt"
	"strings"

	"github.com/ADD-SP/gomirai"
)

// OperationType 操作/事件类型
type OperationType int

const (
	// OperationTypeMemberMute 禁言群成员操作类型
	OperationTypeMemberMute OperationType = iota
	// OperationTypeMemberUnMute 解除群成员禁言操作类型
	OperationTypeMemberUnMute
	// OperationTypeGroupMuteAll 全员禁言操作类型
	OperationTypeGroupMuteAll
	// OperationTypeGroupUnMuteAll 解除全员禁言操作类型
	OperationTypeGroupUnMuteAll
	// OperationTypeMemberJoin 新成员入群事件类型
	OperationTypeMemberJoin
	// OperationTypeMemberLeaveByKick 踢人事件类型
	OperationTypeMemberLeaveByKick
	// OperationTypeMemberLeaveByQuit 成员自行退运事件类型
	OperationTypeMemberLeaveByQuit
	OperationTypeAt
	OperationTypeAtAll
)

type Operation struct {
	InnerType    OperationType `json:"-"`
	Type         string        `yaml:"type" json:"type"`
	GroupID      string        `yaml:"groupID" json:"groupID"`
	GroupName    string        `json:"-"`
	OperatorID   int64         `json:"-"`
	OperatorName string        `json:"-"`
	UserID       string        `yaml:"userID" json:"userID"`
	UserName     string
	Second       string `yaml:"second" json:"second"`
}

func (operation *Operation) ToGoMiraiMessage() (gomirai.Message, bool) {
	var goMiraimessage gomirai.Message
	switch operation.InnerType {
	case OperationTypeAt:
		goMiraimessage.Type = "At"
		userID := CastStringToInt64(operation.UserID)
		if userID == 0 {
			return goMiraimessage, false
		}
		goMiraimessage.Target = CastStringToInt64(operation.UserID)
	case OperationTypeAtAll:
		goMiraimessage.Type = "AtAll"
	default:
		return goMiraimessage, false
	}
	return goMiraimessage, true
}

func (operation Operation) ToString() string {
	// switch operation.Type {
	// case OperationTypeMemberJoin:
	// 	return fmt.Sprintf("Operation: {Type: %s, Member: {id: %s, name: %s}}",
	// 		"MemberJoin", operation.Value["id"], operation.Value["name"])
	// case OperationTypeMemberMute:
	// 	return fmt.Sprintf("Operation: {Type: %s, Member: {id: %s, name: %s}, Second: %s}",
	// 		"MemberMute", operation.Value["id"], operation.Value["name"], operation.Value["second"])
	// case OperationTypeMemberUnMute:
	// 	return fmt.Sprintf("Operation: {Type: %s, Member: {id: %s, name: %s}}",
	// 		"MemberUnMute", operation.Value["id"], operation.Value["name"])
	// default:
	// 	// TODO
	// 	return ""
	// }
	return ""
}

func (operation *Operation) CompleteContent(event Event) {
	for key, value := range event.PreDefVarMap {
		varName := fmt.Sprintf("{%s}", key)
		operation.GroupID = strings.ReplaceAll(operation.GroupID, varName, value)
		operation.UserID = strings.ReplaceAll(operation.UserID, varName, value)
		operation.Second = strings.ReplaceAll(operation.Second, varName, value)
	}

	if operation.GroupID == "" {
		switch event.Type {
		case EventTypeGroupMessage:
			operation.GroupID = event.Sender.GroupIDList[0]
		}
	}
}

func (operation *Operation) CompleteType() {
	if operation.Type != "" {
		switch operation.Type {
		case "At":
			operation.InnerType = OperationTypeAt
		case "AtAll":
			operation.InnerType = OperationTypeAtAll
		case "MemberMute":
			operation.InnerType = OperationTypeMemberMute
		case "MemberUnMute":
			operation.InnerType = OperationTypeMemberUnMute
		case "GroupMuteAll":
			operation.InnerType = OperationTypeGroupMuteAll
		case "GroupUnMuteAll":
			operation.InnerType = OperationTypeGroupUnMuteAll
		case "MemberJoin":
			operation.InnerType = OperationTypeMemberJoin
		case "MemberLeaveByKick":
			operation.InnerType = OperationTypeMemberLeaveByKick
		case "MemberLeaveByQuit":
			operation.InnerType = OperationTypeMemberLeaveByQuit
		}
	}

	switch operation.InnerType {
	case OperationTypeAt:
		operation.Type = "At"
	case OperationTypeAtAll:
		operation.Type = "AtAll"
	case OperationTypeMemberMute:
		operation.Type = "MemberMute"
	case OperationTypeMemberUnMute:
		operation.Type = "MemberUnMute"
	case OperationTypeGroupMuteAll:
		operation.Type = "GroupMuteAll"
	case OperationTypeGroupUnMuteAll:
		operation.Type = "GroupUnMuteAll"
	case OperationTypeMemberJoin:
		operation.Type = "MemberJoin"
	case OperationTypeMemberLeaveByKick:
		operation.Type = "MemberLeaveByKick"
	case OperationTypeMemberLeaveByQuit:
		operation.Type = "MemberLeaveByQuit"
	}
}

func (operation *Operation) DeepCopy() Operation {
	return Operation{
		InnerType:    operation.InnerType,
		Type:         operation.Type,
		GroupID:      operation.GroupID,
		GroupName:    operation.GroupName,
		OperatorID:   operation.OperatorID,
		OperatorName: operation.OperatorName,
		UserID:       operation.UserID,
		UserName:     operation.UserName,
		Second:       operation.Second,
	}
}
