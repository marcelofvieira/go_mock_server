package dynamichandler

import (
	"encoding/json"
	"mock_server_mux/internal/core/domain"
	"mock_server_mux/internal/core/ports"
	"mock_server_mux/pkg/apperrors"
	"mock_server_mux/pkg/logger"
	_interface "mock_server_mux/pkg/utils/interface"
	"mock_server_mux/pkg/utils/response"
	"net/http"
	"strconv"
	"time"
)

type HTTPHandler struct {
	dynamicHandlerService ports.DynamicHandlerService
}

func NewHTTPHandler(mockService ports.DynamicHandlerService) *HTTPHandler {
	return &HTTPHandler{
		dynamicHandlerService: mockService,
	}
}

func (hdl *HTTPHandler) ProcessDynamicHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	httpStatus := http.StatusOK

	handlerResponse, err := hdl.dynamicHandlerService.ProcessDynamicHandler(ctx, r)

	if err != nil {
		if apperrors.Is(err, apperrors.NotImplemented) {
			httpStatus = http.StatusNotImplemented
		} else {
			httpStatus = http.StatusInternalServerError
		}

		err := response.Error(w, r, httpStatus, err)
		if err != nil {
			logger.Error("Error: ", err)
		}

		return
	}

	httpStatus, headers, body := processMockConfigurationResponse(handlerResponse)

	err = response.MockResponse(w, r, httpStatus, headers, body)
	if err != nil {
		response.Error(w, r, httpStatus, err)
		if err != nil {
			logger.Error("Error: ", err)
		}
	}

}

func processMockConfigurationResponse(mockConfig domain.MockConfiguration) (int, map[string]string, []byte) {
	statusCode := processMockHttpStatus(mockConfig)

	headers := processMockHeaders(mockConfig)

	payload := processMockPayload(mockConfig)

	processMockDelay(mockConfig)

	return statusCode, headers, payload
}

func processMockHttpStatus(mockConfig domain.MockConfiguration) int {
	if http.StatusText(mockConfig.Response.StatusCode) == "" {
		return http.StatusOK
	}

	return mockConfig.Response.StatusCode
}

func processMockHeaders(mockConfig domain.MockConfiguration) map[string]string {
	headers := make(map[string]string)

	for key, value := range mockConfig.Response.Headers {
		headers[key], _ = value.(string)
	}

	config, _ := _interface.GetToString(mockConfig.Response.Configuration[domain.ShowMockInformation])

	showConfig, err := strconv.ParseBool(config)
	if err != nil {
		showConfig = false
	}

	if showConfig {
		headers["Mock-Id"] = strconv.Itoa(mockConfig.Id)
		headers["Mock-name"] = mockConfig.Info.TestName
		headers["Mock-Group"] = mockConfig.Info.TestGroup
		headers["Mock-Regex-Process"] = strconv.Itoa(mockConfig.Request.Regex.Count)
	}

	return headers
}

func processMockPayload(mockConfig domain.MockConfiguration) []byte {
	var payload []byte
	var err error

	switch mockConfig.Response.Body.(type) {
	case string: //Working json with string
		var jsonData interface{}

		err = json.Unmarshal([]byte(mockConfig.Response.Body.(string)), &jsonData)
		if err != nil {
			return []byte(mockConfig.Response.Body.(string))
		}

		payload, err = json.MarshalIndent(jsonData, "", "\t")
		if err != nil {
			return []byte(mockConfig.Response.Body.(string))
		}

	default: //Working json with object
		payload, err = json.Marshal(mockConfig.Response.Body)

		if err != nil {
			return []byte(mockConfig.Response.Body.(string))
		}

	}

	return payload
}

func processMockDelay(mockConfig domain.MockConfiguration) {
	delayConfig, _ := _interface.GetToString(mockConfig.Response.Configuration[domain.ConfigResponseDelay])

	delay, err := strconv.Atoi(delayConfig)
	if err != nil {
		delay = 2
	}

	if delay <= 0 || delay >= 60000 {
		delay = 50
	}

	time.Sleep(time.Duration(delay) * time.Millisecond)
}
