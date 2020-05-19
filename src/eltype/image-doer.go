package eltype

type ImageDoer struct {
	configHitList      []Config
	recivedMessageList []Message
	sendedMessageList  []Message
	preDefVarMap       map[string]string
}

func NewImageDoer(configHitList []Config, recivedMessageList []Message, preDefVarMap map[string]string) (IDoer, error) {
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
