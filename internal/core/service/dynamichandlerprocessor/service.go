package dynamichandlerprocessor

import (
	"context"
	"mock_server_mux/internal/core/domain"
	"mock_server_mux/internal/core/ports"
	"mock_server_mux/pkg/apperrors"
	"net/http"
)

type Service struct {
	mockRepository       ports.MockConfigurationRepository
	requestFilterService ports.RequestFilterService
}

func NewService(mockRepository ports.MockConfigurationRepository, requestFilterService ports.RequestFilterService) *Service {
	return &Service{
		mockRepository:       mockRepository,
		requestFilterService: requestFilterService,
	}
}

func (s *Service) ProcessDynamicHandler(ctx context.Context, request *http.Request) (domain.MockConfiguration, error) {
	allMockConfig, err := s.mockRepository.GetAll(ctx)
	if err != nil {
		if apperrors.Is(err, apperrors.NotFound) {
			return domain.MockConfiguration{}, apperrors.New(apperrors.NotImplemented, err, "handler not implemented")
		} else {
			return domain.MockConfiguration{}, apperrors.New(apperrors.InternalServerError, err, "internal server error")
		}
	}

	result, err := s.requestFilterService.FilterMockHandlersByRequest(ctx, request, allMockConfig)
	if err != nil {
		if apperrors.Is(err, apperrors.NotFound) {
			return domain.MockConfiguration{}, apperrors.New(apperrors.NotImplemented, err, "handler not implemented")
		} else {
			return domain.MockConfiguration{}, apperrors.New(apperrors.InternalServerError, err, "internal server error")
		}
	}

	return result, nil
}
