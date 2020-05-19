package eltype

import (
	"bytes"
	"regexp"
)

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
			if handler.checkText(whenMessage) ||
				handler.checkRegex(whenMessage) {
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

func (handler *PlainHandler) checkRegex(whenMessage Message) bool {
	if whenMessage.Type != MessageTypePlain {
		return false
	}
	regex := whenMessage.Value["regex"]

	if regex == "" {
		return false
	}

	var buf bytes.Buffer

	for _, message := range handler.messageList {
		if message.Type == MessageTypePlain && message.Value["text"] != "" {
			buf.WriteString(message.Value["text"])
		}
	}

	isMatch, err := regexp.MatchString(regex, buf.String())

	if err != nil {
		return false
	}

	return isMatch
}

func (handler PlainHandler) GetConfigHitList() []Config {
	return handler.configHitList
}
