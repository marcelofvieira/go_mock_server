package handlerprocessor

import (
	"context"
	"mock_server_mux/internal/core/ports"
	"mock_server_mux/pkg/apperrors"
	"net/http"
)

type Service struct {
	mockRepository ports.MockConfigurationRepository
}

func NewService(mockRepository ports.MockConfigurationRepository) *Service {
	return &Service{
		mockRepository: mockRepository,
	}
}

func (s *Service) ProcessDynamicHandler(ctx context.Context, request *http.Request) (interface{}, error) {

	mockConfig, err := s.mockRepository.GetAll(ctx)
	if err != nil {
		if apperrors.Is(err, apperrors.NotFound) {
			return nil, apperrors.New(apperrors.NotImplemented, err, "handler not implemented")
		} else {
			return nil, apperrors.New(apperrors.InternalServerError, err, "internal server error")
		}
	}

	return mockConfig, nil
}
