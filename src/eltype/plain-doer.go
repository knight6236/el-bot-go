package eltype

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// PlainDoer 表情动作生成类
// @property	configHitList		[]Config			命中的配置列表
// @property	recivedMessageList	[]Message			接收到的消息列表
// @property	sendedMessageList	[]Message			将要发送的消息列表
// @property	preDefVarMap		map[string]string	预定义变量Map
type PlainDoer struct {
	configHitList       []Config
	recivedMessageList  []Message
	sendedMessageList   []Message
	sendedOperationList []Operation
	preDefVarMap        map[string]string
}

// NewPlainDoer 构造一个 PlainDoer
// @param	configHitList		[]Config			命中的配置列表
// @param	recivedMessageList	[]Message			接收到的消息列表
// @param	sendedMessageList	[]Message			将要发送的消息列表
// @param	preDefVarMap		map[string]string	预定义变量Map
func NewPlainDoer(configHitList []Config, recivedMessageList []Message, preDefVarMap map[string]string) (IDoer, error) {
	var doer PlainDoer
	doer.configHitList = configHitList
	doer.recivedMessageList = recivedMessageList
	doer.preDefVarMap = preDefVarMap
	doer.getSendedMessageList()
	return doer, nil
}

func (doer *PlainDoer) getSendedMessageList() {
	for _, config := range doer.configHitList {
		for _, doMessage := range config.DoMessageList {
			if doMessage.Type != MessageTypePlain {
				continue
			}
			if doMessage.Value["url"] == "" {
				sendedMessage, err := doer.getTextMessage(doMessage)
				if err == nil {
					doer.sendedMessageList = append(doer.sendedMessageList, sendedMessage)
				}
			} else if doMessage.Value["url"] != "" {
				sendedMessage, err := doer.getURLMessage(doMessage)
				if err == nil {
					doer.sendedMessageList = append(doer.sendedMessageList, sendedMessage)
				}
			}
		}
	}
}

func (doer *PlainDoer) getTextMessage(message Message) (Message, error) {
	value := make(map[string]string)
	value["text"] = doer.replaceStrByPreDefVarMap(message.Value["text"])
	sendedMessage, err := NewMessage(MessageTypePlain, value)
	if err != nil {
		return sendedMessage, err
	}
	return sendedMessage, nil
}

func (doer *PlainDoer) getURLMessage(message Message) (Message, error) {
	res, err := http.Get(message.Value["url"])
	if err != nil {
		return message, err
	}

	defer res.Body.Close()

	bodyContent, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return message, err
	}
	doer.preDefVarMap["el-url-text"] = string(bodyContent)

	value := make(map[string]string)

	value["text"] = doer.replaceStrByPreDefVarMap(message.Value["text"])

	if message.Value["json"] == "true" {
		value["text"] = doer.replaceStrByJSON(bodyContent, value["text"])
		if err != nil {
			return message, err
		}
	}

	var sendedMessage Message
	sendedMessage, err = NewMessage(MessageTypePlain, value)
	if err != nil {
		return message, err
	}

	return sendedMessage, nil
}

func (doer *PlainDoer) replaceStrByPreDefVarMap(text string) string {
	for varName, value := range doer.preDefVarMap {
		key := fmt.Sprintf("{%s}", varName)
		text = strings.ReplaceAll(text, key, value)
	}
	return text
}

func (doer *PlainDoer) replaceStrByJSON(jsonByteList []byte, text string) string {
	var jsonMap interface{}
	err := json.Unmarshal(jsonByteList, &jsonMap)
	if err != nil {
		return ""
	}

	for nativeKey, nativeValue := range jsonMap.(map[string]interface{}) {
		key := fmt.Sprintf("{%s}", nativeKey)
		switch nativeValue.(type) {
		case string:
			text = strings.ReplaceAll(text, key, nativeValue.(string))
		case int:
			value := strconv.Itoa(nativeValue.(int))
			text = strings.ReplaceAll(text, key, value)
		case int64:
			value := strconv.FormatInt(nativeValue.(int64), 10)
			text = strings.ReplaceAll(text, key, value)
		case float64:
			value := fmt.Sprintf("%.6f", nativeValue.(float64))
			text = strings.ReplaceAll(text, key, value)
		case bool:
			value := strconv.FormatBool(nativeValue.(bool))
			text = strings.ReplaceAll(text, key, value)
		}
	}
	return text
}

// GetSendedMessageList 获取将要发送的信息列表
func (doer PlainDoer) GetSendedMessageList() []Message {
	return doer.sendedMessageList
}

// GetSendedOperationList 获取将要执行的动作列表
func (doer PlainDoer) GetSendedOperationList() []Operation {
	return doer.sendedOperationList
}
