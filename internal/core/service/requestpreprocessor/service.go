package requestpreprocessor

import (
	"context"
	"mock_server_mux/internal/core/domain"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ProcessRequestParams(ctx context.Context, mockConfig domain.MockConfiguration) (domain.MockConfiguration, error) {

	mockConfig, err := processUrlConfiguration(mockConfig)
	if err != nil {
		return domain.MockConfiguration{}, err
	}

	return mockConfig, nil
}

func processUrlConfiguration(mockConfig domain.MockConfiguration) (domain.MockConfiguration, error) {

	variables := make([]domain.Variable, 0)

	URL := mockConfig.Request.URL

	initPos := -1

	for index, char := range URL {

		if char == '{' {
			initPos = index + 1
		}

		if char == '}' {
			variables = append(variables, domain.Variable{Name: URL[initPos:index]})

			initPos = -1
		}

	}

	mockConfig.Request.Variables = variables

	return mockConfig, nil
}
