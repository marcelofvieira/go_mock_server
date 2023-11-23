package variableprocessor

import (
	"context"
	"encoding/json"
	"mock_server_mux/internal/core/domain"
	"mock_server_mux/pkg/logger"
	"mock_server_mux/pkg/regexutil"
	"mock_server_mux/pkg/stringutils"
	"net/http"
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

	found, variables := regexutil.FindStringValuesRegex(regexutil.FindVariablePattern, mockConfig.Request.URL)
	if !found {
		return
	}

	found, variablesValue := regexutil.FindStringValuesRegex(mockConfig.Request.Regex.URL, URL)
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

		found, variables := regexutil.FindStringValuesRegex(regexutil.FindVariablePattern, mockConfig.Request.QueryParameters[key].(string))
		if !found {
			return
		}

		found, variablesValue := regexutil.FindStringValuesRegex(regexParamValue, requestParamValue)
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

		found, variables := regexutil.FindStringValuesRegex(regexutil.FindVariablePattern, mockConfig.Request.Headers[key].(string))
		if !found {
			return
		}

		found, variablesValue := regexutil.FindStringValuesRegex(regexHeaderValue, requestHeaderValue)
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

	requestBody := mockConfig.Request.PreparedBody

	jsonBytes, err := json.Marshal(mockConfig.Request.Body)
	if err != nil {
		logger.Error("Error encoding to JSON", err)
		return
	}

	mockBody := prepareBody(string(jsonBytes))

	found, variables := regexutil.FindStringValuesRegex(regexutil.FindBodyVariablePattern, mockBody)
	if !found {
		return
	}

	found, variablesValue := regexutil.FindStringValuesRegex(mockConfig.Request.Regex.Body.(string), requestBody.(string))
	if !found {
		return
	}

	for i := 1; i < len(variablesValue[0]); i++ {
		variableName := "${" + BodyVariable + "." + variables[i-1][1] + "}"
		mockConfig.MockVariables[BodyVariable][variableName] = variablesValue[0][i]
	}

}

func prepareBody(body string) string {

	body = stringutils.ReplaceTabsToSpaces(body)
	body = stringutils.ReplaceNewLinesToSpaces(body)
	body = stringutils.RemoveSpaces(body)

	if body == "null" {
		return ""
	}

	return body
}
