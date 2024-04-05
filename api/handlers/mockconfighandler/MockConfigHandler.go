package mockconfighandler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"mock_server_mux/api/presenter"
	"mock_server_mux/internal/core/domain"
	"mock_server_mux/internal/core/ports"
	"mock_server_mux/pkg/apperrors"
	"mock_server_mux/pkg/utils/response"
	"net/http"
	"strconv"
)

type HTTPHandler struct {
	mockConfigService ports.MockConfigurationService
}

func NewHTTPHandler(mockService ports.MockConfigurationService) *HTTPHandler {
	return &HTTPHandler{
		mockConfigService: mockService,
	}
}

func (hdl *HTTPHandler) GetMockConfiguration(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	vars := mux.Vars(r)

	Id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, apperrors.New(apperrors.BadRequest, err, "invalid id"))
		return
	}

	mockConfiguration, err := hdl.mockConfigService.GetMockConfigById(ctx, Id)

	if err != nil {
		if apperrors.Is(err, apperrors.NotFound) {
			response.Error(w, r, http.StatusNotFound, err)
		} else {
			response.Error(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	responseBody, err := convertToPresenter(mockConfiguration)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	response.Success(w, r, http.StatusOK, responseBody)

}

func (hdl *HTTPHandler) DeleteMockConfiguration(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	vars := mux.Vars(r)

	Id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, apperrors.New(apperrors.BadRequest, err, "invalid id"))
		return
	}

	err = hdl.mockConfigService.DeleteMockConfiguration(ctx, Id)

	if err != nil {
		if apperrors.Is(err, apperrors.NotFound) {
			response.Error(w, r, http.StatusNotFound, err)
		} else {
			response.Error(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	response.Success(w, r, http.StatusNoContent, nil)

}

func (hdl *HTTPHandler) AddMockConfiguration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body := domain.MockConfiguration{}
	decode := json.NewDecoder(r.Body)
	decode.Decode(&body)

	mockConfiguration, err := hdl.mockConfigService.AddNewMockConfiguration(ctx, body)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	responseBody, err := convertToPresenter(mockConfiguration)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	response.Success(w, r, http.StatusCreated, responseBody)
}

func (hdl *HTTPHandler) UpdateMockConfiguration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, apperrors.New(apperrors.BadRequest, err, "invalid id"))
		return
	}

	body := domain.MockConfiguration{}

	decode := json.NewDecoder(r.Body)

	err = decode.Decode(&body)
	if err != nil {
		err := response.Error(w, r, http.StatusInternalServerError, err)
		if err != nil {
			return
		}
	}

	mockConfiguration, err := hdl.mockConfigService.UpdateMockConfiguration(ctx, body, id)
	if err != nil {
		err := response.Error(w, r, http.StatusInternalServerError, err)
		if err != nil {
			return
		}
		return
	}

	responseBody, err := convertToPresenter(mockConfiguration)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	response.Success(w, r, http.StatusOK, responseBody)

}

func convertToPresenter(mockConfig domain.MockConfiguration) (presenter.MockConfiguration, error) {

	jsonResp, err := json.Marshal(mockConfig)
	if err != nil {
		return presenter.MockConfiguration{}, err
	}

	responseBody := presenter.MockConfiguration{}

	err = json.Unmarshal(jsonResp, &responseBody)
	if err != nil {
		return presenter.MockConfiguration{}, err
	}

	return responseBody, nil
}
