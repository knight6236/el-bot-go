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
	PluginMap  map[string]Plugin
	randKeySet map[string]bool
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
	compiledFolder := "compile"
	msgProcFolder := "msgproc"
	switch reader.os {
	case "freebsd", "linux":
		reader.ReadFolder(
			fmt.Sprintf("%s/%s/%s/%s", PlguinFolder, compiledFolder, reader.os, reader.arch),
			PluginTypeBinary, false)
		reader.ReadFolder(
			fmt.Sprintf("%s/%s/%s/%s", PlguinFolder, msgProcFolder, reader.os, reader.arch),
			PluginTypeBinary, true)
	case "darwin", "windows":
		reader.ReadFolder(
			fmt.Sprintf("%s/%s/%s/%s", PlguinFolder, compiledFolder, reader.os, reader.arch),
			PluginTypeBinary, false)
		reader.ReadFolder(
			fmt.Sprintf("%s/%s/%s/%s", PlguinFolder, msgProcFolder, reader.os, reader.arch),
			PluginTypeBinary, true)

		if reader.arch == "amd64" {
			reader.ReadFolder(fmt.Sprintf("%s/%s/%s/%s", PlguinFolder, compiledFolder, reader.os, "amd64"),
				PluginTypeBinary, false)
			reader.ReadFolder(fmt.Sprintf("%s/%s/%s/%s", PlguinFolder, msgProcFolder, reader.os, "amd64"),
				PluginTypeBinary, true)
		}
	}
	reader.ReadFolder(
		fmt.Sprintf("%s/%s/%s", PlguinFolder, compiledFolder, "java"),
		PluginTypeJava, false)
	reader.ReadFolder(
		fmt.Sprintf("%s/%s/%s", PlguinFolder, msgProcFolder, "java"),
		PluginTypeJava, true)

	reader.ReadFolder(
		fmt.Sprintf("%s/%s/%s", PlguinFolder, compiledFolder, "javaScript"),
		PluginTypeJavaScript, false)
	reader.ReadFolder(
		fmt.Sprintf("%s/%s/%s", PlguinFolder, msgProcFolder, "javaScript"),
		PluginTypeJavaScript, true)

	reader.ReadFolder(
		fmt.Sprintf("%s/%s/%s", PlguinFolder, compiledFolder, "python"),
		PluginTypePython, false)
	reader.ReadFolder(
		fmt.Sprintf("%s/%s/%s", PlguinFolder, msgProcFolder, "python"),
		PluginTypePython, true)
}

func (reader *PluginReader) ReadFolder(folder string, pluginType PluginType, isMsgProc bool) {
	fileInfos, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Println("PluginReader.ReadAllPlugin: " + err.Error())
		return
	}
	var regex *regexp.Regexp
	switch pluginType {
	case PluginTypeBinary:
		if reader.os == "windows" {
			regex, err = regexp.Compile("(.+?)\\-([^\\-]+)\\.exe")
		} else {
			regex, err = regexp.Compile("(.+?)\\-([^\\-]+)\\.bin")
		}
	case PluginTypeJava:
		regex, err = regexp.Compile("(.+?)\\-([^\\-]+)\\.jar")
	case PluginTypeJavaScript:
		regex, err = regexp.Compile("(.+?)\\-([^\\-]+)\\.js")
	case PluginTypePython:
		regex, err = regexp.Compile("(.+?)\\-([^\\-]+)\\.py")
	default:
		return
	}

	if err != nil {
		log.Println("PluginReader.ReadAllPlugin: " + err.Error())
		return
	}

	for _, fileInfo := range fileInfos {
		keyword := ""
		if !isMsgProc {
			matches := regex.FindStringSubmatch(fileInfo.Name())
			if matches == nil {
				continue
			}
			keyword := matches[2]
			if reader.keywordSet[keyword] {
				continue
			}
			reader.keywordSet[keyword] = true
		}

		randKey := RandString(25)
		for !reader.randKeySet[randKey] {
			randKey = RandString(25)
		}
		reader.randKeySet[randKey] = true

		path := fmt.Sprintf("%s/%s", folder, fileInfo.Name())
		plugin := Plugin{
			Type:          pluginType,
			Path:          path,
			ConfigKeyword: keyword,
			IsProcMsg:     isMsgProc,
			RandKey:       randKey,
		}
		reader.PluginMap[randKey] = plugin
	}
}
