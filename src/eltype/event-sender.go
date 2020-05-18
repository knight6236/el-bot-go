package eltype

type EventDoer struct {
	configHitList      []Config
	recivedMessageList []Message
	SendedMessageList  []Message
}

func NewEventDoer(configHitList []Config, recivedMessageList []Message) (EventDoer, error) {
	var doer EventDoer
	doer.configHitList = configHitList
	doer.recivedMessageList = recivedMessageList
	doer.getSendedMessageList()
	return doer, nil
}

func (doer *EventDoer) getSendedMessageList() {

}
