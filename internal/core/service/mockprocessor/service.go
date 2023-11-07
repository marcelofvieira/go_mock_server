package mockprocessor

import (
	"context"
	"mock_server_mux/internal/core/domain"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ProcessMock(ctx context.Context, mockConfig domain.MockConfiguration) (domain.MockConfiguration, error) {

	//TODO: Get variables values

	return mockConfig, nil

}
