package eltype

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// ImageDoer 表情动作生成类
// @property	configHitList		[]Config			命中的配置列表
// @property	recivedMessageList	[]Message			接收到的消息列表
// @property	sendedMessageList	[]Message			将要发送的消息列表
// @property	preDefVarMap		map[string]string	预定义变量Map
type ImageDoer struct {
	configHitList       []Config
	recivedMessageList  []Message
	sendedMessageList   []Message
	sendedOperationList []Operation
	preDefVarMap        map[string]string
}

// NewImageDoer 构造一个 ImageDoer
// @param	configHitList		[]Config			命中的配置列表
// @param	recivedMessageList	[]Message			接收到的消息列表
// @param	sendedMessageList	[]Message			将要发送的消息列表
// @property	sendedOperationList	[]Operation			将要执行的动作列表
// @param	preDefVarMap		map[string]string	预定义变量Map
func NewImageDoer(configHitList []Config, recivedMessageList []Message, preDefVarMap map[string]string) (IDoer, error) {
	var doer ImageDoer
	doer.configHitList = configHitList
	doer.recivedMessageList = recivedMessageList
	doer.preDefVarMap = preDefVarMap
	doer.getSendedMessageList()
	return doer, nil
}

func (doer *ImageDoer) getSendedMessageList() {
	for _, config := range doer.configHitList {
		for _, doMessageDetail := range config.Do.Message.DetailList {
			var willBeSentMessage Message
			var willBeSentMessageDetail MessageDetail
			willBeSentMessage.Sender = config.Do.Message.Sender.DeepCopy()
			willBeSentMessage.Receiver = config.Do.Message.Receiver.DeepCopy()
			willBeSentMessageDetail.innerType = MessageTypeImage
			if doMessageDetail.innerType == MessageTypeImage {
				if doMessageDetail.URL != "" {
					if doMessageDetail.ReDirect == true {
						filename, err := doer.downloadImage(doMessageDetail.URL)
						if err != nil {
							fmt.Println(err)
							continue
						}
						willBeSentMessageDetail.Path = filename
					} else {
						url, isReplace := doer.replaceStrByPreDefVarMap(doMessageDetail.URL)
						if isReplace {
							willBeSentMessageDetail.URL = url
						} else {
							willBeSentMessageDetail.URL = ""
						}
					}

				} else if doMessageDetail.Path != "" {
					willBeSentMessageDetail.Path = doMessageDetail.Path
				}
				willBeSentMessage.AddDetail(willBeSentMessageDetail)
				doer.sendedMessageList = append(doer.sendedMessageList, willBeSentMessage)
			}
		}
	}
}

func (doer *ImageDoer) downloadImage(url string) (string, error) {
	filename := strconv.FormatInt(rand.Int63(), 10) + ".jpg"
	file, err := os.Create(ImageFolder + "/" + filename)
	if err != nil {
		return "", err
	}

	defer file.Close()

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	io.Copy(file, res.Body)
	return filename, nil
}

func (doer ImageDoer) replaceStrByPreDefVarMap(text string) (string, bool) {
	var isReplace bool = false
	for varName, value := range doer.preDefVarMap {
		key := fmt.Sprintf("{%s}", varName)
		temp := text
		text = strings.ReplaceAll(text, key, value)
		if !isReplace && temp != text {
			isReplace = true
		}
	}
	return text, isReplace
}

// GetSendedMessageList 获取将要发送的信息列表
func (doer ImageDoer) GetSendedMessageList() []Message {
	return doer.sendedMessageList
}

// GetSendedOperationList 获取将要执行的动作列表
func (doer ImageDoer) GetSendedOperationList() []Operation {
	return doer.sendedOperationList
}
