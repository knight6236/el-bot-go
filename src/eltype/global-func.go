package eltype

import (
	"el-bot-go/src/gomirai"
	"fmt"
	"strconv"
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

func mergeConfigList(dest *[]Config, lists ...[]Config) {
	for _, list := range lists {
		for _, item := range list {
			*dest = append(*dest, item)
		}
	}
}

func mergeGoMiraiMessageList(dest *[]gomirai.Message, lists ...[]gomirai.Message) {
	for _, list := range lists {
		for _, item := range list {
			*dest = append(*dest, item)
		}
	}
}

func parseJsonObj(obj interface{}, callDepth int) ([]string, []string) {
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
			tmp0, tmp1 := parseJsonObj(obj.([]interface{})[i], callDepth+1)
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
			tmp0, tmp1 := parseJsonObj(value, callDepth+1)
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
