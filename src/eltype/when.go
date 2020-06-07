package eltype

type When struct {
	Message       Message     `yaml:"message" json:"message"`
	OperationList []Operation `yaml:"operation" json:"operation"`
}

func (when *When) AddOperation(operation Operation) {
	when.OperationList = append(when.OperationList, operation)
}

func (when *When) DeepCopy() When {
	var operationList []Operation
	for _, operaiton := range when.OperationList {
		operationList = append(operationList, operaiton.DeepCopy())
	}
	newWhen := When{
		Message:       when.Message.DeepCopy(),
		OperationList: operationList,
	}
	return newWhen
}
