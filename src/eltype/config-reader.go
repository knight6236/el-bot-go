package eltype

import (
	"log"
	"strconv"

	// "os"
	"io/ioutil"

	"gopkg.in/yaml.v2"
	// "reflect"
)

// ConfigReader 配置读取对象
type ConfigReader struct {
	FreqUpperLimit   int64 `yaml:"feqLimit"`
	Port             string
	EnableWebsocket  bool
	folder           string
	AuthKey          string
	Compiler         *Compiler
	GlobalConfigList []Config `yaml:"global"`
	FriendConfigList []Config `yaml:"friend"`
	GroupConfigList  []Config `yaml:"group"`
	CronConfigList   []Config `yaml:"crontab"`
	RssConfigList    []Config `yaml:"rss"`
}

// NewConfigReader 使用配置文件路径构造一个 ConfigReader
// @param	filePath	string			配置文件路径
func NewConfigReader(folder string) (*ConfigReader, error) {
	reader := new(ConfigReader)
	reader.folder = folder
	return reader, nil
}

func (reader *ConfigReader) Load(isDebug bool) {
	var err error
	reader.Compiler, err = NewCompiler(reader.folder)
	if err != nil {
		return
	}
	reader.Compiler.Compile()
	if isDebug {
		filePath := reader.Compiler.WriteFile()
		reader.parseThisFile(filePath)
		reader.CompleteConfigList()
	} else {
		reader.FreqUpperLimit = reader.Compiler.SourceConfig.FreqUpperLimit
		reader.GlobalConfigList = reader.Compiler.SourceConfig.GlobalConfigList
		reader.FriendConfigList = reader.Compiler.SourceConfig.FriendConfigList
		reader.GroupConfigList = reader.Compiler.SourceConfig.GroupConfigList
		reader.CronConfigList = reader.Compiler.SourceConfig.CronConfigList
		reader.RssConfigList = reader.Compiler.SourceConfig.RssConfigList
	}
	reader.parseToSetting()
}

// reLoad 重新加载配置
func (reader *ConfigReader) reLoad() {
	reader.GlobalConfigList = reader.GlobalConfigList[:0]
	reader.FriendConfigList = reader.FriendConfigList[:0]
	reader.GroupConfigList = reader.GroupConfigList[:0]
	reader.CronConfigList = reader.CronConfigList[:0]
	reader.RssConfigList = reader.RssConfigList[:0]
	reader.Load(false)
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
	yaml.Unmarshal(buf, reader)
}

func (reader *ConfigReader) CompleteConfigList() {
	var InnerID int64 = 1
	for i := 0; i < len(reader.GlobalConfigList); i++ {
		temp := reader.GlobalConfigList[i]
		temp.CompleteType()
		temp.InnerID = InnerID
		InnerID++
		reader.GlobalConfigList[i] = temp
	}
	for i := 0; i < len(reader.FriendConfigList); i++ {
		temp := reader.FriendConfigList[i]
		temp.CompleteType()
		temp.InnerID = InnerID
		InnerID++
		reader.FriendConfigList[i] = temp
	}
	for i := 0; i < len(reader.GroupConfigList); i++ {
		temp := reader.GroupConfigList[i]
		temp.CompleteType()
		temp.InnerID = InnerID
		InnerID++
		reader.GroupConfigList[i] = temp
	}
	for i := 0; i < len(reader.CronConfigList); i++ {
		temp := reader.CronConfigList[i]
		temp.CompleteType()
		temp.InnerID = InnerID
		InnerID++
		reader.CronConfigList[i] = temp
	}
}

func (reader *ConfigReader) mergeReader(tempReader ConfigReader) {
	reader.FreqUpperLimit = tempReader.FreqUpperLimit
	MergeConfigList(&reader.GlobalConfigList, tempReader.GlobalConfigList)
	MergeConfigList(&reader.FriendConfigList, tempReader.FriendConfigList)
	MergeConfigList(&reader.GroupConfigList, tempReader.GroupConfigList)
	MergeConfigList(&reader.CronConfigList, tempReader.CronConfigList)
}
