package eltype

type Do struct {
	IsCount       bool        `yaml:"isCount" json:"isCount"`
	Message       Message     `yaml:"message" json:"message"`
	OperationList []Operation `yaml:"operation" json:"operation"`
	ControlList   []Control   `yaml:"control" json:"control"`
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
