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
	processMockResponse  ports.MockResponseProcessor
}

func NewService(mockRepository ports.MockConfigurationRepository, requestFilterService ports.RequestFilterService, processMockResponse ports.MockResponseProcessor) *Service {
	return &Service{
		mockRepository:       mockRepository,
		requestFilterService: requestFilterService,
		processMockResponse:  processMockResponse,
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

	mockConfig, err := s.requestFilterService.FilterMockHandlersByRequest(ctx, request, allMockConfig)
	if err != nil {
		if apperrors.Is(err, apperrors.NotFound) {
			return domain.MockConfiguration{}, apperrors.New(apperrors.NotImplemented, err, "handler not implemented")
		} else {
			return domain.MockConfiguration{}, apperrors.New(apperrors.InternalServerError, err, "internal server error")
		}
	}

	mockConfig, err = s.processMockResponse.ProcessMockResponse(ctx, mockConfig)
	if err != nil {
		return domain.MockConfiguration{}, apperrors.New(apperrors.InternalServerError, err, "internal server error")
	}

	return mockConfig, nil
}
