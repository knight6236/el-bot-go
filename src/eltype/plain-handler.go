package eltype

import (
	"bytes"
	"regexp"
)

// PlainHandler 判断是否命中和表情有关的配置
// @property	configList		[]Config			要判断的配置列表
// @property	messageList		[]Message			要判断的消息列表
// @property	operationList	[]Operation			要判断的配置列表
// @property	configHitList	[]Config			命中的配置列表
// @property	preDefVarMap	map[string]string	预定义变量 Map
type PlainHandler struct {
	configList    []Config
	messageList   []Message
	configHitList []Config
	operationList []Operation
	preDefVarMap  map[string]string
}

// NewPlainHandler 构造一个 PlainHandler
// @param	configList		[]Config			要判断的配置列表
// @param	messageList		[]Message			要判断的消息列表
// @param	operationList	[]Operation			要判断的配置列表
// @param	preDefVarMap	map[string]string	预定义变量 Map
func NewPlainHandler(configList []Config, messageList []Message, operationList []Operation,
	preDefVarMap map[string]string) (IHandler, error) {
	var handler PlainHandler
	handler.configList = configList
	handler.messageList = messageList
	handler.preDefVarMap = preDefVarMap
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

// GetConfigHitList 获取命中的配置列表
func (handler PlainHandler) GetConfigHitList() []Config {
	return handler.configHitList
}
