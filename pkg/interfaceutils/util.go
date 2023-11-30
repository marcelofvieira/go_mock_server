package interfaceutils

import (
	"encoding/json"
)

func GetBytes(value interface{}) ([]byte, error) {

	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetToString(value interface{}) (string, error) {
	data, err := GetBytes(value)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
