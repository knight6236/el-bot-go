package eltype

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"

	"github.com/ADD-SP/gomirai"
)

func CastInt64ToString(nativeValue int64) string {
	return strconv.FormatInt(nativeValue, 10)
}

func CastStringToInt64(str string) int64 {
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {

	}
	return value
}

func MergeConfigList(dest *[]Config, lists ...[]Config) {
	for _, list := range lists {
		for _, item := range list {
			*dest = append(*dest, item)
		}
	}
}

func MergeGoMiraiMessageList(dest *[]gomirai.Message, lists ...[]gomirai.Message) {
	for _, list := range lists {
		for _, item := range list {
			*dest = append(*dest, item)
		}
	}
}

func ParseJsonObj(obj interface{}, callDepth int) ([]string, []string) {
	var varNameList []string
	var valueList []string
	if callDepth >= 20 {
		return nil, nil
	}
	switch obj.(type) {
	case string:
		varNameList = append(varNameList, "")
		valueList = append(valueList, obj.(string))
	case int:
		varNameList = append(varNameList, "")
		valueList = append(valueList, strconv.Itoa(obj.(int)))
	case int8:
		varNameList = append(varNameList, "")
		valueList = append(valueList, fmt.Sprintf("%d", obj.(int8)))
	case int32:
		varNameList = append(varNameList, "")
		valueList = append(valueList, fmt.Sprintf("%d", obj.(int32)))
	case int64:
		varNameList = append(varNameList, "")
		valueList = append(valueList, fmt.Sprintf("%d", obj.(int64)))
	case bool:
		varNameList = append(varNameList, "")
		valueList = append(valueList, strconv.FormatBool(obj.(bool)))
	case float32:
		varNameList = append(varNameList, "")
		valueList = append(valueList, fmt.Sprintf("%.2f", obj.(float32)))
	case float64:
		varNameList = append(varNameList, "")
		valueList = append(valueList, fmt.Sprintf("%.2f", obj.(float64)))
	case []interface{}:
		for i := 0; i < len(obj.([]interface{})); i++ {
			name := fmt.Sprintf("[%d]", i)
			tmp0, tmp1 := ParseJsonObj(obj.([]interface{})[i], callDepth+1)
			for j := 0; j < len(tmp0); j++ {
				if tmp0[j] == "" || tmp0[j][0] == '[' {
					varNameList = append(varNameList, fmt.Sprintf("%s%s", name, tmp0[j]))
				} else {
					varNameList = append(varNameList, fmt.Sprintf("%s.%s", name, tmp0[j]))
				}
				valueList = append(valueList, tmp1[j])
			}
		}
	case map[string]interface{}:
		for key, value := range obj.(map[string]interface{}) {
			name := key
			tmp0, tmp1 := ParseJsonObj(value, callDepth+1)
			for j := 0; j < len(tmp0); j++ {
				if tmp0[j] == "" || tmp0[j][0] == '[' {
					varNameList = append(varNameList, fmt.Sprintf("%s%s", name, tmp0[j]))
				} else {
					varNameList = append(varNameList, fmt.Sprintf("%s.%s", name, tmp0[j]))
				}
				valueList = append(valueList, tmp1[j])
			}
		}
	}
	return varNameList, valueList
}

func ExecCommand(command string, args ...string) (string, error) {
	var cmd *exec.Cmd
	switch len(args) {
	case 1:
		cmd = exec.Command(command, args[0])
	case 2:
		cmd = exec.Command(command, args[0], args[1])
	case 3:
		cmd = exec.Command(command, args[0], args[1], args[2])
	case 4:
		cmd = exec.Command(command, args[0], args[1], args[2], args[3])

	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		return "", err
	}
	if opBytes, err := ioutil.ReadAll(stdout); err != nil {
		return "", err
	} else {
		return string(opBytes), nil
	}
}
