package mockconfighandler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"mock_server_mux/internal/core/domain"
	"mock_server_mux/internal/core/ports"
	"mock_server_mux/pkg/apperrors"
	"mock_server_mux/pkg/response"
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

	//ctx := r.Context()

	vars := mux.Vars(r)

	Id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, apperrors.New(apperrors.BadRequest, err, "invalid id"))
		return
	}

	configMock, err := hdl.mockConfigService.GetMockConfigById(Id)

	if err != nil {
		if apperrors.Is(err, apperrors.NotFound) {
			response.Error(w, r, http.StatusNotFound, err)
		} else {
			response.Error(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	response.Success(w, r, http.StatusOK, configMock.ToPresenter())

}

func (hdl *HTTPHandler) AddMockConfiguration(w http.ResponseWriter, r *http.Request) {

	body := domain.MockConfiguration{}
	decode := json.NewDecoder(r.Body)
	decode.Decode(&body)

	configMock, err := hdl.mockConfigService.AddNewMockConfiguration(body)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	response.Success(w, r, http.StatusCreated, configMock.ToPresenter())
}

func (hdl *HTTPHandler) UpdateMockConfiguration(w http.ResponseWriter, r *http.Request) {

	body := domain.MockConfiguration{}

	decode := json.NewDecoder(r.Body)

	err := decode.Decode(&body)
	if err != nil {
		err := response.Error(w, r, http.StatusInternalServerError, err)
		if err != nil {
			return
		}
	}

	configMock, err := hdl.mockConfigService.UpdateMockConfiguration(body)
	if err != nil {
		err := response.Error(w, r, http.StatusInternalServerError, err)
		if err != nil {
			return
		}
		return
	}

	err = response.Success(w, r, http.StatusCreated, configMock.ToPresenter())
	if err != nil {
		return
	}
}
