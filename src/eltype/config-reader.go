package eltype

import (
	"strconv"

	"gopkg.in/yaml.v2"

	// "io"

	// "os"
	"io/ioutil"
	// "reflect"
)

// ConfigReader 配置读取对象
type ConfigReader struct {
	filePath         string
	GlobalConfigList []Config
	FriendConifgList []Config
	GroupConfigList  []Config
}

// NewConfigReader 使用配置文件路径构造一个 ConfigReader
// @param	filePath	string			配置文件路径
// @return				ConfigReader	构造完毕的 ConfigReader
func NewConfigReader(filePath string) ConfigReader {
	var reader ConfigReader
	reader.filePath = filePath
	reader.parseYml()
	return reader
}

func (reader *ConfigReader) parseYml() {
	buf, err := ioutil.ReadFile(reader.filePath)
	if err != nil {
	}

	result := make(map[string]interface{})
	err = yaml.Unmarshal(buf, &result)
	if err != nil {
	}

	// temp := result ["global"].([]interface{})[0].(map[interface{}]interface{}) ["when"].(map[interface{}]interface{}) ["message"].([]interface{})[0]

	// fmt.Printf("\n\n%v\n\n", temp)
	// fmt.Println(reflect.TypeOf(temp))

	reader.parseToConfigList(result)
}

func (reader *ConfigReader) parseToConfigList(ymlObject map[string]interface{}) {
	nativeGlobal := ymlObject["global"]
	if nativeGlobal != nil {
		for _, item := range nativeGlobal.([]interface{}) {
			nativeWhen := item.(map[interface{}]interface{})["when"].(map[interface{}]interface{})
			nativeDo := item.(map[interface{}]interface{})["do"].(map[interface{}]interface{})
			reader.GlobalConfigList = append(reader.GlobalConfigList,
				reader.parseToConfig(ConfigTypeGlobal, nativeWhen, nativeDo))
		}
	}

	natvieFriend := ymlObject["friend"]
	if natvieFriend != nil {
		for _, item := range natvieFriend.([]interface{}) {
			nativeWhen := item.(map[interface{}]interface{})["when"].(map[interface{}]interface{})
			nativeDo := item.(map[interface{}]interface{})["do"].(map[interface{}]interface{})
			reader.FriendConifgList = append(reader.FriendConifgList, reader.parseToConfig(ConfigTypeFriend, nativeWhen, nativeDo))
		}
	}

	nativeGroup := ymlObject["group"]
	if nativeGroup != nil {
		for _, item := range nativeGroup.([]interface{}) {
			nativeWhen := item.(map[interface{}]interface{})["when"].(map[interface{}]interface{})
			nativeDo := item.(map[interface{}]interface{})["do"].(map[interface{}]interface{})
			reader.GroupConfigList = append(reader.GroupConfigList, reader.parseToConfig(ConfigTypeGroup, nativeWhen, nativeDo))
		}
	}
}

func (reader *ConfigReader) parseToMessageList(nativeMessageList []interface{}) []Message {
	var messageList []Message
	for i := 0; i < len(nativeMessageList); i++ {
		message := reader.parseToMessage(nativeMessageList[i].(map[interface{}]interface{}))
		messageList = append(messageList, message)
	}
	return messageList
}

func (reader *ConfigReader) parseToOperationList(nativeOperationLisst []interface{}) []Operation {
	var operationList []Operation
	for i := 0; i < len(nativeOperationLisst); i++ {
		message := reader.parseToOperation(nativeOperationLisst[i].(map[interface{}]interface{}))
		operationList = append(operationList, message)
	}
	return operationList
}

func (reader *ConfigReader) parseToConfig(configType ConfigType,
	nativeWhen map[interface{}]interface{},
	nativeDo map[interface{}]interface{}) Config {

	nativeSenderList := nativeWhen["sender"]
	var senderList []Sender
	if nativeSenderList != nil {
		senderList = reader.parseToSenderOrReciverList(nativeSenderList.(map[interface{}]interface{}))
	}

	nativeWhenMessageList := nativeWhen["message"]
	var whenMessageList []Message
	if nativeWhenMessageList != nil {
		whenMessageList = reader.parseToMessageList(nativeWhenMessageList.([]interface{}))
	}

	nativeWhenOperation := nativeWhen["operation"]
	var whenOperationList []Operation
	if nativeWhenOperation != nil {
		whenOperationList = reader.parseToOperationList(nativeWhenOperation.([]interface{}))
	}

	nativeReceiverList := nativeDo["receiver"]
	var receiverList []Sender
	if nativeReceiverList != nil {
		receiverList = reader.parseToSenderOrReciverList(nativeReceiverList.(map[interface{}]interface{}))
	}

	nativeDoMessageList := nativeDo["message"]
	var doMessageList []Message
	if nativeDoMessageList != nil {
		doMessageList = reader.parseToMessageList(nativeDoMessageList.([]interface{}))
	}

	nativeDoOperation := nativeDo["operation"]
	var doOperationList []Operation
	if nativeDoOperation != nil {
		doOperationList = reader.parseToOperationList(nativeDoOperation.([]interface{}))
	}
	config, err := NewConfig(configType, whenMessageList, whenOperationList,
		doMessageList, doOperationList, senderList, receiverList)
	if err != nil {

	}
	return config
}

func (reader *ConfigReader) parseToSenderOrReciverList(nativeSender map[interface{}]interface{}) []Sender {
	var senderOrReciverList []Sender

	groupList := nativeSender["group"]
	if groupList != nil {
		for _, groupID := range groupList.([]interface{}) {
			sender, err := NewSender(SenderTypeGroup, int64(groupID.(int)), "", "")
			if err != nil {
				return nil
			}
			senderOrReciverList = append(senderOrReciverList, sender)
		}
	}

	userList := nativeSender["user"]
	if userList != nil {
		for _, userID := range userList.([]interface{}) {
			sender, err := NewSender(SenderTypeUser, int64(userID.(int)), "", "")
			if err != nil {
				return nil
			}
			senderOrReciverList = append(senderOrReciverList, sender)
		}
	}
	return senderOrReciverList
}

// func (reader *ConfigReader) parseToSender(nativeSender []interface{}) Sender {
// 	var sender Sender
// 	return sender
// }

func (reader *ConfigReader) parseToMessage(nativeMessage map[interface{}]interface{}) Message {
	var messageType MessageType
	switch nativeMessage["type"] {
	case "Plain":
		messageType = MessageTypePlain
	case "Image":
		messageType = MessageTypeImage
	case "Face":
		messageType = MessageTypeFace
	case "Event":
		messageType = MessageTypeEvent
	case "Xml":
		messageType = MessageTypeXML
	}

	msgValue := make(map[string]string)

	for key, nativeValue := range nativeMessage {
		value := ""
		switch nativeValue.(type) {
		case string:
			value = nativeValue.(string)
		case int64:
			value = strconv.FormatInt(nativeValue.(int64), 10)
		case bool:
			value = strconv.FormatBool(nativeValue.(bool))
		}
		msgValue[key.(string)] = value
	}

	buf, err := ioutil.ReadFile("../../config/face-map.yml")
	if err != nil {
	}

	faceMap := make(map[string]interface{})
	err = yaml.Unmarshal(buf, &faceMap)
	if err != nil {
	}

	if messageType == MessageTypeFace {
		msgValue["id"] = strconv.Itoa(faceMap[msgValue["name"]].(int))
	}

	msg, err := NewMessage(messageType, msgValue)
	if err != nil {

	}
	return msg
}

func (reader *ConfigReader) parseToOperation(natvieOperation map[interface{}]interface{}) Operation {
	var operationType OperationType
	operationType = CastConfigOperationTypeToOperationType(natvieOperation["type"].(string))

	operationValue := make(map[string]string)
	for key, value := range natvieOperation {
		operationValue[key.(string)] = value.(string)
	}

	operation, err := NewOperation(operationType, operationValue)
	if err != nil {

	}
	return operation
}
