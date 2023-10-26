package dynamichandler

import (
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
		return
	}

}

func processMockResponse(mockConfig domain.MockConfiguration) (int, map[string]string, interface{}, int) {

	statusCode := 0
	headers := make(map[string]string)
	delay := 0

	statusCode = mockConfig.Response.StatusCode

	if http.StatusText(mockConfig.Response.StatusCode) == "" {
		statusCode = http.StatusOK
	}

	body := mockConfig.Response.Body

	for _, header := range mockConfig.Response.Headers {
		headers[header.Key] = header.Value
	}

	delay = mockConfig.Response.Configuration.ResponseDelay

	if delay <= 0 || delay >= 20000 {
		delay = 50
	}

	return statusCode, headers, body, delay
}
