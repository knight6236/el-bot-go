package eltype

// OperationType 操作/事件类型
type OperationType int

const (
	// OperationTypeMemberMute 禁言群成员操作类型
	OperationTypeMemberMute OperationType = iota
	// OperationTypeMemberUnmute 解除群成员禁言操作类型
	OperationTypeMemberUnmute
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
		return OperationTypeMemberUnmute
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
