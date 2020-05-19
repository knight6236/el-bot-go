package eltype

type OperationDoer struct {
	configHitList      []Config
	recivedMessageList []Message
	sendedMessageList  []Message
	preDefVarMap       map[string]string
}

func NewOperationDoer(configHitList []Config, recivedMessageList []Message, preDefVarMap map[string]string) (IDoer, error) {
	var doer OperationDoer
	doer.configHitList = configHitList
	doer.recivedMessageList = recivedMessageList
	doer.getSendedMessageList()
	return doer, nil
}

func (doer *OperationDoer) getSendedMessageList() {

}

func (doer OperationDoer) GetSendedMessageList() []Message {
	return doer.sendedMessageList
}
