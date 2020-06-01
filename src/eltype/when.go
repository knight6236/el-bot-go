package eltype

type When struct {
	Message       Message     `yaml:"message"`
	OperationList []Operation `yaml:"operation"`
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
