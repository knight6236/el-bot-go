package eltype

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"runtime"

	"gopkg.in/yaml.v2"
)

type Listen struct {
	GroupIDList []string `yaml:"group"`
	// UserIDList  []string `yaml:"user"`
}

type Target struct {
	GroupIDList []string `yaml:"group"`
	UserIDList  []string `yaml:"user"`
}

type Transfer struct {
	Listen Listen `yaml:"listen"`
	Target Target `yaml:"target"`
}

// type Echo struct {
// 	Enable bool `yaml:"enable"`
// }

type SourceConfig struct {
	// Echo             Echo           `yaml:"echo" json:"echo"`
	FreqUpperLimit   int64      `yaml:"freqLimit" json:"freqLimit"`
	GlobalConfigList []Config   `yaml:"global" json:"global"`
	FriendConfigList []Config   `yaml:"friend" json:"friend"`
	GroupConfigList  []Config   `yaml:"group" json:"group"`
	CronConfigList   []Config   `yaml:"crontab" json:"crontab"`
	RssConfigList    []Config   `yaml:"rss" json:"rss"`
	TransferList     []Transfer `yaml:"transfer" json:"transfer"`
}

type Compiler struct {
	pluginReader *PluginReader
	SourceConfig SourceConfig
	folder       string
}

func NewCompiler(folder string) (Compiler, error) {
	var compiler Compiler
	compiler.pluginReader, _ = NewPluginReader()
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

func (compiler *Compiler) callPlugin(configMap map[string]interface{}) {
	for _, plugin := range compiler.pluginReader.PluginMap {
		obj := configMap[plugin.ConfigKeyword]
		if obj == nil {
			continue
		}
		jsonMap := JsonParse(obj.(map[interface{}]interface{}), 0)

		jsonStr, err := json.Marshal(jsonMap)
		fmt.Println(string(jsonStr))
		var ret string
		switch plugin.Type {
		case PluginTypeBinary:
			if runtime.GOOS == "windows" {
				ret, err = ExecCommand(plugin.Path, string(jsonStr))
			} else {
				ret, err = ExecCommand("/bin/bash", "-c", string(jsonStr))
			}
		case PluginTypeJava:
			if runtime.GOOS == "windows" {
				ret, err = ExecCommand("java", "-jar", plugin.Path, string(jsonStr))
			} else {
				ret, err = ExecCommand("/bin/bash", "-c", fmt.Sprintf("java -jar %s %s", plugin.Path, string(jsonStr)))
			}
		case PluginTypePython:
			if runtime.GOOS == "windows" {
				ret, err = ExecCommand("python", plugin.Path, string(jsonStr))
			} else {
				ret, err = ExecCommand("/bin/bash", "-c", fmt.Sprintf("%s %s %s", PythonCommand, plugin.Path, string(jsonStr)))
			}
		case PluginTypeJavaScript:
			if runtime.GOOS == "windows" {
				ret, err = ExecCommand("node", plugin.Path, string(jsonStr))
			} else {
				ret, err = ExecCommand("/bin/bash", "-c", fmt.Sprintf("node %s %s", plugin.Path, string(jsonStr)))
			}
		default:
			return
		}

		if err != nil {
			continue
		}
		fmt.Println(ret)
		var tempSourceConfig SourceConfig
		json.Unmarshal([]byte(ret), &tempSourceConfig)
		compiler.mergeSourceConfig(tempSourceConfig)
	}
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
	var configMap map[string]interface{}
	yaml.Unmarshal(buf, &configMap)
	compiler.callPlugin(configMap)

	for _, transfer := range tempSourceConfig.TransferList {
		compiler.SourceConfig.GroupConfigList = append(compiler.SourceConfig.GroupConfigList, transfer.toConfig())
	}

	// compiler.SourceConfig.GlobalConfigList = append(compiler.SourceConfig.GlobalConfigList, compiler.SourceConfig.Echo.toConfig())

	compiler.mergeSourceConfig(tempSourceConfig)
}

func (compiler *Compiler) CompleteConfigList() {
	var InnerID int64 = 1
	for i := 0; i < len(compiler.SourceConfig.GlobalConfigList); i++ {
		temp := compiler.SourceConfig.GlobalConfigList[i]
		temp.CompleteType()
		temp.InnerID = InnerID
		InnerID++
		compiler.SourceConfig.GlobalConfigList[i] = temp
	}
	for i := 0; i < len(compiler.SourceConfig.FriendConfigList); i++ {
		temp := compiler.SourceConfig.FriendConfigList[i]
		temp.CompleteType()
		temp.InnerID = InnerID
		InnerID++
		compiler.SourceConfig.FriendConfigList[i] = temp
	}
	for i := 0; i < len(compiler.SourceConfig.GroupConfigList); i++ {
		temp := compiler.SourceConfig.GroupConfigList[i]
		temp.CompleteType()
		temp.InnerID = InnerID
		InnerID++
		compiler.SourceConfig.GroupConfigList[i] = temp
	}
	for i := 0; i < len(compiler.SourceConfig.CronConfigList); i++ {
		temp := compiler.SourceConfig.CronConfigList[i]
		temp.CompleteType()
		temp.InnerID = InnerID
		InnerID++
		compiler.SourceConfig.CronConfigList[i] = temp
	}
}

func (compiler *Compiler) mergeSourceConfig(tempSourceConfig SourceConfig) {
	compiler.SourceConfig.FreqUpperLimit = tempSourceConfig.FreqUpperLimit
	// compiler.SourceConfig.Echo = tempSourceConfig.Echo
	MergeConfigList(&compiler.SourceConfig.GlobalConfigList, tempSourceConfig.GlobalConfigList)
	MergeConfigList(&compiler.SourceConfig.FriendConfigList, tempSourceConfig.FriendConfigList)
	MergeConfigList(&compiler.SourceConfig.GroupConfigList, tempSourceConfig.GroupConfigList)
	MergeConfigList(&compiler.SourceConfig.CronConfigList, tempSourceConfig.CronConfigList)
	MergeConfigList(&compiler.SourceConfig.RssConfigList, tempSourceConfig.RssConfigList)
}

func (transfer *Transfer) toConfig() Config {
	var config Config
	for _, groupID := range transfer.Listen.GroupIDList {
		config.When.Message.Sender.AddGroupID(groupID)
	}
	// for _, UserID := range transfer.Listen.UserIDList {
	// 	config.When.Message.Sender.AddUserID(UserID)
	// }
	for _, groupID := range transfer.Target.GroupIDList {
		config.Do.Message.Receiver.AddGroupID(groupID)
	}
	for _, UserID := range transfer.Target.UserIDList {
		config.Do.Message.Receiver.AddUserID(UserID)
	}

	var messageDetail MessageDetail
	messageDetail.InnerType = MessageTypePlain
	messageDetail.Regex = "(?:.|\\n)+"
	config.When.Message.AddDetail(messageDetail)
	messageDetail.Regex = ""

	messageDetail.InnerType = MessageTypeImage
	config.When.Message.AddDetail(messageDetail)

	messageDetail.InnerType = MessageTypeXML
	config.When.Message.AddDetail(messageDetail)

	messageDetail.InnerType = MessageTypePlain
	messageDetail.Text = "{el-message-text}"
	config.Do.Message.AddDetail(messageDetail)
	messageDetail.Text = ""

	messageDetail.InnerType = MessageTypeImage
	for i := 0; i < 20; i++ {
		messageDetail.URL = fmt.Sprintf("{el-message-image-url-%d}", i)
		config.Do.Message.AddDetail(messageDetail)
	}
	messageDetail.URL = ""

	messageDetail.InnerType = MessageTypeXML
	messageDetail.Text = "{el-message-xml}"
	config.Do.Message.AddDetail(messageDetail)
	messageDetail.Text = ""

	return config
}

// func (echo *Echo) toConfig() Config {
// 	var config Config
// 	messageDetail := MessageDetail{
// 		InnerType: MessageTypePlain,
// 		Regex:     "echo\\s(.+)",
// 	}
// 	config.When.Message.AddDetail(messageDetail)
// 	messageDetail.Regex = ""
// 	messageDetail.Text = "{el-regex-0}"
// 	config.Do.Message.AddDetail(messageDetail)
// 	config.Do.Message.IsQuote = true
// 	return config
// }
