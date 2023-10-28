package dynamichandler

import (
	"encoding/json"
	"mock_server_mux/internal/core/domain"
	"mock_server_mux/internal/core/ports"
	"mock_server_mux/pkg/apperrors"
	"mock_server_mux/pkg/logger"
	"mock_server_mux/pkg/response"
	"net/http"
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

	httpStatus, headers, body, delay := processMockResponse(handlerResponse)

	err = response.MockResponse(w, r, httpStatus, headers, body, delay)
	if err != nil {
		response.Error(w, r, httpStatus, err)
		if err != nil {
			logger.Error("Error: ", err)
		}
	}

}

func processMockResponse(mockConfig domain.MockConfiguration) (int, map[string]string, []byte, int) {
	var body []byte

	statusCode := 0
	headers := make(map[string]string)
	delay := 0

	statusCode = mockConfig.Response.StatusCode

	if http.StatusText(mockConfig.Response.StatusCode) == "" {
		statusCode = http.StatusOK
	}

	payload, err := processMockPayload(mockConfig.Response.Body)
	if err == nil {
		body = payload
	} else {
		body = []byte(mockConfig.Response.Body.(string))
	}

	for _, header := range mockConfig.Response.Headers {
		headers[header.Key] = header.Value
	}

	delay = mockConfig.Response.Configuration.ResponseDelay

	if delay <= 0 || delay >= 20000 {
		delay = 50
	}

	return statusCode, headers, body, delay
}

func processMockPayload(body interface{}) ([]byte, error) {
	var payload []byte
	var err error

	switch body.(type) {
	case string:
		var jsonData interface{}

		err = json.Unmarshal([]byte(body.(string)), &jsonData)
		if err != nil {
			return nil, err
		}

		payload, err = json.MarshalIndent(jsonData, "", "\t")
		if err != nil {
			return nil, err
		}

	default:
		payload, err = json.Marshal(body)

		if err != nil {
			return nil, err
		}

	}

	return payload, nil
}
