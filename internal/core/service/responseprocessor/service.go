package responseprocessor

import (
	"context"
	"fmt"
	"mock_server_mux/internal/core/domain"
	"mock_server_mux/pkg/utils/regex"
	"mock_server_mux/pkg/utils/request"
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

		found, variables := regex.FindStringValuesRegex(regex.FindResponseVariablePattern, value)
		if !found {
			continue
		}

		for _, variable := range variables {

			found, variableContext := regex.FindStringValuesRegex(regex.FindVariableContextPattern, variable[0])
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

	body, err := request.MockBodyToString(mockConfig.Response.Body)
	if err != nil {
		return err
	}

	found, variables := regex.FindStringValuesRegex(regex.FindResponseVariablePattern, body)
	if !found {
		return nil
	}

	for _, variable := range variables {

		found, variableContext := regex.FindStringValuesRegex(regex.FindVariableContextPattern, variable[0])
		if !found {
			continue
		}

		variableName, replaceContent := analyseVariable(variable[0])

		if variableValue, ok := mockConfig.MockVariables[variableContext[0][1]][variableName]; ok {
			body = strings.Replace(body, replaceContent, fmt.Sprintf("%v", variableValue), 1)
		} else {
			body = strings.Replace(body, replaceContent, "", 1)
		}

	}

	mockConfig.Response.Body = body

	return nil
}

func analyseVariable(variable string) (string, string) {

	var variableName, replaceContent string

	variableName = variable
	replaceContent = variable

	found, _ := regex.FindStringValuesRegex(regex.FindNumberBooleanVariablePattern, variable)
	if !found {
		return variableName, replaceContent
	}

	found, variableInfo := regex.FindStringValuesRegex(regex.FindResponseVariablePattern, variable)
	if !found {
		return variableName, replaceContent
	}

	variableName = "${" + variableInfo[0][1] + "}"
	replaceContent = "\"" + variable + "\""

	return variableName, replaceContent
}
