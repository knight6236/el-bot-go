package eltype

type PlainDoer struct {
	configHitList      []Config
	recivedMessageList []Message
	sendedMessageList  []Message
}

func NewPlainDoer(configHitList []Config, recivedMessageList []Message) (IDoer, error) {
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
				doer.sendedMessageList = append(doer.sendedMessageList, doer.getTextMessage(doMessage))
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

func (doer PlainDoer) GetSendedMessageList() []Message {
	return doer.sendedMessageList
}
