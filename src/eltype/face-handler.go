package eltype

type FaceHandler struct {
	configList    []Config
	messageList   []Message
	ConfigHitList []Config
}

func NewFaceHandler(configList []Config, messageList []Message) (FaceHandler, error) {
	var handler FaceHandler
	handler.configList = configList
	handler.messageList = messageList
	handler.searchHitConfig()
	return handler, nil
}

func (handler *FaceHandler) searchHitConfig() {

}
