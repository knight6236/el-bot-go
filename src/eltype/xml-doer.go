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
	recivedMessageList  []Message
	sendedMessageList   []Message
	sendedOperationList []Operation
	preDefVarMap        map[string]string
}

// NewXMLDoer 构造一个 XMLDoer
// @param	configHitList		[]Config			命中的配置列表
// @param	recivedMessageList	[]Message			接收到的消息列表
// @param	sendedMessageList	[]Message			将要发送的消息列表
// @param	preDefVarMap		map[string]string	预定义变量Map
func NewXMLDoer(configHitList []Config, recivedMessageList []Message, preDefVarMap map[string]string) (IDoer, error) {
	var doer XMLDoer
	doer.configHitList = configHitList
	doer.recivedMessageList = recivedMessageList
	doer.preDefVarMap = preDefVarMap
	doer.getSendedMessageList()
	return doer, nil
}

func (doer *XMLDoer) getSendedMessageList() {
	for _, config := range doer.configHitList {
		for _, doMessage := range config.DoMessageList {
			if doMessage.Type == MessageTypeXML {
				value := make(map[string]string)
				xml, isReplace := doer.replaceStrByPreDefVarMap(doMessage.Value["text"])
				if isReplace {
					value["xml"] = xml
				}
				message, err := NewMessage(MessageTypeXML, value)
				if err != nil {
					continue
				}
				doer.sendedMessageList = append(doer.sendedMessageList, message)
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
func (doer XMLDoer) GetSendedMessageList() []Message {
	return doer.sendedMessageList
}

// GetSendedOperationList 获取将要执行的动作列表
func (doer XMLDoer) GetSendedOperationList() []Operation {
	return doer.sendedOperationList
}
