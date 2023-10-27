package mockconfiguration

import (
	"context"
	"errors"
	"mock_server_mux/internal/core/domain"
	"mock_server_mux/internal/core/ports"
	"mock_server_mux/pkg/apperrors"
)

type Service struct {
	mockRepository          ports.MockConfigurationRepository
	mockRequestPreProcessor ports.MockRequestPreProcessorService
}

func NewService(mockRepository ports.MockConfigurationRepository, mockRequestPreProcessor ports.MockRequestPreProcessorService) *Service {
	return &Service{
		mockRepository:          mockRepository,
		mockRequestPreProcessor: mockRequestPreProcessor,
	}
}

func (s *Service) GetMockConfigById(ctx context.Context, Id int) (domain.MockConfiguration, error) {

	mockConfig, err := s.mockRepository.GetById(ctx, Id)
	if err != nil {
		return domain.MockConfiguration{}, err
	}

	return mockConfig, nil
}

func (s *Service) DeleteMockConfiguration(ctx context.Context, Id int) error {
	err := s.mockRepository.DeleteById(ctx, Id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) AddNewMockConfiguration(ctx context.Context, mockConfig domain.MockConfiguration) (domain.MockConfiguration, error) {

	mockConfig, err := s.mockRequestPreProcessor.ProcessRequestParams(ctx, mockConfig)
	if err != nil {
		return domain.MockConfiguration{}, apperrors.New(apperrors.InternalServerError, err, "Error to process request params")
	}

	mockConfig, err = s.mockRepository.Save(ctx, mockConfig)
	if err != nil {
		return domain.MockConfiguration{}, errors.New("create mock configuration into repository has failed")
	}

	return mockConfig, nil
}

func (s *Service) UpdateMockConfiguration(ctx context.Context, mockConfig domain.MockConfiguration) (domain.MockConfiguration, error) {
	mockConfig, err := s.mockRepository.Save(ctx, mockConfig)

	if err != nil {
		return domain.MockConfiguration{}, errors.New("update mock configuration into repository has failed")
	}

	return mockConfig, nil
}
