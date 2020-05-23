package eltype

import (
	"fmt"
	"strconv"

	"gopkg.in/yaml.v2"

	// "io"

	// "os"
	"io/ioutil"
	// "reflect"
)

var SettingFullPath string = "../../plugins/MiraiAPIHTTP/setting.yml"
var FaceMapFullPath string = "../../config/face-map.yml"
var DefaultConfigFullPath string = "../../config/default.yml"

// ConfigReader 配置读取对象
type ConfigReader struct {
	Port              string
	EnableWebsocket   bool
	folder            string
	AuthKey           string
	GlobalConfigList  []Config
	FriendConifgList  []Config
	GroupConfigList   []Config
	CrontabConfigList []Config
	CounterConfigList []Config
}

// NewConfigReader 使用配置文件路径构造一个 ConfigReader
// @param	filePath	string			配置文件路径
func NewConfigReader(folder string) ConfigReader {
	var reader ConfigReader
	reader.folder = folder
	reader.parseYml()
	// fmt.Printf("%v\n", reader.GroupConfigList)
	return reader
}

func (reader *ConfigReader) parseYml() {
	reader.parseToSetting()

	if reader.folder == "" {
		// fmt.Println("使用默认配置，不使用" + reader.folder + "\n")
		reader.parseThisFile(DefaultConfigFullPath)
	} else {
		files, err := ioutil.ReadDir(reader.folder)
		if err != nil {
			// fmt.Println("使用默认配置，不使用" + reader.folder + "\n")
			reader.parseThisFile(DefaultConfigFullPath)
		}

		// fmt.Println("使用自定义配置: " + reader.folder + "\n")
		for _, file := range files {
			if !file.IsDir() {
				// fmt.Printf("正在读取配置：%s/%s\n", reader.folder, file.Name())
				reader.parseThisFile(fmt.Sprintf("%s/%s", reader.folder, file.Name()))
			}
		}
	}
}

func (reader *ConfigReader) parseThisFile(fileFullPath string) {
	buf, err := ioutil.ReadFile(fileFullPath)
	if err != nil {
		fmt.Printf("跳过 %s, 因为未能打开文件。\n", fileFullPath)
		return
	}
	result := make(map[string]interface{})
	err = yaml.Unmarshal(buf, &result)
	if err != nil {
		fmt.Printf("跳过 %s, 因为解析失败，配置文件可能存在语法错误。\n", fileFullPath)
		return
	}
	reader.parseToConfigList(result)
}

func (reader *ConfigReader) parseToSetting() {
	buf, err := ioutil.ReadFile(SettingFullPath)
	if err != nil {
		fmt.Printf("跳过 %s, 因为未能打开文件。\n", SettingFullPath)
		return
	}
	result := make(map[string]interface{})
	err = yaml.Unmarshal(buf, &result)
	if err != nil {
		fmt.Printf("跳过 %s, 因为解析失败，配置文件可能存在语法错误。\n", SettingFullPath)
		return
	}
	reader.Port = strconv.Itoa(result["port"].(int))
	reader.AuthKey = result["authKey"].(string)
	reader.EnableWebsocket = result["enableWebsocket"].(bool)
}

func (reader *ConfigReader) parseToConfigList(ymlObject map[string]interface{}) {
	nativeCrontab := ymlObject["crontab"]
	// fmt.Printf("%v", nativeCrontab)
	if nativeCrontab != nil {
		for _, item := range nativeCrontab.([]interface{}) {
			cron := item.(map[interface{}]interface{})["cron"].(string)
			nativeDo := item.(map[interface{}]interface{})["do"].(map[interface{}]interface{})
			nativeReceiverList := nativeDo["receiver"]
			var receiverList []Receiver
			if nativeReceiverList != nil {
				receiverList = reader.parseToReceiverList(nativeReceiverList.(map[interface{}]interface{}))
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

			config, err := NewConfig(ConfigTypeCrontab, nil, nil, doMessageList,
				doOperationList, nil, receiverList, cron, false, "")
			if err != nil {
				continue
			}
			reader.CrontabConfigList = append(reader.CrontabConfigList, config)
		}
	}

	nativeGlobal := ymlObject["global"]
	if nativeGlobal != nil {
		for _, nativeConfig := range nativeGlobal.([]interface{}) {
			reader.GlobalConfigList = append(reader.GlobalConfigList,
				reader.parseToConfig(ConfigTypeGlobal, nativeConfig))
		}
	}

	natvieFriend := ymlObject["friend"]
	if natvieFriend != nil {
		for _, nativeConfig := range natvieFriend.([]interface{}) {
			reader.GlobalConfigList = append(reader.GlobalConfigList,
				reader.parseToConfig(ConfigTypeFriend, nativeConfig))
		}
	}

	nativeGroup := ymlObject["group"]
	if nativeGroup != nil {
		for _, nativeConfig := range nativeGroup.([]interface{}) {
			reader.GroupConfigList = append(reader.GroupConfigList,
				reader.parseToConfig(ConfigTypeGroup, nativeConfig))
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

func (reader *ConfigReader) parseToConfig(configType ConfigType, nativeConfig interface{}) Config {

	nativeWhen := nativeConfig.(map[interface{}]interface{})["when"].(map[interface{}]interface{})
	nativeDo := nativeConfig.(map[interface{}]interface{})["do"].(map[interface{}]interface{})
	var isCount bool
	if nativeDo["count"] == "" || nativeDo["count"] == nil {
		isCount = false
	} else {
		isCount = nativeDo["count"].(bool)
	}
	var countID string
	if nativeConfig.(map[interface{}]interface{})["countID"] == "" ||
		nativeConfig.(map[interface{}]interface{})["countID"] == nil {
		countID = ""
	} else {
		countID = nativeConfig.(map[interface{}]interface{})["countID"].(string)
	}

	nativeSenderList := nativeWhen["sender"]
	var senderList []Sender
	if nativeSenderList != nil {
		senderList = reader.parseToSenderList(nativeSenderList.(map[interface{}]interface{}))
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
	var receiverList []Receiver
	if nativeReceiverList != nil {
		receiverList = reader.parseToReceiverList(nativeReceiverList.(map[interface{}]interface{}))
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
		doMessageList, doOperationList, senderList, receiverList, "", isCount, countID)
	if err != nil {

	}
	return config
}

func (reader *ConfigReader) parseToSenderList(nativeSender map[interface{}]interface{}) []Sender {
	var senderList []Sender

	groupList := nativeSender["group"]
	if groupList != nil {
		for _, groupID := range groupList.([]interface{}) {
			sender, err := NewSender(SenderTypeGroup, int64(groupID.(int)), "", "")
			if err != nil {
				return nil
			}
			senderList = append(senderList, sender)
		}
	}

	friendList := nativeSender["friend"]
	if friendList != nil {
		for _, userID := range friendList.([]interface{}) {
			sender, err := NewSender(SenderTypeFriend, int64(userID.(int)), "", "")
			if err != nil {
				return nil
			}
			senderList = append(senderList, sender)
		}
	}
	return senderList
}

func (reader *ConfigReader) parseToReceiverList(nativeSender map[interface{}]interface{}) []Receiver {
	var reciverList []Receiver

	groupList := nativeSender["group"]
	if groupList != nil {
		for _, groupID := range groupList.([]interface{}) {
			receiver, err := NewReceiver(ReceiverTypeGroup, int64(groupID.(int)), "", "")
			if err != nil {
				return nil
			}
			reciverList = append(reciverList, receiver)
		}
	}

	userList := nativeSender["user"]
	if userList != nil {
		for _, userID := range userList.([]interface{}) {
			receiver, err := NewReceiver(ReceiverTypeUser, int64(userID.(int)), "", "")
			if err != nil {
				return nil
			}
			reciverList = append(reciverList, receiver)
		}
	}
	return reciverList
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
	case "Xml":
		messageType = MessageTypeXML
	}

	msgValue := make(map[string]string)

	for key, nativeValue := range nativeMessage {
		value := ""
		switch nativeValue.(type) {
		case string:
			value = nativeValue.(string)
		case int:
			value = strconv.Itoa(nativeValue.(int))
		case int64:
			value = strconv.FormatInt(nativeValue.(int64), 10)
		case bool:
			value = strconv.FormatBool(nativeValue.(bool))
		}
		msgValue[key.(string)] = value
	}

	buf, err := ioutil.ReadFile(FaceMapFullPath)
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
