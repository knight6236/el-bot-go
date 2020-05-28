package eltype

import (
	"fmt"
	"log"
	"strconv"

	"gopkg.in/yaml.v2"

	// "io"

	// "os"
	"io/ioutil"
	// "reflect"
)

// ConfigReader 配置读取对象
type ConfigReader struct {
	Port              string
	EnableWebsocket   bool
	folder            string
	AuthKey           string
	GlobalConfigList  []Config `yaml:"global"`
	FriendConfigList  []Config `yaml:"friend"`
	GroupConfigList   []Config `yaml:"group"`
	CrontabConfigList []Config `yaml:"crontab"`
	// CounterConfigList []Config
}

// NewConfigReader 使用配置文件路径构造一个 ConfigReader
// @param	filePath	string			配置文件路径
func NewConfigReader(folder string) ConfigReader {
	var reader ConfigReader
	reader.folder = folder
	reader.parseToSetting()
	reader.parseYml()
	reader.initConfigList()
	// go reader.monitorFolder()
	// fmt.Printf("%v\n", reader.GroupConfigList)
	return reader
}

// reLoad 重新加载配置
func (reader *ConfigReader) reLoad() {
	reader.GlobalConfigList = reader.GlobalConfigList[:0]
	reader.FriendConfigList = reader.FriendConfigList[:0]
	reader.GroupConfigList = reader.GroupConfigList[:0]
	reader.CrontabConfigList = reader.CrontabConfigList[:0]
	// reader.CounterConfigList = reader.CounterConfigList[:0]
	reader.parseYml()
	reader.initConfigList()
}

func (reader *ConfigReader) parseYml() {
	// reader.parseToSetting()

	files, err := ioutil.ReadDir(reader.folder)
	if err != nil {
	}

	// fmt.Println("使用自定义配置: " + reader.folder + "\n")
	for _, file := range files {
		if !file.IsDir() {
			// fmt.Printf("正在读取配置：%s/%s\n", reader.folder, file.Name())
			reader.parseThisFile(fmt.Sprintf("%s/%s", reader.folder, file.Name()))
		}
	}

}

func (reader *ConfigReader) initConfigList() {
	for i := 0; i < len(reader.GlobalConfigList); i++ {
		temp := reader.GlobalConfigList[i]
		temp.Init()
		reader.GlobalConfigList[i] = temp
	}
	for i := 0; i < len(reader.FriendConfigList); i++ {
		temp := reader.FriendConfigList[i]
		temp.Init()
		reader.FriendConfigList[i] = temp
	}
	for i := 0; i < len(reader.GroupConfigList); i++ {
		temp := reader.GroupConfigList[i]
		temp.Init()
		reader.GroupConfigList[i] = temp
	}
	for i := 0; i < len(reader.CrontabConfigList); i++ {
		temp := reader.CrontabConfigList[i]
		temp.Init()
		reader.CrontabConfigList[i] = temp
	}
}

func (reader *ConfigReader) parseToSetting() {
	buf, err := ioutil.ReadFile(SettingFullPath)
	if err != nil {
		log.Printf("跳过 %s, 因为未能打开文件。\n", SettingFullPath)
		return
	}
	result := make(map[string]interface{})
	err = yaml.Unmarshal(buf, &result)
	if err != nil {
		log.Printf("跳过 %s, 因为解析失败，配置文件可能存在语法错误。\n", SettingFullPath)
		return
	}
	reader.Port = strconv.Itoa(result["port"].(int))
	reader.AuthKey = result["authKey"].(string)
	reader.EnableWebsocket = result["enableWebsocket"].(bool)
}

func (reader *ConfigReader) parseThisFile(fileFullPath string) {
	buf, err := ioutil.ReadFile(fileFullPath)
	if err != nil {
		log.Printf("跳过 %s, 因为未能打开文件。\n", fileFullPath)
		return
	}
	result := make(map[string]interface{})
	err = yaml.Unmarshal(buf, &result)
	if err != nil {
		log.Printf("跳过 %s, 因为解析失败，配置文件可能存在语法错误。\n", fileFullPath)
		return
	}
	var tempReader ConfigReader
	yaml.Unmarshal(buf, &tempReader)
	reader.mergeReader(tempReader)
}

func (reader *ConfigReader) mergeReader(tempReader ConfigReader) {
	mergeConfigList(&reader.GlobalConfigList, tempReader.GlobalConfigList)
	mergeConfigList(&reader.FriendConfigList, tempReader.FriendConfigList)
	mergeConfigList(&reader.GroupConfigList, tempReader.GroupConfigList)
}
