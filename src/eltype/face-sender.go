package eltype

type FaceDoer struct {
	configHitList      []Config
	recivedMessageList []Message
	SendedMessageList  []Message
}

func NewFaceDoer(configHitList []Config, recivedMessageList []Message) (FaceDoer, error) {
	var doer FaceDoer
	doer.configHitList = configHitList
	doer.recivedMessageList = recivedMessageList
	doer.getSendedMessageList()
	return doer, nil
}

func (doer *FaceDoer) getSendedMessageList() {

}
