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
}

func NewPlainDoer(configHitList []Config, recivedMessageList []Message) (IDoer, error) {
	var doer PlainDoer
	doer.configHitList = configHitList
	doer.recivedMessageList = recivedMessageList
	doer.getSendedMessageList()
	return doer, nil
}

func (doer *PlainDoer) getSendedMessageList() {
	for _, config := range doer.configHitList {
		for _, doMessage := range config.DoMessageList {
			if doMessage.Type != MessageTypePlain {
				continue
			}
			if doMessage.Value["text"] != "" {
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
	value["text"] = message.Value["text"]
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

	value := make(map[string]string)
	var sendedMessage Message

	if message.Value["jtext"] == "" {
		value["text"] = string(bodyContent)
		sendedMessage, err = NewMessage(MessageTypePlain, value)
		if err != nil {
			return message, err
		}
	} else {
		value["text"] = doer.replaceStrByJson(bodyContent, message.Value["jtext"])
		sendedMessage, err = NewMessage(MessageTypePlain, value)
		if err != nil {
			return message, err
		}
	}

	return sendedMessage, nil
}

func (doer *PlainDoer) replaceStrByJson(jsonByteList []byte, jtext string) string {
	var jsonMap interface{}
	err := json.Unmarshal(jsonByteList, &jsonMap)
	if err != nil {
		return ""
	}

	for nativeKey, nativeValue := range jsonMap.(map[string]interface{}) {
		key := fmt.Sprintf("{%s}", nativeKey)
		switch nativeValue.(type) {
		case string:
			jtext = strings.ReplaceAll(jtext, key, nativeValue.(string))
		case int:
			value := strconv.Itoa(nativeValue.(int))
			jtext = strings.ReplaceAll(jtext, key, value)
		case int64:
			value := strconv.FormatInt(nativeValue.(int64), 10)
			jtext = strings.ReplaceAll(jtext, key, value)
		case float64:
			value := fmt.Sprintf("%.6f", nativeValue.(float64))
			jtext = strings.ReplaceAll(jtext, key, value)
		case bool:
			value := strconv.FormatBool(nativeValue.(bool))
			jtext = strings.ReplaceAll(jtext, key, value)
		}
	}
	return jtext
}

func (doer PlainDoer) GetSendedMessageList() []Message {
	return doer.sendedMessageList
}
