package eltype

import (
	"strconv"

	"gopkg.in/yaml.v2"

	// "io"
	// "fmt"
	// "os"
	"io/ioutil"
	// "reflect"
)

type ConfigReader struct {
	filePath         string
	GlobalConfigList []Config
	FriendConifgList []Config
	GroupConfigList  []Config
}

func NewConfigReader(filePath string) ConfigReader {
	var reader ConfigReader
	reader.filePath = filePath
	reader.ParseYml()
	return reader
}

func (reader *ConfigReader) ParseYml() {
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
	config, err := NewConfig(configType, whenMessageList, whenOperationList, doMessageList, doOperationList)
	if err != nil {

	}
	return config
}

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
	}

	msggValue := make(map[string]string)

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
		msggValue[key.(string)] = value
	}

	msg, err := NewMessage(messageType, msggValue)
	if err != nil {

	}
	return msg
}

func (reader *ConfigReader) parseToOperation(natvieOperation map[interface{}]interface{}) Operation {
	var operationType OperationType
	switch natvieOperation["type"] {
	case "mute":
		operationType = OperationTypeMute
	}

	operationValue := make(map[string]string)

	for key, value := range natvieOperation {
		operationValue[key.(string)] = value.(string)
	}

	operation, err := NewOperation(operationType, operationValue)
	if err != nil {

	}
	return operation
}
