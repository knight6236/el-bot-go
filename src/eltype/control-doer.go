package eltype

import (
	"fmt"
	"strings"
)

// 「」

// ControlDoer
// @property	configHitList		[]Config			命中的配置列表
// @property	recivedMessageList	[]Message			接收到的消息列表
// @property	sendedMessageList	[]Message			将要发送的消息列表
// @property	sendedOperationList	[]Operation			将要执行的动作列表
// @property	preDefVarMap		map[string]string	预定义变量Map
type ControlDoer struct {
	configHitList       []Config
	recivedMessage      Message
	willBeSentMessage   []Message
	willBeSentOperation []Operation
	willBeSentControl   []Control
	preDefVarMap        map[string]string
}

// NewControlDoer 构造一个 ControlDoer
// @param	configHitList		[]Config			命中的配置列表
// @param	recivedMessageList	[]Message			接收到的消息列表
// @param	preDefVarMap		map[string]string	预定义变量 Map
func NewControlDoer(configHitList []Config, recivedMessage Message, preDefVarMap map[string]string) (IDoer, error) {
	var doer ControlDoer
	doer.configHitList = configHitList
	doer.recivedMessage = recivedMessage
	doer.preDefVarMap = preDefVarMap
	doer.getWillBeSentMessageList()
	doer.searchWillBeSentControlList()
	return doer, nil
}

func (doer *ControlDoer) getWillBeSentMessageList() {
}

func (doer ControlDoer) replaceStrByPreDefVarMap(text string) (string, bool) {
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

func (doer *ControlDoer) searchWillBeSentControlList() {
	for _, config := range doer.configHitList {
		for _, control := range config.Do.ControlList {
			doer.willBeSentControl = append(doer.willBeSentControl, control.DeepCopy())
		}
	}
}

// GetWillBeSentMessageList 获取将要发送的信息列表
func (doer ControlDoer) GetWillBeSentMessageList() []Message {
	return doer.willBeSentMessage
}

// GetWillBeSentOperationList 获取将要执行的动作列表
func (doer ControlDoer) GetWillBeSentOperationList() []Operation {
	return doer.willBeSentOperation
}

// GetwillBeSentControlList 获取将要执行的动作列表
func (doer ControlDoer) GetwillBeSentControlList() []Control {
	return doer.willBeSentControl
}
