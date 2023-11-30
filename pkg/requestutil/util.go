package requestutil

import (
	"bytes"
	"encoding/json"
	"io"
	"mock_server_mux/pkg/interfaceutils"
	"mock_server_mux/pkg/logger"
	"net/http"
)

func MockBodyToString(mockBody interface{}) (string, error) {
	json, err := json.Marshal(mockBody)
	if err != nil {
		return "", err
	}

	body := string(json)

	return body, nil
}

func ReadBodyToString(request *http.Request) (string, error) {
	switch getContentType(request) {
	case "application/json":
		return getJsonBodyToString(request)

	default:
		return getBodyToString(request)
	}
}

func getContentType(request *http.Request) string {
	contentType := request.Header.Get("Content-Type")

	return contentType
}

func getJsonBodyToString(request *http.Request) (string, error) {
	body, err := io.ReadAll(request.Body)

	if err != nil {
		logger.Error("Error reading body", err)
		return "", err
	}

	//set body again
	request.Body = io.NopCloser(bytes.NewBuffer(body))

	var iBody interface{}

	json.Unmarshal(body, &iBody)
	if err != nil {
		logger.Error("Error transform to json", err)
		return "", err
	}

	requestBody, err := interfaceutils.GetToString(iBody)
	if err != nil {
		logger.Error("Error transform json to string", err)
		return "", err
	}

	return requestBody, nil
}

func getBodyToString(request *http.Request) (string, error) {
	body, err := io.ReadAll(request.Body)

	if err != nil {
		logger.Error("Error reading body", err)
		return "", err
	}

	//set body again
	request.Body = io.NopCloser(bytes.NewBuffer(body))

	requestBody := string(body)

	return requestBody, nil
}
