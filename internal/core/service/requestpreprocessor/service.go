package requestpreprocessor

import (
	"context"
	"mock_server_mux/internal/core/domain"
	"strings"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ProcessMockParameters(ctx context.Context, mockConfig domain.MockConfiguration) (domain.MockConfiguration, error) {

	// interfaceutils.GetToString(mockConfig.Request.Body)

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
			variables = append(variables, domain.Variable{
				Name:    URL[initPos:index],
				Context: "path",
			})
			initPos = -1
		}
	}

	mockConfig.Request.Variables = variables

	if len(variables) > 0 {
		for _, variable := range variables {
			find := "{" + variable.Name + "}"
			URL = strings.Replace(URL, find, "([^/]+)", 1)
		}

		mockConfig.Request.RegexURL = URL + "+$"

	}

	return mockConfig, nil
}

func processQueryParamConfig(mockConfig domain.MockConfiguration) (domain.MockConfiguration, error) {

	variables := make([]domain.Variable, 0)

	URL := mockConfig.Request.URL

	initPos := -1

	for index, char := range URL {
		if char == '{' {
			initPos = index + 1
		}

		if char == '}' {
			variables = append(variables, domain.Variable{
				Name:    URL[initPos:index],
				Context: "path",
			})

			initPos = -1
		}
	}

	mockConfig.Request.Variables = variables

	return mockConfig, nil
}
