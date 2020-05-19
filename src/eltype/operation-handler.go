package eltype

type OperationHandler struct {
	configList    []Config
	messageList   []Message
	operationList []Operation
	configHitList []Config
}

func NewOperationHandler(configList []Config, messageList []Message, operationList []Operation) (IHandler, error) {
	var handler OperationHandler
	handler.configList = configList
	handler.operationList = operationList
	handler.searchHitConfig()
	return handler, nil
}

func (handler *OperationHandler) searchHitConfig() {
	for _, config := range handler.configList {
		goto SECOND_LOOP
	TOP_LOOP:
		continue
	SECOND_LOOP:
		for _, operation := range handler.operationList {
			for _, doOperation := range config.WhenOperationList {
				if operation.Type == doOperation.Type {
					handler.configHitList = append(handler.configHitList, config)
					goto TOP_LOOP
				}
			}
		}
	}
}

func (handler OperationHandler) GetConfigHitList() []Config {
	return handler.configHitList
}