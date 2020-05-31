package eltype

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"

	"gopkg.in/yaml.v2"
)

type Listen struct {
	GroupIDList []string `yaml:"group"`
	UserIDList  []string `yaml:"user"`
}

type Target struct {
	GroupIDList []string `yaml:"group"`
	UserIDList  []string `yaml:"user"`
}

type Transfer struct {
	Listen Listen `yaml:"listen"`
	Target Target `yaml:"target"`
}

type SourceConfig struct {
	GlobalConfigList []Config   `yaml:"global"`
	FriendConfigList []Config   `yaml:"friend"`
	GroupConfigList  []Config   `yaml:"group"`
	CronConfigList   []Config   `yaml:"crontab"`
	RssConfigList    []Config   `yaml:"rss"`
	TransferList     []Transfer `yaml:"transfer"`
}

type Compiler struct {
	SourceConfig SourceConfig
	folder       string
}

func NewCompiler(folder string) (Compiler, error) {
	var compiler Compiler
	compiler.folder = folder
	return compiler, nil
}

func (compiler *Compiler) Compile() {
	if compiler.folder == "default" {
		compiler.compileThisFile(ConfigRoot + "/" + DefaultConfigFileName)
	} else {
		compiler.compileFolder()
	}

	compiler.CompleteConfigList()
}

func (compiler *Compiler) WriteFile() string {
	ymlStr, err := yaml.Marshal(compiler.SourceConfig)
	if err != nil {
		return ""
	}
	err = ioutil.WriteFile(ConfigRoot+"/obj/obj.yml", ymlStr, 0777)
	if err != nil {
		return ""
	}
	return ConfigRoot + "/obj/obj.yml"
}

func (compiler *Compiler) compileFolder() {
	files, err := ioutil.ReadDir(ConfigRoot + "/" + compiler.folder)
	if err != nil {
		return
	}
	for _, fileInfo := range files {
		isMatch, err := regexp.MatchString(".+\\.yml", fileInfo.Name())
		if isMatch && err == nil {
			compiler.compileThisFile(fmt.Sprintf("%s/%s/%s", ConfigRoot, compiler.folder, fileInfo.Name()))
		}
	}
}

func (compiler *Compiler) compileThisFile(filePath string) {
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("跳过 %s, 因为未能读取文件。\n", filePath)
		return
	}
	var tempSourceConfig SourceConfig
	yaml.Unmarshal(buf, &tempSourceConfig)

	for _, transfer := range tempSourceConfig.TransferList {
		compiler.SourceConfig.GlobalConfigList = append(compiler.SourceConfig.GlobalConfigList, transfer.toConfig())
	}

	compiler.mergeSourceConfig(tempSourceConfig)
}

func (compiler *Compiler) CompleteConfigList() {
	var innerID int = 0
	for i := 0; i < len(compiler.SourceConfig.GlobalConfigList); i++ {
		temp := compiler.SourceConfig.GlobalConfigList[i]
		temp.Init()
		temp.innerID = innerID
		innerID++
		compiler.SourceConfig.GlobalConfigList[i] = temp
	}
	for i := 0; i < len(compiler.SourceConfig.FriendConfigList); i++ {
		temp := compiler.SourceConfig.FriendConfigList[i]
		temp.Init()
		temp.innerID = innerID
		innerID++
		compiler.SourceConfig.FriendConfigList[i] = temp
	}
	for i := 0; i < len(compiler.SourceConfig.GroupConfigList); i++ {
		temp := compiler.SourceConfig.GroupConfigList[i]
		temp.Init()
		temp.innerID = innerID
		innerID++
		compiler.SourceConfig.GroupConfigList[i] = temp
	}
	for i := 0; i < len(compiler.SourceConfig.CronConfigList); i++ {
		temp := compiler.SourceConfig.CronConfigList[i]
		temp.Init()
		temp.innerID = innerID
		innerID++
		compiler.SourceConfig.CronConfigList[i] = temp
	}
}

func (compiler *Compiler) mergeSourceConfig(tempSourceConfig SourceConfig) {
	mergeConfigList(&compiler.SourceConfig.GlobalConfigList, tempSourceConfig.GlobalConfigList)
	mergeConfigList(&compiler.SourceConfig.FriendConfigList, tempSourceConfig.FriendConfigList)
	mergeConfigList(&compiler.SourceConfig.GroupConfigList, tempSourceConfig.GroupConfigList)
	mergeConfigList(&compiler.SourceConfig.CronConfigList, tempSourceConfig.CronConfigList)
	mergeConfigList(&compiler.SourceConfig.RssConfigList, tempSourceConfig.RssConfigList)
}

func (transfer *Transfer) toConfig() Config {
	var config Config
	for _, groupID := range transfer.Listen.GroupIDList {
		config.When.Message.Sender.AddGroupID(groupID)
	}
	for _, UserID := range transfer.Listen.UserIDList {
		config.When.Message.Sender.AddUserID(UserID)
	}
	for _, groupID := range transfer.Target.GroupIDList {
		config.Do.Message.Receiver.AddGroupID(groupID)
	}
	for _, UserID := range transfer.Target.UserIDList {
		config.Do.Message.Receiver.AddUserID(UserID)
	}

	var messageDetail MessageDetail
	messageDetail.innerType = MessageTypePlain
	messageDetail.Regex = "(?:.|\\n)+"
	config.When.Message.AddDetail(messageDetail)
	messageDetail.Regex = ""

	messageDetail.innerType = MessageTypeImage
	config.When.Message.AddDetail(messageDetail)

	messageDetail.innerType = MessageTypeXML
	config.When.Message.AddDetail(messageDetail)

	messageDetail.innerType = MessageTypePlain
	messageDetail.Text = "{el-message-text}"
	config.Do.Message.AddDetail(messageDetail)
	messageDetail.Text = ""

	messageDetail.innerType = MessageTypeImage
	for i := 0; i < 20; i++ {
		messageDetail.URL = fmt.Sprintf("{el-message-image-url-%d}", i)
		config.Do.Message.AddDetail(messageDetail)
	}
	messageDetail.URL = ""

	messageDetail.innerType = MessageTypeXML
	messageDetail.Text = "{el-message-xml}"
	config.Do.Message.AddDetail(messageDetail)
	messageDetail.Text = ""

	return config
}
