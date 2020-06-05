package eltype

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	willBeSentControl   []Control
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
	var newDetail MessageDetail = detail.DeepCopy()
	text, isReplace := doer.replaceStrByPreDefVarMap(detail.Text)
	if isReplace {
		newDetail.Text = text
	} else {
		newDetail.Text = detail.Text
	}
	return newDetail, nil
}

func (doer *PlainDoer) getURLMessageDetail(detail MessageDetail) (MessageDetail, error) {
	var newDetail MessageDetail = detail.DeepCopy()
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

	if detail.JSON {
		doer.addPreDefVarByJSON(bodyContent)
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

func (doer *PlainDoer) addPreDefVarByJSON(jsonByteList []byte) {
	var jsonMap interface{}
	err := json.Unmarshal(jsonByteList, &jsonMap)
	if err != nil {
		return
	}

	varNameList, valueList := parseJsonObj(jsonMap, 0)

	for i := 0; i < len(varNameList); i++ {
		doer.preDefVarMap[varNameList[i]] = valueList[i]
	}
}

// GetSendedMessageList 获取将要发送的信息列表
func (doer PlainDoer) GetWillBeSentMessageList() []Message {
	return doer.willBeSentMessage
}

// GetSendedOperationList 获取将要执行的动作列表
func (doer PlainDoer) GetWillBeSentOperationList() []Operation {
	return doer.willBeSentOperation
}

// GetwillBeSentControlList 获取将要执行的动作列表
func (doer PlainDoer) GetwillBeSentControlList() []Control {
	return doer.willBeSentControl
}
