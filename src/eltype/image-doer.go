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
		for _, doMessage := range config.DoMessageList {
			if doMessage.Type == MessageTypeImage {
				value := make(map[string]string)
				if doMessage.Value["url"] != "" {
					if doMessage.Value["direct"] == "true" {
						filename, err := doer.downloadImage(doMessage.Value["url"])
						if err != nil {
							continue
						}
						value["path"] = filename
					} else {
						value["url"] = doer.replaceStrByPreDefVarMap(doMessage.Value["url"])
					}

				} else if doMessage.Value["path"] != "" {
					value = doMessage.Value
				}
				message, err := NewMessage(MessageTypeImage, value)
				if err != nil {
					continue
				}
				doer.sendedMessageList = append(doer.sendedMessageList, message)
			}
		}
	}
}

func (doer *ImageDoer) downloadImage(url string) (string, error) {
	filename := strconv.FormatInt(rand.Int63(), 10) + ".jpg"
	file, err := os.Create("../../plugins/MiraiAPIHTTP/images/" + filename)
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

func (doer ImageDoer) replaceStrByPreDefVarMap(text string) string {
	for varName, value := range doer.preDefVarMap {
		key := fmt.Sprintf("{%s}", varName)
		text = strings.ReplaceAll(text, key, value)
	}
	return text
}

// GetSendedMessageList 获取将要发送的信息列表
func (doer ImageDoer) GetSendedMessageList() []Message {
	return doer.sendedMessageList
}

// GetSendedOperationList 获取将要执行的动作列表
func (doer ImageDoer) GetSendedOperationList() []Operation {
	return doer.sendedOperationList
}
