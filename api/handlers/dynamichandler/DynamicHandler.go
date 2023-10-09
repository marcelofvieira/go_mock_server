package dynamichandler

import (
	"mock_server_mux/internal/core/ports"
	"mock_server_mux/pkg/apperrors"
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
			return
		}
		return
	}

	err = response.Success(w, r, httpStatus, handlerResponse)
	if err != nil {
		return
	}
}
