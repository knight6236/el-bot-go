package eltype

type ImageDoer struct {
	configHitList      []Config
	recivedMessageList []Message
	SendedMessageList  []Message
}

func NewImageDoer(configHitList []Config, recivedMessageList []Message) (ImageDoer, error) {
	var doer ImageDoer
	doer.configHitList = configHitList
	doer.recivedMessageList = recivedMessageList
	doer.getSendedMessageList()
	return doer, nil
}

func (doer *ImageDoer) getSendedMessageList() {

}
