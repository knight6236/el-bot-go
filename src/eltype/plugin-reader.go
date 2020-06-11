package eltype

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	reader.PluginMap = make(map[string]Plugin)
	reader.randKeySet = make(map[string]bool)
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
	dirInfos, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Println("PluginReader.ReadFolder: " + err.Error())
		return
	}
	var regex *regexp.Regexp
	regex, err = regexp.Compile("(.+?)\\-([^\\-]+)")

	if err != nil {
		log.Println("PluginReader.ReadFolder: " + err.Error())
		return
	}

	for _, dirInfo := range dirInfos {
		if !dirInfo.IsDir() {
			continue
		}
		if err != nil {
			log.Println("PluginReader.ReadFolder: " + err.Error())
			continue
		}
		keyword := ""
		randKey := ""
		path := ""
		entry := ""
		switch pluginType {
		case PluginTypeBinary:
			entry = "start.exe"
		case PluginTypeJava:
			entry = "start.jar"
		case PluginTypeJavaScript:
			entry = "start.js"
		case PluginTypePython:
			entry = "start.py"
		}
		path = fmt.Sprintf("%s/%s/%s", folder, dirInfo.Name(), entry)
		if _, err := os.Lstat(path); err != nil {
			log.Printf("no such file or directory: %s", path)
			continue
		}
		if isMsgProc {
			randKey = RandStringBytesMaskImpr(5)
			for reader.randKeySet[randKey] {
				randKey = RandStringBytesMaskImpr(5)
			}
			reader.randKeySet[randKey] = true

		} else {
			matches := regex.FindStringSubmatch(dirInfo.Name())
			if matches == nil {
				continue
			}
			keyword := matches[2]
			if reader.keywordSet[keyword] {
				continue
			}
			reader.keywordSet[keyword] = true
		}
		plugin := Plugin{
			Type:          pluginType,
			Path:          path,
			Name:          dirInfo.Name(),
			ConfigKeyword: keyword,
			IsProcMsg:     isMsgProc,
			RandKey:       randKey,
		}

		if isMsgProc {
			reader.PluginMap[randKey] = plugin
		}
	}
}
