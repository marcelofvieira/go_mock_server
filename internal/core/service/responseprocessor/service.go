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

	err := s.processResponseHeaders(&mockConfig)
	if err != nil {
		return domain.MockConfiguration{}, err
	}

	err = s.processResponsePayload(&mockConfig)
	if err != nil {
		return domain.MockConfiguration{}, err
	}

	return mockConfig, nil
}

func (s *Service) processResponseHeaders(mockConfig *domain.MockConfiguration) error {

	//TODO: Set headers variables values
	return nil
}

func (s *Service) processResponsePayload(mockConfig *domain.MockConfiguration) error {

	//TODO: Set payload variables values
	return nil
}
