package requestpreprocessor

import (
	"context"
	"mock_server_mux/internal/core/domain"
	"mock_server_mux/pkg/utils/regex"
	"mock_server_mux/pkg/utils/request"
	stringutil "mock_server_mux/pkg/utils/string"
	"strings"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

const (
	PathVariable   = "path"
	QueryVariable  = "query"
	HeaderVariable = "header"
	BodyVariable   = "body"
)

func (s *Service) ProcessMockParameters(ctx context.Context, mockConfig domain.MockConfiguration) (domain.MockConfiguration, error) {

	processUserVariables(&mockConfig)

	processUrlVariables(&mockConfig)

	processQueryVariables(&mockConfig)

	processHeaderVariables(&mockConfig)

	processBodyVariables(&mockConfig)

	return mockConfig, nil
}

func addVariableControl(mockConfig *domain.MockConfiguration, variable string, value interface{}, context string) {
	if mockConfig.MockVariables == nil {
		mockConfig.MockVariables = make(map[string]map[string]interface{})
	}

	if mockConfig.MockVariables[context] == nil {
		mockConfig.MockVariables[context] = make(map[string]interface{})
	}

	mockConfig.MockVariables[context]["${"+context+"."+variable+"}"] = value
}

func processUserVariables(mockConfig *domain.MockConfiguration) {

	for contextKey, userVariables := range mockConfig.UserVariables {
		for key, value := range userVariables {
			addVariableControl(mockConfig, key, value, contextKey)
		}
	}
}

func processUrlVariables(mockConfig *domain.MockConfiguration) {

	URL := mockConfig.Request.URL

	found, variables := regex.FindStringValuesRegex(regex.FindVariablePattern, URL)
	if !found {
		return
	}

	if len(variables) > 0 {
		for _, variable := range variables {
			URL = strings.Replace(URL, variable[0], regex.FindVariableValuePattern, 1)

			addVariableControl(mockConfig, variable[1], variable[0], PathVariable)
		}

		mockConfig.Request.Regex.URL = URL + regex.FindToFinalPattern

		mockConfig.Request.Regex.Count += len(variables)

	}
}

func processQueryVariables(mockConfig *domain.MockConfiguration) {

	for parameterName, parameterValue := range mockConfig.Request.QueryParameters {

		value, _ := parameterValue.(string)

		found, variables := regex.FindStringValuesRegex(regex.FindVariablePattern, value)
		if !found {
			continue
		}

		if len(variables) > 0 {

			if mockConfig.Request.Regex.QueryParameters == nil {
				mockConfig.Request.Regex.QueryParameters = make(map[string]string)
			}

			for _, variable := range variables {
				value = strings.Replace(value, variable[0], regex.FindVariableValuePattern, 1)

				addVariableControl(mockConfig, variable[1], variable[0], QueryVariable)
			}

			mockConfig.Request.Regex.QueryParameters[parameterName] = value + regex.FindToFinalPattern

			mockConfig.Request.Regex.Count += len(variables)
		}
	}
}

func processHeaderVariables(mockConfig *domain.MockConfiguration) {

	for headerName, headerValue := range mockConfig.Request.Headers {

		value, _ := headerValue.(string)

		found, variables := regex.FindStringValuesRegex(regex.FindVariablePattern, value)
		if !found {
			continue
		}

		if len(variables) > 0 {

			if mockConfig.Request.Regex.Headers == nil {
				mockConfig.Request.Regex.Headers = make(map[string]string)
			}

			for _, variable := range variables {
				value = strings.Replace(value, variable[0], regex.FindVariableValuePattern, 1)

				addVariableControl(mockConfig, variable[1], variable[0], HeaderVariable)
			}

			mockConfig.Request.Regex.Headers[headerName] = value + regex.FindToFinalPattern

			mockConfig.Request.Regex.Count += len(variables)
		}
	}
}

func processBodyVariables(mockConfig *domain.MockConfiguration) {

	mockBody, err := request.MockBodyToString(mockConfig.Request.Body)
	if err != nil {
		return
	}

	found, variables := regex.FindStringValuesRegex(regex.FindBodyVariablePattern, mockBody)
	if !found {
		return
	}

	if len(variables) > 0 {
		mockBody = stringutil.RemoveParenthesis(mockBody)

		for _, variable := range variables {
			value := variable[0]

			variableName := analyzeValue(variable[1], variable[2])

			mockBody = strings.Replace(mockBody, value, regex.FindVariableValuePattern, 1)

			addVariableControl(mockConfig, variableName, value, BodyVariable)

			mockConfig.Request.Regex.Count++
		}

		mockConfig.Request.Regex.Body = mockBody + regex.FindToFinalPattern
	}
}

func analyzeValue(value1, value2 string) string {
	if value1 == "" {
		return value2
	}

	return value1
}
