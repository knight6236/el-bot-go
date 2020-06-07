package eltype

import (
	"fmt"
	"strings"
)

// XMLDoer 表情动作生成类
// @property	configHitList		[]Config			命中的配置列表
// @property	recivedMessageList	[]Message			接收到的消息列表
// @property	sendedMessageList	[]Message			将要发送的消息列表
// @property	sendedOperationList	[]Operation			将要执行的动作列表
// @property	preDefVarMap		map[string]string	预定义变量Map
type XMLDoer struct {
	configHitList       []Config
	recivedMessage      Message
	willBeSentMessage   []Message
	willBeSentOperation []Operation
	willBeSentControl   []Control
	preDefVarMap        map[string]string
}

// NewXMLDoer 构造一个 XMLDoer
// @param	configHitList		[]Config			命中的配置列表
// @param	recivedMessageList	[]Message			接收到的消息列表
// @param	sendedMessageList	[]Message			将要发送的消息列表
// @param	preDefVarMap		map[string]string	预定义变量Map
func NewXMLDoer(configHitList []Config, recivedMessage Message, preDefVarMap map[string]string) (IDoer, error) {
	var doer XMLDoer
	doer.configHitList = configHitList
	doer.recivedMessage = recivedMessage
	doer.preDefVarMap = preDefVarMap
	doer.getWillBeSentMessageList()
	return doer, nil
}

func (doer *XMLDoer) getWillBeSentMessageList() {
	for _, config := range doer.configHitList {
		for _, doMessageDetail := range config.Do.Message.DetailList {
			var willBeSentMessage Message
			var willBeSentMessageDetail MessageDetail
			willBeSentMessage.IsQuote = config.Do.Message.IsQuote
			willBeSentMessage.Sender = config.Do.Message.Sender.DeepCopy()
			willBeSentMessage.Receiver = config.Do.Message.Receiver.DeepCopy()
			willBeSentMessageDetail.InnerType = MessageTypeXML
			if doMessageDetail.InnerType == MessageTypeXML {
				xml, isReplace := doer.replaceStrByPreDefVarMap(doMessageDetail.Text)
				if isReplace {
					willBeSentMessageDetail.Text = xml
				}
				willBeSentMessage.AddDetail(willBeSentMessageDetail)
				doer.willBeSentMessage = append(doer.willBeSentMessage, willBeSentMessage)
			}
		}
	}
}

func (doer XMLDoer) replaceStrByPreDefVarMap(text string) (string, bool) {
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

// GetSendedMessageList 获取将要发送的信息列表
func (doer XMLDoer) GetWillBeSentMessageList() []Message {
	return doer.willBeSentMessage
}

// GetSendedOperationList 获取将要执行的动作列表
func (doer XMLDoer) GetWillBeSentOperationList() []Operation {
	return doer.willBeSentOperation
}

// GetwillBeSentControlList 获取将要执行的动作列表
func (doer XMLDoer) GetwillBeSentControlList() []Control {
	return doer.willBeSentControl
}
