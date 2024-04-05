package requestfilter

import (
	"context"
	"mock_server_mux/internal/core/domain"
	"mock_server_mux/pkg/apperrors"
	"mock_server_mux/pkg/logger"
	"mock_server_mux/pkg/utils/regex"
	requestutil "mock_server_mux/pkg/utils/request"
	stringutil "mock_server_mux/pkg/utils/string"
	"net/http"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) FilterMockHandlersByRequest(ctx context.Context, request *http.Request, mockConfigurations []domain.MockConfiguration) (domain.MockConfiguration, error) {

	var filteredMockConfigurations []domain.MockConfiguration

	for _, mockConfig := range mockConfigurations {
		result, found := filterByMethodAndPath(request, mockConfig)
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

	return getBestResult(filteredMockConfigurations), nil
}

func filterByMethodAndPath(request *http.Request, configuration domain.MockConfiguration) (domain.MockConfiguration, bool) {
	path := request.URL.Path
	method := request.Method

	if configuration.Request.Method == method && configuration.Request.URL == path {
		return configuration, true
	}

	if len(configuration.Request.Regex.URL) > 0 {
		pattern := configuration.Request.Method + " " + configuration.Request.Regex.URL

		findString := method + " " + path

		if regex.FindStringRegex(pattern, findString) {
			return configuration, true
		}
	}

	return domain.MockConfiguration{}, false
}

func filterByQueryParam(request *http.Request, configuration domain.MockConfiguration) (domain.MockConfiguration, bool) {

	if len(configuration.Request.QueryParameters) == 0 {
		return configuration, true
	}

	for key, value := range configuration.Request.QueryParameters {

		queryParamValue := request.URL.Query().Get(key)

		if queryParamValue != value {

			if regexValue, ok := configuration.Request.Regex.QueryParameters[key]; ok {

				if !regex.FindStringRegex(regexValue, queryParamValue) {
					return domain.MockConfiguration{}, false
				}

			} else {
				return domain.MockConfiguration{}, false
			}
		}
	}

	return configuration, true
}

func filterByHeader(request *http.Request, configuration domain.MockConfiguration) (domain.MockConfiguration, bool) {

	if len(configuration.Request.Headers) == 0 {
		return configuration, true
	}

	for key, value := range configuration.Request.Headers {

		headerValue := request.Header.Get(key)

		if headerValue != value {

			if regexValue, ok := configuration.Request.Regex.Headers[key]; ok {

				if !regex.FindStringRegex(regexValue, headerValue) {
					return domain.MockConfiguration{}, false
				}

			} else {
				return domain.MockConfiguration{}, false
			}
		}
	}

	return configuration, true

}

func filterByBody(request *http.Request, configuration domain.MockConfiguration) (domain.MockConfiguration, bool) {

	requestBody, err := requestutil.ReadBodyToString(request)
	if err != nil {
		logger.Error("Error reading body", err)
		return domain.MockConfiguration{}, false
	}

	requestBody = prepareBody(requestBody, false)

	mockBody, err := requestutil.MockBodyToString(configuration.Request.Body)
	if err != nil {
		logger.Error("Error encoding to JSON", err)
		return domain.MockConfiguration{}, false
	}

	mockBody = prepareBody(mockBody, false)

	if requestBody != mockBody {

		if configuration.Request.Regex.Body == nil {
			return domain.MockConfiguration{}, false
		}

		mockBody = prepareBody(configuration.Request.Regex.Body.(string), false)

		if !regex.FindStringRegex(mockBody, requestBody) {
			return domain.MockConfiguration{}, false
		}
	}

	return configuration, true
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

// Best Result Order
// 1. Specialized with greater ID
// 2. Specialized
// 3. Generic with less regex variables
// 4. Generic with greater ID
// 5. Generic
func getBestResult(mockConfigs []domain.MockConfiguration) domain.MockConfiguration {

	bestIndex := -1

	for index, mockConfig := range mockConfigs {

		if bestIndex == -1 {
			bestIndex = index
		}

		if mockConfig.Request.Regex.Count == 0 {
			if mockConfig.Id > mockConfigs[bestIndex].Id {
				bestIndex = index
			}

		} else {

			if mockConfig.Request.Regex.Count < mockConfigs[bestIndex].Request.Regex.Count {
				bestIndex = index
			} else {
				if mockConfig.Request.Regex.Count == mockConfigs[bestIndex].Request.Regex.Count &&
					mockConfig.Id > mockConfigs[bestIndex].Id {
					bestIndex = index
				}
			}
		}
	}

	return mockConfigs[bestIndex]
}
