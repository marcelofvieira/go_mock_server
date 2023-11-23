package ports

import (
	"context"
	"mock_server_mux/internal/core/domain"
	"net/http"
)

type MockConfigurationService interface {
	GetMockConfigById(context.Context, int) (domain.MockConfiguration, error)
	AddNewMockConfiguration(context.Context, domain.MockConfiguration) (domain.MockConfiguration, error)
	UpdateMockConfiguration(context.Context, domain.MockConfiguration) (domain.MockConfiguration, error)
	DeleteMockConfiguration(context.Context, int) error
}

type MockRequestPreProcessorService interface {
	ProcessMockParameters(context.Context, domain.MockConfiguration) (domain.MockConfiguration, error)
}

type DynamicHandlerService interface {
	ProcessDynamicHandler(context.Context, *http.Request) (domain.MockConfiguration, error)
}

type RequestFilterService interface {
	FilterMockHandlersByRequest(context.Context, *http.Request, []domain.MockConfiguration) (domain.MockConfiguration, error)
}

type MockProcessorService interface {
	GetVariablesValues(context.Context, *http.Request, domain.MockConfiguration) (domain.MockConfiguration, error)
}

type MockResponseProcessor interface {
	ProcessMockResponse(context.Context, domain.MockConfiguration) (domain.MockConfiguration, error)
}
