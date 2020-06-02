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
// @property	sendedOperationList	[]Operation			将要执行的动作列表
// @property	preDefVarMap		map[string]string	预定义变量Map
type PlainDoer struct {
	configHitList       []Config
	recivedMessageList  []Message
	willBeSentMessage   []Message
	willBeSentOperation []Operation
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
	doer.getWillBeSentMessageList()
	return doer, nil
}

func (doer *PlainDoer) getWillBeSentMessageList() {
	for _, config := range doer.configHitList {
		for _, doMessageDetail := range config.Do.Message.DetailList {
			var willBeSentMessage Message
			var willBeSentMessageDetail MessageDetail
			willBeSentMessage.IsQuote = config.Do.Message.IsQuote
			willBeSentMessage.Sender = config.Do.Message.Sender.DeepCopy()
			willBeSentMessage.Receiver = config.Do.Message.Receiver.DeepCopy()
			willBeSentMessageDetail.innerType = MessageTypePlain
			if doMessageDetail.innerType == MessageTypePlain {
				if doMessageDetail.URL == "" {
					willBeSentMessageDetail, err := doer.getTextMessageDetail(doMessageDetail)
					willBeSentMessage.AddDetail(willBeSentMessageDetail)
					if err == nil {
						doer.willBeSentMessage = append(doer.willBeSentMessage, willBeSentMessage)
					}
				} else if doMessageDetail.URL != "" {
					willBeSentMessageDetail, err := doer.getURLMessageDetail(doMessageDetail)
					willBeSentMessage.AddDetail(willBeSentMessageDetail)
					if err == nil {
						doer.willBeSentMessage = append(doer.willBeSentMessage, willBeSentMessage)
					}
				}
			}
		}
	}
}

func (doer *PlainDoer) getTextMessageDetail(detail MessageDetail) (MessageDetail, error) {
	var newDetail MessageDetail
	text, isReplace := doer.replaceStrByPreDefVarMap(detail.Text)
	if isReplace {
		newDetail.Text = text
	} else {
		newDetail.Text = detail.Text
	}
	return newDetail, nil
}

func (doer *PlainDoer) getURLMessageDetail(detail MessageDetail) (MessageDetail, error) {
	var newDetail MessageDetail
	res, err := http.Get(detail.URL)
	if err != nil {
		return detail, err
	}

	defer res.Body.Close()

	bodyContent, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return detail, err
	}
	doer.preDefVarMap["el-url-text"] = string(bodyContent)

	text, isReplace := doer.replaceStrByPreDefVarMap(detail.Text)
	if isReplace {
		newDetail.Text = text
	}

	if detail.JSON {
		newDetail.Text = doer.replaceStrByJSON(bodyContent, newDetail.Text)
	}

	return newDetail, nil
}

func (doer PlainDoer) replaceStrByPreDefVarMap(text string) (string, bool) {
	var isReplace bool = false
	for varName, value := range doer.preDefVarMap {
		key := fmt.Sprintf("{%s}", varName)
		temp := text
		text = strings.ReplaceAll(text, key, value)
		if !isReplace && temp == text {
			isReplace = true
		}
	}
	return text, isReplace
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
func (doer PlainDoer) GetWillBeSentMessageList() []Message {
	return doer.willBeSentMessage
}

// GetSendedOperationList 获取将要执行的动作列表
func (doer PlainDoer) GetWillBeSentOperationList() []Operation {
	return doer.willBeSentOperation
}
