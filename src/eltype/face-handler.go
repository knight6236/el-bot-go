package eltype

type FaceHandler struct {
	configList    []Config
	messageList   []Message
	operationList []Operation
	configHitList []Config
	preDefVarMap  map[string]string
}

func NewFaceHandler(configList []Config, messageList []Message, operationList []Operation,
	preDefVarMap map[string]string) (IHandler, error) {
	var handler FaceHandler
	handler.configList = configList
	handler.messageList = messageList
	handler.preDefVarMap = preDefVarMap
	handler.searchHitConfig()
	return handler, nil
}

func (handler *FaceHandler) searchHitConfig() {

}

func (handler FaceHandler) GetConfigHitList() []Config {
	return handler.configHitList
}
