package eltype

type PlainHandler struct {
	configList    []Config
	messageList   []Message
	ConfigHitList []Config
}

func NewPlainHandler(configList []Config, messageList []Message) (PlainHandler, error) {
	var handler PlainHandler
	handler.configList = configList
	handler.messageList = messageList
	handler.searchHitConfig()
	return handler, nil
}

func (handler *PlainHandler) searchHitConfig() {
	for _, config := range handler.configList {
		for _, whenMessage := range config.WhenMessageList {
			if handler.checkText(whenMessage) {
				handler.ConfigHitList = append(handler.ConfigHitList, config)
				break
			}
		}
	}
}

func (handler *PlainHandler) checkText(whenMessage Message) bool {
	if whenMessage.Type != MessageTypePlain {
		return false
	}
	text := whenMessage.Value["text"]
	if text == "" {
		return false
	}
	for _, message := range handler.messageList {
		if message.Value["text"] == text {
			return true
		}
	}
	return false
}
