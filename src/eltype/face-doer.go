package eltype

// 「」

// FaceDoer 表情动作生成类
// @property	configHitList		[]Config			命中的配置列表
// @property	recivedMessageList	[]Message			接收到的消息列表
// @property	sendedMessageList	[]Message			将要发送的消息列表
// @property	preDefVarMap		map[string]string	预定义变量Map
type FaceDoer struct {
	configHitList      []Config
	recivedMessageList []Message
	sendedMessageList  []Message
	preDefVarMap       map[string]string
}

// NewFaceDoer 构造一个 FaceDoer
// @param	configHitList		[]Config			命中的配置列表
// @param	recivedMessageList	[]Message			接收到的消息列表
// @param	preDefVarMap		map[string]string	预定义变量 Map
func NewFaceDoer(configHitList []Config, recivedMessageList []Message, preDefVarMap map[string]string) (IDoer, error) {
	var doer FaceDoer
	doer.configHitList = configHitList
	doer.recivedMessageList = recivedMessageList
	doer.getSendedMessageList()
	return doer, nil
}

func (doer *FaceDoer) getSendedMessageList() {

}

// GetSendedMessageList 获取将要发送的信息列表
func (doer FaceDoer) GetSendedMessageList() []Message {
	return doer.sendedMessageList
}
