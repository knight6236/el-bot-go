package eltype

type FaceDoer struct {
	configHitList      []Config
	recivedMessageList []Message
	sendedMessageList  []Message
}

func NewFaceDoer(configHitList []Config, recivedMessageList []Message) (IDoer, error) {
	var doer FaceDoer
	doer.configHitList = configHitList
	doer.recivedMessageList = recivedMessageList
	doer.getSendedMessageList()
	return doer, nil
}

func (doer *FaceDoer) getSendedMessageList() {

}

func (doer FaceDoer) GetSendedMessageList() []Message {
	return doer.sendedMessageList
}
