package eltype

type PlainDoer struct {
	configHitList      []Config
	recivedMessageList []Message
	SendedMessageList  []Message
}

func NewPlainDoer(configHitList []Config, recivedMessageList []Message) (PlainDoer, error) {
	var doer PlainDoer
	doer.configHitList = configHitList
	doer.recivedMessageList = recivedMessageList
	doer.getSendedMessageList()
	return doer, nil
}

func (doer *PlainDoer) getSendedMessageList() {
	for _, config := range doer.configHitList {
		for _, doMessage := range config.DoMessageList {
			if doMessage.Type == MessageTypePlain && doMessage.Value["text"] != "" {
				doer.SendedMessageList = append(doer.SendedMessageList, doer.getTextMessage(doMessage))
			}
		}
	}
}

func (doer *PlainDoer) getTextMessage(message Message) Message {
	value := make(map[string]string)
	value["text"] = message.Value["text"]
	message, err := NewMessage(MessageTypePlain, value)
	if err != nil {
	}
	return message
}
