package eltype

type EventDoer struct {
	configHitList      []Config
	recivedMessageList []Message
	sendedMessageList  []Message
}

func NewEventDoer(configHitList []Config, recivedMessageList []Message) (IDoer, error) {
	var doer EventDoer
	doer.configHitList = configHitList
	doer.recivedMessageList = recivedMessageList
	doer.getSendedMessageList()
	return doer, nil
}

func (doer *EventDoer) getSendedMessageList() {

}

func (doer EventDoer) GetSendedMessageList() []Message {
	return doer.sendedMessageList
}
