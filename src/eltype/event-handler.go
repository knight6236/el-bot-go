package eltype

type EventHandler struct {
	configList    []Config
	messageList   []Message
	configHitList []Config
}

func NewEventHandler(configList []Config, messageList []Message) (IHandler, error) {
	var handler EventHandler
	handler.configList = configList
	handler.messageList = messageList
	handler.searchHitConfig()
	return handler, nil
}

func (handler *EventHandler) searchHitConfig() {

}

func (handler EventHandler) GetConfigHitList() []Config {
	return handler.configHitList
}
