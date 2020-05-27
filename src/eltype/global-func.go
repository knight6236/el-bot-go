package eltype

import (
	"errors"
	"fmt"
	"strconv"
)

func toString(nativeValue interface{}) (string, error) {
	value := ""
	switch nativeValue.(type) {
	case string:
		value = nativeValue.(string)
	case int:
		value = strconv.Itoa(nativeValue.(int))
	case int64:
		value = strconv.FormatInt(nativeValue.(int64), 10)
	case float64:
		value = fmt.Sprintf("%.6f", nativeValue.(float64))
	case float32:
		value = fmt.Sprintf("%.6f", nativeValue.(float32))
	case bool:
		value = strconv.FormatBool(nativeValue.(bool))
	default:
		return value, errors.New("")
	}
	return value, nil
}

func toInt64(nativeValue interface{}) (int64, error) {
	var value int64
	var err error
	switch nativeValue.(type) {
	case string:
		value, err = strconv.ParseInt(nativeValue.(string), 10, 64)
		if err != nil {
			return value, err
		}
	case int:
		value = int64(nativeValue.(int))
	case int64:
		value = nativeValue.(int64)
	case float64:
		value = int64(nativeValue.(float64))
	case float32:
		value = int64(nativeValue.(float32))
	case bool:
		if nativeValue.(bool) {
			value = 1
		} else {
			value = 0
		}
	default:
		return value, errors.New("")
	}
	return value, nil
}
