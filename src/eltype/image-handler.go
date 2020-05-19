package eltype

type ImageHandler struct {
	configList    []Config
	messageList   []Message
	configHitList []Config
	operationList []Operation
	preDefVarMap  map[string]string
}

func NewImageHandler(configList []Config, messageList []Message, operationList []Operation,
	preDefVarMap map[string]string) (IHandler, error) {
	var handler ImageHandler
	handler.configList = configList
	handler.messageList = messageList
	handler.preDefVarMap = preDefVarMap
	handler.searchHitConfig()
	return handler, nil
}

func (handler *ImageHandler) searchHitConfig() {

}

func (handler ImageHandler) GetConfigHitList() []Config {
	return handler.configHitList
}
