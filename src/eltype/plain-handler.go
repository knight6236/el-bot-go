package eltype

type PlainHandler struct {
	configList    []Config
	messageList   []Message
	configHitList []Config
}

func NewPlainHandler(configList []Config, messageList []Message) (IHandler, error) {
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
				handler.configHitList = append(handler.configHitList, config)
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

func (handler PlainHandler) GetConfigHitList() []Config {
	return handler.configHitList
}
