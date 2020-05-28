package eltype

import (
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

// NewOperation 构造一个 Operation
// @param	operationType	OperationType		操作/事件类型
// @param	value	map[string]string	与事件/操作相关的信息
// func NewOperation(operationType OperationType, value map[string]string) (Operation, error) {
// 	var operation Operation
// 	operation.Type = operationType
// 	operation.Value = value
// 	return operation, nil
// }

// CastConfigOperationTypeToOperationType 将 ConfigOperationType 转换为 OperationType
func CastConfigOperationTypeToOperationType(configEventType string) OperationType {
	switch configEventType {
	case "MemberMute":
		return OperationTypeMemberMute
	case "MemberUnmute":
		return OperationTypeMemberUnMute
	case "GroupMuteAll":
		return OperationTypeGroupMuteAll
	case "GroupUnMuteAll":
		return OperationTypeGroupUnMuteAll
	case "MemberJoin":
		return OperationTypeMemberJoin
	case "MemberLeaveByKick":
		return OperationTypeMemberLeaveByKick
	case "MemberLeaveByQuit":
		return OperationTypeMemberLeaveByQuit
	default:
		panic("")
	}
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

func (operation *Operation) Complete(preDefVarMap map[string]string) {
	for key, value := range preDefVarMap {
		varName := fmt.Sprintf("{%s}", key)
		operation.GroupID = strings.ReplaceAll(operation.GroupID, varName, value)
		operation.UserID = strings.ReplaceAll(operation.UserID, varName, value)
		operation.Second = strings.ReplaceAll(operation.Second, varName, value)
	}
}

func (operation *Operation) Init() {
	switch operation.innerType {
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
		operation.Type = "MemberLeaveEventKick"
	case OperationTypeMemberLeaveByQuit:
		operation.Type = "MemberLeaveEventQuit"
	}
	switch operation.Type {
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
	case "MemberLeaveEventKick":
		operation.innerType = OperationTypeMemberLeaveByKick
	case "MemberLeaveEventQuit":
		operation.innerType = OperationTypeMemberLeaveByQuit
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
