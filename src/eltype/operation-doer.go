package eltype

type OperationDoer struct {
	configHitList      []Config
	recivedMessageList []Message
	sendedMessageList  []Message
}

func NewOperationDoer(configHitList []Config, recivedMessageList []Message) (IDoer, error) {
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
