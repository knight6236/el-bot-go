package eltype

type EventHandler struct {
	configList    []Config
	messageList   []Message
	ConfigHitList []Config
}

func NewEventHandler(configList []Config, messageList []Message) (EventHandler, error) {
	var handler EventHandler
	handler.configList = configList
	handler.messageList = messageList
	handler.searchHitConfig()
	return handler, nil
}

func (handler *EventHandler) searchHitConfig() {

}
