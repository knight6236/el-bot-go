package eltype

import (
	"fmt"
	"strings"
)

// 「」

// AtDoer 表情动作生成类
// @property	configHitList		[]Config			命中的配置列表
// @property	recivedMessageList	[]Message			接收到的消息列表
// @property	sendedMessageList	[]Message			将要发送的消息列表
// @property	sendedOperationList	[]Operation			将要执行的动作列表
// @property	preDefVarMap		map[string]string	预定义变量Map
type AtDoer struct {
	configHitList       []Config
	recivedMessageList  []Message
	sendedMessageList   []Message
	sendedOperationList []Operation
	preDefVarMap        map[string]string
}

// NewAtDoer 构造一个 AtDoer
// @param	configHitList		[]Config			命中的配置列表
// @param	recivedMessageList	[]Message			接收到的消息列表
// @param	preDefVarMap		map[string]string	预定义变量 Map
func NewAtDoer(configHitList []Config, recivedMessageList []Message, preDefVarMap map[string]string) (IDoer, error) {
	var doer AtDoer
	doer.configHitList = configHitList
	doer.recivedMessageList = recivedMessageList
	doer.preDefVarMap = preDefVarMap
	doer.getSendedMessageList()
	return doer, nil
}

func (doer *AtDoer) getSendedMessageList() {
	for _, config := range doer.configHitList {
		for _, doMessageDetail := range config.Do.Message.DetailList {
			var willBeSentMessage Message
			var willBeSentMessageDetail MessageDetail
			willBeSentMessage.Sender = config.Do.Message.Sender.DeepCopy()
			willBeSentMessage.Receiver = config.Do.Message.Receiver.DeepCopy()
			switch doMessageDetail.innerType {
			case MessageTypeAt:
				willBeSentMessageDetail.innerType = MessageTypeAt
				groupID, isReplace := doer.replaceStrByPreDefVarMap(doMessageDetail.GroupID)
				if !isReplace {
					groupID = doMessageDetail.GroupID
				}
				userID, isReplace := doer.replaceStrByPreDefVarMap(doMessageDetail.UserID)
				if !isReplace {
					userID = doMessageDetail.UserID
				}
				willBeSentMessageDetail.GroupID = groupID
				willBeSentMessageDetail.UserID = userID
				willBeSentMessage.AddDetail(willBeSentMessageDetail)
				doer.sendedMessageList = append(doer.sendedMessageList, willBeSentMessage)
			case MessageTypeAtAll:
				willBeSentMessageDetail.innerType = MessageTypeAtAll
				groupID, isReplace := doer.replaceStrByPreDefVarMap(doMessageDetail.GroupID)
				if !isReplace {
					groupID = doMessageDetail.GroupID
				}
				willBeSentMessageDetail.GroupID = groupID
				willBeSentMessage.AddDetail(willBeSentMessageDetail)
				doer.sendedMessageList = append(doer.sendedMessageList, willBeSentMessage)
			}
		}
	}
}

func (doer AtDoer) replaceStrByPreDefVarMap(text string) (string, bool) {
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
func (doer AtDoer) GetSendedMessageList() []Message {
	return doer.sendedMessageList
}

// GetSendedOperationList 获取将要执行的动作列表
func (doer AtDoer) GetSendedOperationList() []Operation {
	return doer.sendedOperationList
}
