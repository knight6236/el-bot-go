package eltype

type OperationType int

const (
	OperationTypeMute OperationType = iota
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
