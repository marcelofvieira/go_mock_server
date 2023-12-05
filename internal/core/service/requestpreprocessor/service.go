package requestpreprocessor

import (
	"context"
	"mock_server_mux/internal/core/domain"
	"mock_server_mux/pkg/regexutil"
	"mock_server_mux/pkg/requestutil"
	"mock_server_mux/pkg/stringutils"
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

		mockConfig.Request.Regex.Count += len(variables)

	}
}

func processQueryVariables(mockConfig *domain.MockConfiguration) {

	for parameterName, parameterValue := range mockConfig.Request.QueryParameters {

		value, _ := parameterValue.(string)

		found, variables := regexutil.FindStringValuesRegex(regexutil.FindVariablePattern, value)
		if !found {
			continue
		}

		if len(variables) > 0 {
			mockConfig.Request.Regex.QueryParameters = make(map[string]string)

			for _, variable := range variables {
				value = strings.Replace(value, variable[0], regexutil.FindVariableValuePattern, 1)

				addVariableControl(mockConfig, variable[1], variable[0], QueryVariable)
			}

			mockConfig.Request.Regex.QueryParameters[parameterName] = value + regexutil.FindToFinalPattern

			mockConfig.Request.Regex.Count += len(variables)
		}
	}
}

func processHeaderVariables(mockConfig *domain.MockConfiguration) {

	for headerName, headerValue := range mockConfig.Request.Headers {

		value, _ := headerValue.(string)

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

			mockConfig.Request.Regex.Count += len(variables)
		}
	}
}

func processBodyVariables(mockConfig *domain.MockConfiguration) {

	mockBody, err := requestutil.MockBodyToString(mockConfig.Request.Body)
	if err != nil {
		return
	}

	found, variables := regexutil.FindStringValuesRegex(regexutil.FindBodyVariablePattern, mockBody)
	if !found {
		return
	}

	if len(variables) > 0 {
		mockBody = stringutils.RemoveParenthesis(mockBody)

		for _, variable := range variables {
			mockBody = strings.Replace(mockBody, variable[0], regexutil.FindVariableValuePattern, 1)

			addVariableControl(mockConfig, variable[1], variable[0], BodyVariable)

			mockConfig.Request.Regex.Count++
		}

		mockConfig.Request.Regex.Body = mockBody + regexutil.FindToFinalPattern
	}
}
