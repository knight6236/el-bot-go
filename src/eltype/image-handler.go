package eltype

type ImageHandler struct {
	configList    []Config
	messageList   []Message
	ConfigHitList []Config
}

func NewImageHandler(configList []Config, messageList []Message) (ImageHandler, error) {
	var handler ImageHandler
	handler.configList = configList
	handler.messageList = messageList
	handler.searchHitConfig()
	return handler, nil
}

func (handler *ImageHandler) searchHitConfig() {

}
