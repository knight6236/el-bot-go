package eltype

import (
	"el-bot-go/src/gomirai"
	"fmt"
	"strings"
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
	innerType    OperationType
	Type         string `yaml:"type"`
	GroupID      string `yaml:"groupID"`
	GroupName    string
	OperatorID   int64
	OperatorName string
	UserID       string `yaml:"userID"`
	UserName     string
	Second       string `yaml:"second"`
}

func (operation *Operation) ToGoMiraiMessage() (gomirai.Message, bool) {
	var goMiraimessage gomirai.Message
	switch operation.innerType {
	case OperationTypeAt:
		goMiraimessage.Type = "At"
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

func (operation *Operation) CompleteContent(preDefVarMap map[string]string) {
	for key, value := range preDefVarMap {
		varName := fmt.Sprintf("{%s}", key)
		operation.GroupID = strings.ReplaceAll(operation.GroupID, varName, value)
		operation.UserID = strings.ReplaceAll(operation.UserID, varName, value)
		operation.Second = strings.ReplaceAll(operation.Second, varName, value)
	}
}

func (operation *Operation) CompleteType() {
	if operation.Type != "" {
		switch operation.Type {
		case "At":
			operation.innerType = OperationTypeAt
		case "AtAll":
			operation.innerType = OperationTypeAtAll
		case "MemberMute":
			operation.innerType = OperationTypeMemberMute
		case "MemberUnMute":
			operation.innerType = OperationTypeMemberUnMute
		case "GroupMuteAll":
			operation.innerType = OperationTypeGroupMuteAll
		case "GroupUnMuteAll":
			operation.innerType = OperationTypeGroupUnMuteAll
		case "MemberJoin":
			operation.innerType = OperationTypeMemberJoin
		case "MemberLeaveByKick":
			operation.innerType = OperationTypeMemberLeaveByKick
		case "MemberLeaveByQuit":
			operation.innerType = OperationTypeMemberLeaveByQuit
		}
	}

	switch operation.innerType {
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
		innerType:    operation.innerType,
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
