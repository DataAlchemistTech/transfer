package maputil

import (
	"fmt"
	"strconv"
)

func GetKeyFromMap(obj map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if len(obj) == 0 {
		return defaultValue
	}

	val, isOk := obj[key]
	if !isOk {
		return defaultValue
	}

	return val
}

func GetIntegerFromMap(obj map[string]interface{}, key string) (int, error) {
	if len(obj) == 0 {
		return 0, fmt.Errorf("object is empty")
	}

	valInterface, isOk := obj[key]
	if !isOk {
		return 0, fmt.Errorf("key: %s does not exist in object", key)
	}

	val, err := strconv.Atoi(fmt.Sprint(valInterface))
	if err != nil {
		return 0, fmt.Errorf("key: %s is not type integer, err: %v", key, err)
	}

	return val, nil
}
