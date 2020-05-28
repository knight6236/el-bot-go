package eltype

import (
	"el-bot-go/src/gomirai"
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
