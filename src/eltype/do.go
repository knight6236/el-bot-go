package eltype

type Do struct {
	IsCount       bool        `yaml:"isCount"`
	Message       Message     `yaml:"message"`
	OperationList []Operation `yaml:"operation"`
	ControlList   []Control   `yaml:"control"`
}

func (do *Do) AddOperation(operation Operation) {
	do.OperationList = append(do.OperationList, operation)
}

func (do *Do) DeepCopy() Do {
	var operationList []Operation
	for _, operaiton := range do.OperationList {
		operationList = append(operationList, operaiton.DeepCopy())
	}

	var controlList []Control
	for _, control := range do.ControlList {
		controlList = append(controlList, control.DeepCopy())
	}

	newDo := Do{
		IsCount:       do.IsCount,
		Message:       do.Message.DeepCopy(),
		OperationList: operationList,
		ControlList:   controlList,
	}
	return newDo
}
