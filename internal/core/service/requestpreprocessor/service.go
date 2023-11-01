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
)

func (s *Service) ProcessMockParameters(ctx context.Context, mockConfig domain.MockConfiguration) (domain.MockConfiguration, error) {

	userVariables := mockConfig.Variables

	mockConfig.Variables = nil

	processUserVariables(&mockConfig, userVariables)

	processUrlVariables(&mockConfig)

	processQueryVariables(&mockConfig)

	processHeaderVariables(&mockConfig)

	processBodyVariables(&mockConfig)

	return mockConfig, nil
}

func processVariable(mockConfig *domain.MockConfiguration, variable string, value interface{}, context string) {
	if mockConfig.Variables == nil {
		mockConfig.Variables = make(map[string][]domain.Variable)
	}

	mockConfig.Variables[context] = append(mockConfig.Variables[context],
		domain.Variable{
			Name:  "${" + context + "." + variable + "}",
			Value: value,
		})
}

func processUserVariables(mockConfig *domain.MockConfiguration, userVariables map[string][]domain.Variable) {
	for userContext := range userVariables {

		for _, variable := range userVariables[userContext] {

			processVariable(mockConfig, variable.Name, variable.Value, userContext)

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

			processVariable(mockConfig, variable[1], variable[0], PathVariable)
		}
		mockConfig.Request.Regex.URL = URL + regexutil.FindToFinalPattern
	}
}

func processQueryVariables(mockConfig *domain.MockConfiguration) {
	for _, query := range mockConfig.Request.QueryParameters {

		queryValue, _ := interfaceutils.GetToString(query.Value)

		found, variables := regexutil.FindStringValuesRegex(regexutil.FindVariablePattern, queryValue)
		if !found {
			continue
		}

		if len(variables) > 0 {
			for _, variable := range variables {
				queryValue = strings.Replace(queryValue, variable[0], regexutil.FindVariableValuePattern, 1)

				processVariable(mockConfig, variable[1], variable[0], QueryVariable)
			}

			mockConfig.Request.Regex.QueryParameters = append(mockConfig.Request.Regex.QueryParameters,
				domain.QueryParameter{
					Key:   query.Key,
					Value: queryValue,
				})
		}
	}
}

func processHeaderVariables(mockConfig *domain.MockConfiguration) {
	for _, header := range mockConfig.Request.Headers {

		headerValue, _ := interfaceutils.GetToString(header.Value)

		found, variables := regexutil.FindStringValuesRegex(regexutil.FindVariablePattern, headerValue)
		if !found {
			continue
		}

		if len(variables) > 0 {
			for _, variable := range variables {
				headerValue = strings.Replace(headerValue, variable[0], regexutil.FindVariableValuePattern, 1)

				processVariable(mockConfig, variable[1], variable[0], HeaderVariable)
			}

			mockConfig.Request.Regex.Headers = append(mockConfig.Request.Regex.Headers,
				domain.Header{
					Key:   header.Key,
					Value: headerValue,
				})
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

			processVariable(mockConfig, variable[1], variable[0], BodyVariable)
		}

		mockConfig.Request.Regex.Body = body + regexutil.FindToFinalPattern
	}
}
