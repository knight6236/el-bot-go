package eltype

import (
	"fmt"
	"strings"
)

// OperationDoer 表情动作生成类
// @property	configHitList		[]Config			命中的配置列表
// @property	recivedMessageList	[]Message			接收到的消息列表
// @property	sendedMessageList	[]Message			将要发送的消息列表
// @property	sendedOperationList	[]Operation			将要执行的动作列表
// @property	preDefVarMap		map[string]string	预定义变量Map
type OperationDoer struct {
	configHitList       []Config
	recivedMessageList  []Message
	willBeSentMessage   []Message
	willBeSentOperation []Operation
	preDefVarMap        map[string]string
}

// NewOperationDoer 构造一个 OperationDoer
// @param	configHitList		[]Config			命中的配置列表
// @param	recivedMessageList	[]Message			接收到的消息列表
// @param	sendedMessageList	[]Message			将要发送的消息列表
// @param	preDefVarMap		map[string]string	预定义变量Map
func NewOperationDoer(configHitList []Config, recivedMessageList []Message, preDefVarMap map[string]string) (IDoer, error) {
	var doer OperationDoer
	doer.configHitList = configHitList
	doer.recivedMessageList = recivedMessageList
	doer.preDefVarMap = preDefVarMap
	doer.getWillBeSentMessageList()
	doer.getSendedOperationList()
	return doer, nil
}

func (doer *OperationDoer) getWillBeSentMessageList() {
	for _, config := range doer.configHitList {
		for _, doOperation := range config.Do.OperationList {
			var willBeSentMessage Message
			var willBeSentMessageDetail MessageDetail
			switch doOperation.innerType {
			case OperationTypeAt:
				willBeSentMessage.Receiver.AddGroupID(doOperation.GroupID)
				willBeSentMessage.Receiver.AddUserID(doOperation.UserID)
				willBeSentMessageDetail.innerType = MessageTypeAt
				willBeSentMessageDetail.GroupID = doOperation.GroupID
				willBeSentMessageDetail.UserID = doOperation.UserID
			case OperationTypeAtAll:
				willBeSentMessage.Receiver.AddGroupID(doOperation.GroupID)
				willBeSentMessageDetail.innerType = MessageTypeAtAll
				willBeSentMessageDetail.GroupID = doOperation.GroupID
			default:
				continue
			}
			willBeSentMessage.AddDetail(willBeSentMessageDetail)
			doer.willBeSentMessage = append(doer.willBeSentMessage, willBeSentMessage)
		}
	}
}

func (doer *OperationDoer) getSendedOperationList() {
	for _, config := range doer.configHitList {
		for _, doOperation := range config.Do.OperationList {
			var operation Operation
			switch doOperation.innerType {
			case OperationTypeMemberMute:
				operation.innerType = OperationTypeMemberMute
				groupID, isReplace := doer.replaceStrByPreDefVarMap(doOperation.GroupID)
				if !isReplace {
					operation.GroupID = doOperation.GroupID
				}
				operation.GroupID = groupID

				userID, isReplace := doer.replaceStrByPreDefVarMap(doOperation.UserID)
				if !isReplace {
					operation.UserID = doOperation.UserID
				}
				operation.UserID = userID

				operation.Second = doOperation.Second
			case OperationTypeMemberUnMute:
				operation.innerType = OperationTypeMemberUnMute
				groupID, isReplace := doer.replaceStrByPreDefVarMap(doOperation.GroupID)
				if !isReplace {
					operation.GroupID = doOperation.GroupID
				}
				operation.GroupID = groupID

				userID, isReplace := doer.replaceStrByPreDefVarMap(doOperation.UserID)
				if !isReplace {
					operation.UserID = doOperation.UserID
				}
				operation.UserID = userID
			case OperationTypeGroupMuteAll:
				operation.innerType = OperationTypeGroupMuteAll
				groupID, isReplace := doer.replaceStrByPreDefVarMap(doOperation.GroupID)
				if !isReplace {
					operation.GroupID = doOperation.GroupID
				}
				operation.GroupID = groupID
			case OperationTypeGroupUnMuteAll:
				operation.innerType = OperationTypeGroupUnMuteAll
				groupID, isReplace := doer.replaceStrByPreDefVarMap(doOperation.GroupID)
				if !isReplace {
					operation.GroupID = doOperation.GroupID
				}
				operation.GroupID = groupID
			default:
				continue
			}
			doer.willBeSentOperation = append(doer.willBeSentOperation, operation)
		}
	}
}

func (doer OperationDoer) replaceStrByPreDefVarMap(text string) (string, bool) {
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
func (doer OperationDoer) GetWillBeSentMessageList() []Message {
	return doer.willBeSentMessage
}

// GetSendedOperationList 获取将要执行的动作列表
func (doer OperationDoer) GetWillBeSentOperationList() []Operation {
	return doer.willBeSentOperation
}
