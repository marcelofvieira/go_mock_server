package ports

import "mock_server_mux/internal/core/domain"

type MockConfigurationRepository interface {
	GetAll() ([]domain.MockConfiguration, error)
	GetById(int) (domain.MockConfiguration, error)
	Save(mockConfig domain.MockConfiguration) (domain.MockConfiguration, error)
	Update(mockConfig domain.MockConfiguration) (domain.MockConfiguration, error)
}
