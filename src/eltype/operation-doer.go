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
	sendedMessageList   []Message
	sendedOperationList []Operation
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
	doer.getSendedMessageList()
	doer.getSendedOperationList()
	return doer, nil
}

func (doer *OperationDoer) getSendedMessageList() {
}

func (doer *OperationDoer) getSendedOperationList() {
	for _, config := range doer.configHitList {
		for _, doOperation := range config.DoOperationList {
			var operation Operation
			var err error
			var value map[string]string
			if doOperation.Type == OperationTypeAt {
				value = make(map[string]string)
				value["id"] = doer.replaceStrByPreDefVarMap(doOperation.Value["id"])
				operation, err = NewOperation(OperationTypeAt, value)
				if err != nil {
					continue
				}
			}
			doer.sendedOperationList = append(doer.sendedOperationList, operation)
		}
	}
}

func (doer OperationDoer) replaceStrByPreDefVarMap(text string) string {
	for varName, value := range doer.preDefVarMap {
		key := fmt.Sprintf("{%s}", varName)
		text = strings.ReplaceAll(text, key, value)
	}
	return text
}

// GetSendedMessageList 获取将要发送的信息列表
func (doer OperationDoer) GetSendedMessageList() []Message {
	return doer.sendedMessageList
}

// GetSendedOperationList 获取将要执行的动作列表
func (doer OperationDoer) GetSendedOperationList() []Operation {
	return doer.sendedOperationList
}
