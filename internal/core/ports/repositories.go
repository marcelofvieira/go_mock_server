package ports

import (
	"context"
	"mock_server_mux/internal/core/domain"
)

type MockConfigurationRepository interface {
	GetAll(context.Context) ([]domain.MockConfiguration, error)
	GetById(context.Context, int) (domain.MockConfiguration, error)
	Save(context.Context, domain.MockConfiguration) (domain.MockConfiguration, error)
	Update(context.Context, domain.MockConfiguration) (domain.MockConfiguration, error)
	DeleteById(context.Context, int) error
}
