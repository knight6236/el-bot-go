package eltype

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
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
	preDefVarMap  *map[string]string
}

// NewPlainHandler 构造一个 PlainHandler
// @param	configList		[]Config			要判断的配置列表
// @param	messageList		[]Message			要判断的消息列表
// @param	operationList	[]Operation			要判断的配置列表
// @param	preDefVarMap	map[string]string	预定义变量 Map
func NewPlainHandler(configList []Config, messageList []Message, operationList []Operation,
	preDefVarMap *map[string]string) (IHandler, error) {
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
	pattern := whenMessage.Value["regex"]

	if pattern == "" {
		return false
	}

	var buf bytes.Buffer

	for _, message := range handler.messageList {
		if message.Type == MessageTypePlain && message.Value["text"] != "" {
			buf.WriteString(message.Value["text"])
		}
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println(err)
		return false
	}

	matches := regex.FindStringSubmatch(buf.String())
	if matches == nil {
		return false
	}
	for i := 1; i < len(matches); i++ {
		varName := fmt.Sprintf("el-regex-%d", i-1)
		handler.addPerDefVar(varName, matches[i])
	}
	return true
}

func (handler *PlainHandler) addPerDefVar(varName string, value interface{}) {
	switch value.(type) {
	case string:
		(*(handler.preDefVarMap))[varName] = value.(string)
	case int:
		(*(handler.preDefVarMap))[varName] = strconv.Itoa(value.(int))
	case int64:
		(*(handler.preDefVarMap))[varName] = strconv.FormatInt(value.(int64), 10)
	case float64:
		(*(handler.preDefVarMap))[varName] = fmt.Sprintf("%.2f", value.(float64))
	}
}

// GetConfigHitList 获取命中的配置列表
func (handler PlainHandler) GetConfigHitList() []Config {
	return handler.configHitList
}
