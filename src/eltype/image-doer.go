package eltype

// ImageDoer 表情动作生成类
// @property	configHitList		[]Config			命中的配置列表
// @property	recivedMessageList	[]Message			接收到的消息列表
// @property	sendedMessageList	[]Message			将要发送的消息列表
// @property	preDefVarMap		map[string]string	预定义变量Map
type ImageDoer struct {
	configHitList      []Config
	recivedMessageList []Message
	sendedMessageList  []Message
	preDefVarMap       map[string]string
}

// NewImageDoer 构造一个 ImageDoer
// @param	configHitList		[]Config			命中的配置列表
// @param	recivedMessageList	[]Message			接收到的消息列表
// @param	sendedMessageList	[]Message			将要发送的消息列表
// @param	preDefVarMap		map[string]string	预定义变量Map
func NewImageDoer(configHitList []Config, recivedMessageList []Message, preDefVarMap map[string]string) (IDoer, error) {
	var doer ImageDoer
	doer.configHitList = configHitList
	doer.recivedMessageList = recivedMessageList
	doer.getSendedMessageList()
	return doer, nil
}

func (doer *ImageDoer) getSendedMessageList() {

}

// GetSendedMessageList 获取将要发送的信息列表
func (doer ImageDoer) GetSendedMessageList() []Message {
	return doer.sendedMessageList
}
