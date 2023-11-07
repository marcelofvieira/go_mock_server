package responseprocessor

import (
	"context"
	"mock_server_mux/internal/core/domain"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ProcessMockResponse(ctx context.Context, mockConfig domain.MockConfiguration) (domain.MockConfiguration, error) {
	return mockConfig, nil
}
