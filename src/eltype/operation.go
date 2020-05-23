package eltype

import (
	"el-bot-go/src/gomirai"
	"errors"
	"fmt"
	"strconv"
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
	// OperationTypeAt At某人
	OperationTypeAt
	// OperationTypeAtAll At全体成员
	OperationTypeAtAll
)

// Operation 操作/事件
// @property	Type	OperationType		操作/事件类型
// @property	Value	map[string]string	与事件/操作相关的信息
type Operation struct {
	Type  OperationType
	Value map[string]string
}

// NewOperation 构造一个 Operation
// @param	operationType	OperationType		操作/事件类型
// @param	value	map[string]string	与事件/操作相关的信息
func NewOperation(operationType OperationType, value map[string]string) (Operation, error) {
	var operation Operation
	operation.Type = operationType
	operation.Value = value
	return operation, nil
}

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
	case "At":
		return OperationTypeAt
	case "AtAll":
		return OperationTypeAtAll
	default:
		panic("")
	}
}

func (operation *Operation) ToGoMiraiMessage() (gomirai.Message, error) {
	var goMiraiMessage gomirai.Message
	switch operation.Type {
	case OperationTypeAt:
		goMiraiMessage.Type = "At"
		id, err := strconv.ParseInt(operation.Value["id"], 10, 64)
		if err != nil {
			return goMiraiMessage, err
		}
		goMiraiMessage.Target = id
	case OperationTypeAtAll:
		goMiraiMessage.Type = "AtAll"
	default:
		return goMiraiMessage, errors.New("不受支持的 Operation 类型")
	}
	return goMiraiMessage, nil
}

// ToString ...
func (operation Operation) ToString() string {
	switch operation.Type {
	case OperationTypeMemberJoin:
		return fmt.Sprintf("Operation: {Type: %s, Member: {id: %s, name: %s}}",
			"MemberJoin", operation.Value["id"], operation.Value["name"])
	case OperationTypeMemberMute:
		return fmt.Sprintf("Operation: {Type: %s, Member: {id: %s, name: %s}, Second: %s}",
			"MemberMute", operation.Value["id"], operation.Value["name"], operation.Value["second"])
	case OperationTypeMemberUnMute:
		return fmt.Sprintf("Operation: {Type: %s, Member: {id: %s, name: %s}}",
			"MemberUnMute", operation.Value["id"], operation.Value["name"])
	default:
		// TODO
		return ""
	}
}
