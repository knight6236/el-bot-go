package eltype

// IDoer 所有 Doer 均实现此接口
// @method	GetSendedMessageList	[]Message	获取将要发送的消息列表
type IDoer interface {
	replaceStrByPreDefVarMap(text string) string
	GetSendedMessageList() []Message
	GetSendedOperationList() []Operation
}
