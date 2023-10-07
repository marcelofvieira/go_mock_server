package ports

import (
	"mock_server_mux/internal/core/domain"
	"net/http"
)

type MockConfigurationService interface {
	GetMockConfigById(int) (domain.MockConfiguration, error)
	AddNewMockConfiguration(domain.MockConfiguration) (domain.MockConfiguration, error)
	UpdateMockConfiguration(domain.MockConfiguration) (domain.MockConfiguration, error)
}

type DynamicHandlerService interface {
	ProcessDynamicHandler(*http.Request) (interface{}, error)
}
