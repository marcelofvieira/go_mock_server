package responseprocessor

import (
	"context"
	"encoding/json"
	"mock_server_mux/internal/core/domain"
	"mock_server_mux/pkg/regexutil"
	"strings"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ProcessMockResponse(ctx context.Context, mockConfig domain.MockConfiguration) (domain.MockConfiguration, error) {

	err := s.processHeaders(&mockConfig)
	if err != nil {
		return domain.MockConfiguration{}, err
	}

	err = s.processPayload(&mockConfig)
	if err != nil {
		return domain.MockConfiguration{}, err
	}

	return mockConfig, nil
}

func (s *Service) processHeaders(mockConfig *domain.MockConfiguration) error {

	if len(mockConfig.Response.Headers) == 0 {
		return nil
	}

	for headerName, headerValue := range mockConfig.Response.Headers {

		value, _ := headerValue.(string)

		found, variables := regexutil.FindStringValuesRegex(regexutil.FindResponseVariablePattern, value)
		if !found {
			continue
		}

		for _, variable := range variables {

			found, variableContext := regexutil.FindStringValuesRegex(regexutil.FindVariableContextPattern, variable[0])
			if !found {
				continue
			}

			if variableValue, ok := mockConfig.MockVariables[variableContext[0][1]][variable[0]].(string); ok {
				value = strings.Replace(value, variable[0], variableValue, 1)
			} else {
				value = strings.Replace(value, variable[0], "", 1)
			}

		}

		mockConfig.Response.Headers[headerName] = value
	}

	return nil
}

func (s *Service) processPayload(mockConfig *domain.MockConfiguration) error {
	if mockConfig.Response.Body == nil {
		return nil
	}

	json, err := json.Marshal(mockConfig.Response.Body)
	if err != nil {
		return err
	}

	body := string(json)

	found, variables := regexutil.FindStringValuesRegex(regexutil.FindResponseVariablePattern, body)
	if !found {
		return nil
	}

	for _, variable := range variables {

		found, variableContext := regexutil.FindStringValuesRegex(regexutil.FindVariableContextPattern, variable[0])
		if !found {
			continue
		}

		if variableValue, ok := mockConfig.MockVariables[variableContext[0][1]][variable[0]].(string); ok {
			body = strings.Replace(body, variable[0], variableValue, 1)
		} else {
			body = strings.Replace(body, variable[0], "", 1)
		}

	}

	mockConfig.Response.Body = body

	return nil
}
