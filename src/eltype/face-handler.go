package eltype

type FaceHandler struct {
	configList    []Config
	messageList   []Message
	configHitList []Config
}

func NewFaceHandler(configList []Config, messageList []Message) (IHandler, error) {
	var handler FaceHandler
	handler.configList = configList
	handler.messageList = messageList
	handler.searchHitConfig()
	return handler, nil
}

func (handler *FaceHandler) searchHitConfig() {

}

func (handler FaceHandler) GetConfigHitList() []Config {
	return handler.configHitList
}
