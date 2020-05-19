package eltype

type ImageDoer struct {
	configHitList      []Config
	recivedMessageList []Message
	sendedMessageList  []Message
}

func NewImageDoer(configHitList []Config, recivedMessageList []Message) (IDoer, error) {
	var doer ImageDoer
	doer.configHitList = configHitList
	doer.recivedMessageList = recivedMessageList
	doer.getSendedMessageList()
	return doer, nil
}

func (doer *ImageDoer) getSendedMessageList() {

}

func (doer ImageDoer) GetSendedMessageList() []Message {
	return doer.sendedMessageList
}
