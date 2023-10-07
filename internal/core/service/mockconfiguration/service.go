package mockconfiguration

import (
	"errors"
	"mock_server_mux/internal/core/domain"
	"mock_server_mux/internal/core/ports"
)

type Service struct {
	mockRepository ports.MockConfigurationRepository
}

func NewService(mockRepository ports.MockConfigurationRepository) *Service {
	return &Service{
		mockRepository: mockRepository,
	}
}

func (s *Service) GetMockConfigById(Id int) (domain.MockConfiguration, error) {
	mockConfig, err := s.mockRepository.GetById(Id)

	if err != nil {
		return domain.MockConfiguration{}, err
	}

	return mockConfig, nil
}

func (s *Service) AddNewMockConfiguration(mockConfig domain.MockConfiguration) (domain.MockConfiguration, error) {
	mockConfig, err := s.mockRepository.Save(mockConfig)

	if err != nil {
		return domain.MockConfiguration{}, errors.New("create mock configuration into repository has failed")
	}

	return mockConfig, nil
}

func (s *Service) UpdateMockConfiguration(mockConfig domain.MockConfiguration) (domain.MockConfiguration, error) {
	mockConfig, err := s.mockRepository.Save(mockConfig)

	if err != nil {
		return domain.MockConfiguration{}, errors.New("update mock configuration into repository has failed")
	}

	return mockConfig, nil
}
