package eltype

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type PlainDoer struct {
	configHitList      []Config
	recivedMessageList []Message
	sendedMessageList  []Message
	preDefVarMap       map[string]string
}

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
				sendedMessage, err := doer.getUrlMessage(doMessage)
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

func (doer *PlainDoer) getUrlMessage(message Message) (Message, error) {
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
		value["text"] = doer.replaceStrByJson(bodyContent, value["text"])
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

func (doer *PlainDoer) replaceStrByJson(jsonByteList []byte, text string) string {
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

func (doer PlainDoer) GetSendedMessageList() []Message {
	return doer.sendedMessageList
}
