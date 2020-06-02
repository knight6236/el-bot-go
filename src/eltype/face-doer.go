package eltype

import (
	"fmt"
	"strings"
)

// 「」

// FaceDoer 表情动作生成类
// @property	configHitList		[]Config			命中的配置列表
// @property	recivedMessageList	[]Message			接收到的消息列表
// @property	sendedMessageList	[]Message			将要发送的消息列表
// @property	sendedOperationList	[]Operation			将要执行的动作列表
// @property	preDefVarMap		map[string]string	预定义变量Map
type FaceDoer struct {
	configHitList       []Config
	recivedMessageList  []Message
	willBeSentMessage   []Message
	willBeSentOperation []Operation
	preDefVarMap        map[string]string
}

// NewFaceDoer 构造一个 FaceDoer
// @param	configHitList		[]Config			命中的配置列表
// @param	recivedMessageList	[]Message			接收到的消息列表
// @param	preDefVarMap		map[string]string	预定义变量 Map
func NewFaceDoer(configHitList []Config, recivedMessageList []Message, preDefVarMap map[string]string) (IDoer, error) {
	var doer FaceDoer
	doer.configHitList = configHitList
	doer.recivedMessageList = recivedMessageList
	doer.preDefVarMap = preDefVarMap
	doer.getWillBeSentMessageList()
	return doer, nil
}

func (doer *FaceDoer) getWillBeSentMessageList() {
	for _, config := range doer.configHitList {
		for _, doMessageDetail := range config.Do.Message.DetailList {
			// var willBeSentMessage Message
			// var willBeSentMessageDetail MessageDetail
			// willBeSentMessage = config.Do.Message.DeepCopy()
			// willBeSentMessage.Sender = config.Do.Message.Sender.DeepCopy()
			// willBeSentMessage.Receiver = config.Do.Message.Receiver.DeepCopy()
			// willBeSentMessageDetail.innerType = MessageTypeFace
			if doMessageDetail.innerType == MessageTypeFace {
				doer.willBeSentMessage = append(doer.willBeSentMessage, config.Do.Message.DeepCopy())
			}
		}
	}
}

func (doer FaceDoer) replaceStrByPreDefVarMap(text string) (string, bool) {
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
func (doer FaceDoer) GetWillBeSentMessageList() []Message {
	return doer.willBeSentMessage
}

// GetSendedOperationList 获取将要执行的动作列表
func (doer FaceDoer) GetWillBeSentOperationList() []Operation {
	return doer.willBeSentOperation
}
