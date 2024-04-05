package variableprocessor

import (
	"context"
	"mock_server_mux/internal/core/domain"
	"mock_server_mux/pkg/logger"
	"mock_server_mux/pkg/utils/regex"
	requestutil "mock_server_mux/pkg/utils/request"
	stringutil "mock_server_mux/pkg/utils/string"
	"net/http"
	"strings"
)

const (
	PathVariable   = "path"
	QueryVariable  = "query"
	HeaderVariable = "header"
	BodyVariable   = "body"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetVariablesValues(ctx context.Context, request *http.Request, mockConfig domain.MockConfiguration) (domain.MockConfiguration, error) {

	getURLVariablesValues(request, &mockConfig)

	getQueryParamVariablesValues(request, &mockConfig)

	getHeaderVariablesValues(request, &mockConfig)

	getBodyVariablesValues(request, &mockConfig)

	return mockConfig, nil

}

func getURLVariablesValues(request *http.Request, mockConfig *domain.MockConfiguration) {

	if len(mockConfig.MockVariables[PathVariable]) == 0 {
		return
	}

	URL := request.URL.Path

	found, variables := regex.FindStringValuesRegex(regex.FindVariablePattern, mockConfig.Request.URL)
	if !found {
		return
	}

	found, variablesValue := regex.FindStringValuesRegex(mockConfig.Request.Regex.URL, URL)
	if !found {
		return
	}

	for i := 1; i < len(variablesValue[0]); i++ {
		variableName := "${" + PathVariable + "." + variables[i-1][1] + "}"
		mockConfig.MockVariables[PathVariable][variableName] = variablesValue[0][i]
	}

}

func getQueryParamVariablesValues(request *http.Request, mockConfig *domain.MockConfiguration) {

	if len(mockConfig.MockVariables[QueryVariable]) == 0 {
		return
	}

	for key, regexParamValue := range mockConfig.Request.Regex.QueryParameters {

		requestParamValue := request.URL.Query().Get(key)

		found, variables := regex.FindStringValuesRegex(regex.FindVariablePattern, mockConfig.Request.QueryParameters[key].(string))
		if !found {
			return
		}

		found, variablesValue := regex.FindStringValuesRegex(regexParamValue, requestParamValue)
		if !found {
			return
		}

		for i := 1; i < len(variablesValue[0]); i++ {
			variableName := "${" + QueryVariable + "." + variables[i-1][1] + "}"
			mockConfig.MockVariables[QueryVariable][variableName] = variablesValue[0][i]
		}
	}
}

func getHeaderVariablesValues(request *http.Request, mockConfig *domain.MockConfiguration) {

	if len(mockConfig.MockVariables[HeaderVariable]) == 0 {
		return
	}

	for key, regexHeaderValue := range mockConfig.Request.Regex.Headers {

		requestHeaderValue := request.Header.Get(key)

		found, variables := regex.FindStringValuesRegex(regex.FindVariablePattern, mockConfig.Request.Headers[key].(string))
		if !found {
			return
		}

		found, variablesValue := regex.FindStringValuesRegex(regexHeaderValue, requestHeaderValue)
		if !found {
			return
		}

		for i := 1; i < len(variablesValue[0]); i++ {
			variableName := "${" + HeaderVariable + "." + variables[i-1][1] + "}"
			mockConfig.MockVariables[HeaderVariable][variableName] = variablesValue[0][i]
		}
	}
}

func getBodyVariablesValues(request *http.Request, mockConfig *domain.MockConfiguration) {

	if len(mockConfig.MockVariables[BodyVariable]) == 0 {
		return
	}

	requestBody, err := requestutil.ReadBodyToString(request)
	if err != nil {
		logger.Error("Error reading body", err)
		return
	}

	requestBody = prepareBody(requestBody, false)

	mockBody, err := requestutil.MockBodyToString(mockConfig.Request.Body)
	if err != nil {
		return
	}

	mockBody = prepareBody(mockBody, true)

	found, variables := regex.FindStringValuesRegex(regex.FindBodyVariablePattern, mockBody)
	if !found {
		return
	}

	found, variablesValue := regex.FindStringValuesRegex(prepareBody(mockConfig.Request.Regex.Body.(string), false), requestBody)
	if !found {
		return
	}

	for i := 1; i < len(variablesValue[0]); i++ {
		variableName := "${" + BodyVariable + "." + analyseValue(variables[i-1][1], variables[i-1][2]) + "}"
		mockConfig.MockVariables[BodyVariable][variableName] = strings.Replace(variablesValue[0][i], "\"", "", -1)
	}

}

func prepareBody(body string, workParenthesis bool) string {

	body = stringutil.ReplaceTabsToSpaces(body)
	body = stringutil.ReplaceNewLinesToSpaces(body)
	body = stringutil.RemoveSpaces(body)

	if workParenthesis {
		body = stringutil.RemoveParenthesis(body)
	}

	if body == "null" {
		return ""
	}

	return body
}

func analyseValue(value1, value2 string) string {
	if value1 == "" {
		return value2
	} else {
		return value1
	}
}
