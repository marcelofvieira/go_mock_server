package requestfilter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mock_server_mux/internal/core/domain"
	"mock_server_mux/pkg/apperrors"
	"mock_server_mux/pkg/stringutils"
	"net/http"
	"regexp"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) FilterMockHandlersByRequest(ctx context.Context, request *http.Request, mockConfigurations []domain.MockConfiguration) (domain.MockConfiguration, error) {

	var filteredMockConfigurations []domain.MockConfiguration

	for _, mockConfig := range mockConfigurations {
		result, found := filterByMethodAndLiteralPath(request, mockConfig)
		if !found {
			continue
		}

		result, found = filterByQueryParam(request, result)
		if !found {
			continue
		}

		result, found = filterByHeader(request, result)
		if !found {
			continue
		}

		result, found = filterByBody(request, result)
		if !found {
			continue
		}

		filteredMockConfigurations = append(filteredMockConfigurations, result)
	}

	if len(filteredMockConfigurations) == 0 {
		return domain.MockConfiguration{}, apperrors.New(apperrors.NotFound, nil, "not found handler")
	}

	return filteredMockConfigurations[0], nil
}

func filterByMethodAndLiteralPath(request *http.Request, configuration domain.MockConfiguration) (domain.MockConfiguration, bool) {

	path := request.URL.Path
	method := request.Method

	if configuration.Request.Method == method && configuration.Request.URL == path {
		return configuration, true
	}

	pattern := configuration.Request.Method + " " + configuration.Request.URL
	findString := method + " " + path

	if findStringRegex(pattern, findString) {
		return configuration, true
	}

	return domain.MockConfiguration{}, false
}

func filterByQueryParam(request *http.Request, configuration domain.MockConfiguration) (domain.MockConfiguration, bool) {

	if len(configuration.Request.QueryParameters) == 0 {
		return configuration, true
	}

	for _, queryParam := range configuration.Request.QueryParameters {

		queryParamValue := request.URL.Query().Get(queryParam.Key)

		if queryParamValue != queryParam.Value {
			if !findStringRegex(queryParam.Value, queryParamValue) {
				return domain.MockConfiguration{}, false
			}
		}
	}

	return configuration, true
}

func findStringRegex(pattern, text string) bool {
	validRegex := regexp.MustCompile(pattern)

	return validRegex.MatchString(text)
}

func filterByHeader(request *http.Request, configuration domain.MockConfiguration) (domain.MockConfiguration, bool) {

	if len(configuration.Request.Headers) == 0 {
		return configuration, true
	}

	for _, header := range configuration.Request.Headers {

		headerValue := request.Header.Get(header.Key)

		if headerValue != header.Value {

			if !findStringRegex(header.Value, headerValue) {
				return domain.MockConfiguration{}, false
			}
		}
	}

	return configuration, true

}

func filterByBody(request *http.Request, configuration domain.MockConfiguration) (domain.MockConfiguration, bool) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		//TODO: log error
		return domain.MockConfiguration{}, false
	}

	requestBody := prepareBody(string(body))

	jsonBytes, err := json.Marshal(configuration.Request.Body)
	if err != nil {
		fmt.Println("Error encoding to JSON:", err)
	}

	mockBody := prepareBody(string(jsonBytes))

	if requestBody != mockBody {
		if !findStringRegex(mockBody, requestBody) {
			return domain.MockConfiguration{}, false
		}
	}

	return configuration, true
}

func prepareBody(body string) string {

	body = stringutils.ReplaceTabsToSpaces(body)
	body = stringutils.ReplaceNewLinesToSpaces(body)
	body = stringutils.RemoveSpaces(body)

	return body
}