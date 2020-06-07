package eltype

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"runtime"
)

type PluginReader struct {
	os         string
	arch       string
	PluginList []Plugin
	keywordSet map[string]bool
}

func NewPluginReader() (*PluginReader, error) {
	reader := new(PluginReader)
	reader.os = runtime.GOOS
	reader.arch = runtime.GOARCH
	reader.keywordSet = make(map[string]bool)
	reader.ReadAllPlugin()
	return reader, nil
}

func (reader *PluginReader) ReadAllPlugin() {
	switch reader.os {
	case "freebsd":
		reader.ReadFolder(fmt.Sprintf("%s/%s/%s", PlguinFolder, reader.os, reader.arch), Binary)
	case "linux":
		reader.ReadFolder(fmt.Sprintf("%s/%s/%s", PlguinFolder, reader.os, reader.arch), Binary)
	case "windows":
		if reader.arch == "386" {
			reader.ReadFolder(fmt.Sprintf("%s/%s/%s", PlguinFolder, reader.os, "386"), Binary)
		} else {
			reader.ReadFolder(fmt.Sprintf("%s/%s/%s", PlguinFolder, reader.os, "386"), Binary)
			reader.ReadFolder(fmt.Sprintf("%s/%s/%s", PlguinFolder, reader.os, "amd64"), Binary)
		}
	case "darwin":
		if reader.arch == "386" {
			reader.ReadFolder(fmt.Sprintf("%s/%s/%s", PlguinFolder, reader.os, "386"), Binary)
		} else {
			reader.ReadFolder(fmt.Sprintf("%s/%s/%s", PlguinFolder, reader.os, "386"), Binary)
			reader.ReadFolder(fmt.Sprintf("%s/%s/%s", PlguinFolder, reader.os, "amd64"), Binary)
		}
	}
	reader.ReadFolder(fmt.Sprintf("%s/%s", PlguinFolder, "java"), Java)
	reader.ReadFolder(fmt.Sprintf("%s/%s", PlguinFolder, "javaScript"), JavaScript)
	reader.ReadFolder(fmt.Sprintf("%s/%s", PlguinFolder, "python"), Python)
}

func (reader *PluginReader) ReadFolder(folder string, pluginType PluginType) {
	fileInfos, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Println("PluginReader.ReadAllPlugin: " + err.Error())
		return
	}
	var regex *regexp.Regexp
	switch pluginType {
	case Binary:
		if reader.os == "windows" {
			regex, err = regexp.Compile("(.+?)\\-([^\\-]+)\\.exe")
		} else {
			regex, err = regexp.Compile("(.+?)\\-([^\\-]+)\\.bin")
		}
	case Java:
		regex, err = regexp.Compile("(.+?)\\-([^\\-]+)\\.jar")
	case JavaScript:
		regex, err = regexp.Compile("(.+?)\\-([^\\-]+)\\.js")
	case Python:
		regex, err = regexp.Compile("(.+?)\\-([^\\-]+)\\.py")
	default:
		return
	}

	if err != nil {
		log.Println("PluginReader.ReadAllPlugin: " + err.Error())
		return
	}

	for _, fileInfo := range fileInfos {
		matches := regex.FindStringSubmatch(fileInfo.Name())
		if matches == nil {
			continue
		}
		keyword := matches[2]
		if reader.keywordSet[keyword] {
			continue
		}
		reader.keywordSet[keyword] = true
		path := fmt.Sprintf("%s/%s", folder, fileInfo.Name())
		plugin := Plugin{
			Type:          pluginType,
			Path:          path,
			ConfigKeyword: keyword,
		}
		reader.PluginList = append(reader.PluginList, plugin)
	}
}
