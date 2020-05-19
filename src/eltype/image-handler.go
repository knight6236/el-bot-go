package eltype

type ImageHandler struct {
	configList    []Config
	messageList   []Message
	configHitList []Config
}

func NewImageHandler(configList []Config, messageList []Message) (IHandler, error) {
	var handler ImageHandler
	handler.configList = configList
	handler.messageList = messageList
	handler.searchHitConfig()
	return handler, nil
}

func (handler *ImageHandler) searchHitConfig() {

}

func (handler ImageHandler) GetConfigHitList() []Config {
	return handler.configHitList
}
