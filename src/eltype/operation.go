package eltype

type OperationType int

const (
	OperationTypeMemberMute OperationType = iota
	OperationTypeMemberUnmute
	OperationTypeGroupMuteAll
)

type Operation struct {
	Type  OperationType
	Value map[string]string
}

func NewOperation(operationType OperationType, value map[string]string) (Operation, error) {
	var operation Operation
	operation.Type = operationType
	operation.Value = value
	return operation, nil
}

func CastConfigOperationTypeToOperationType(configEventType string) OperationType {
	switch configEventType {
	case "MuteMember":
		return OperationTypeMemberMute
	case "MemberUnmute":
		return OperationTypeMemberUnmute
	case "GroupMuteAll":
		return OperationTypeGroupMuteAll
	default:
		panic("")
	}
}
