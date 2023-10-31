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

	processUrlConfiguration(&mockConfig)

	processQueryParameter(&mockConfig)

	return mockConfig, nil
}

func processVariable(mockConfig *domain.MockConfiguration, variable, value, context string) {
	mockConfig.Request.Variables = append(mockConfig.Request.Variables,
		domain.Variable{
			Name:      "${" + context + "." + variable + "}",
			ValueFrom: value,
		})
}

func processUrlConfiguration(mockConfig *domain.MockConfiguration) {
	URL := mockConfig.Request.URL

	found, variables := regexutil.FindStringValuesRegex(regexutil.FindVariablePattern, mockConfig.Request.URL)
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

func processQueryParameter(mockConfig *domain.MockConfiguration) {
	for _, query := range mockConfig.Request.QueryParameters {

		queryValue, _ := interfaceutils.GetToString(query.Value)

		found, variables := regexutil.FindStringValuesRegex(regexutil.FindVariablePattern, queryValue)
		if !found {
			return
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
