package eltype

type Do struct {
	IsCount       bool        `yaml:"isCount"`
	Message       Message     `yaml:"message"`
	OperationList []Operation `yaml:"operation"`
}

func (do *Do) AddOperation(operation Operation) {
	do.OperationList = append(do.OperationList, operation)
}

func (do *Do) DeepCopy() Do {
	var operationList []Operation
	for _, operaiton := range do.OperationList {
		operationList = append(operationList, operaiton.DeepCopy())
	}
	newDo := Do{
		IsCount:       do.IsCount,
		Message:       do.Message.DeepCopy(),
		OperationList: operationList,
	}
	return newDo
}
