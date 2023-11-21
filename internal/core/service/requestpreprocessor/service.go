package requestpreprocessor

import (
	"context"
	"mock_server_mux/internal/core/domain"
	"mock_server_mux/pkg/interfaceutils"
	"mock_server_mux/pkg/regexutil"
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
	UserVariable   = "var"
)

func (s *Service) ProcessMockParameters(ctx context.Context, mockConfig domain.MockConfiguration) (domain.MockConfiguration, error) {

	processMockVariables(&mockConfig)

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

func processMockVariables(mockConfig *domain.MockConfiguration) {
	mockVariables := mockConfig.MockVariables

	mockConfig.MockVariables = nil

	for contextKey, contextVariables := range mockVariables {
		for key, value := range contextVariables {
			addVariableControl(mockConfig, key, value, contextKey)
		}
	}
}

func processUrlVariables(mockConfig *domain.MockConfiguration) {
	URL := mockConfig.Request.URL

	found, variables := regexutil.FindStringValuesRegex(regexutil.FindVariablePattern, URL)
	if !found {
		return
	}

	if len(variables) > 0 {
		for _, variable := range variables {
			URL = strings.Replace(URL, variable[0], regexutil.FindVariableValuePattern, 1)

			addVariableControl(mockConfig, variable[1], variable[0], PathVariable)
		}

		mockConfig.Request.Regex.URL = URL + regexutil.FindToFinalPattern
	}
}

func processQueryVariables(mockConfig *domain.MockConfiguration) {

	for parameterName, parameterValue := range mockConfig.Request.QueryParameters {

		value, _ := interfaceutils.GetToString(parameterValue)

		found, variables := regexutil.FindStringValuesRegex(regexutil.FindVariablePattern, value)
		if !found {
			continue
		}

		if len(variables) > 0 {
			mockConfig.Request.Regex.QueryParameters = make(map[string]string)

			for _, variable := range variables {
				parameterValue = strings.Replace(value, variable[0], regexutil.FindVariableValuePattern, 1)

				addVariableControl(mockConfig, variable[1], variable[0], QueryVariable)
			}

			mockConfig.Request.Regex.QueryParameters[parameterName] = value + regexutil.FindToFinalPattern
		}
	}
}

func processHeaderVariables(mockConfig *domain.MockConfiguration) {

	for headerName, headerValue := range mockConfig.Request.Headers {

		value, _ := interfaceutils.GetToString(headerValue)

		found, variables := regexutil.FindStringValuesRegex(regexutil.FindVariablePattern, value)
		if !found {
			continue
		}

		if len(variables) > 0 {
			mockConfig.Request.Regex.Headers = make(map[string]string)

			for _, variable := range variables {
				value = strings.Replace(value, variable[0], regexutil.FindVariableValuePattern, 1)

				addVariableControl(mockConfig, variable[1], variable[0], HeaderVariable)
			}

			mockConfig.Request.Regex.Headers[headerName] = value + regexutil.FindToFinalPattern
		}
	}
}

func processBodyVariables(mockConfig *domain.MockConfiguration) {

	body, _ := interfaceutils.GetToString(mockConfig.Request.Body)

	found, variables := regexutil.FindStringValuesRegex(regexutil.FindBodyVariablePattern, body)
	if !found {
		return
	}

	if len(variables) > 0 {
		for _, variable := range variables {
			body = strings.Replace(body, variable[0], regexutil.FindVariableValuePattern, 1)

			addVariableControl(mockConfig, variable[1], variable[0], BodyVariable)
		}

		mockConfig.Request.Regex.Body = body + regexutil.FindToFinalPattern
	}
}
