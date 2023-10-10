package requestfilter

import (
	"context"
	"mock_server_mux/internal/core/domain"
	"mock_server_mux/pkg/apperrors"
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

		result, found = filterByQueryHeader(request, result)
		if !found {
			continue
		}

		//result, found = s.filterByQueryBody(body, result)

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

	if findStringRegex(configuration.Request.Method+" "+configuration.Request.URL,
		method+" "+path) {
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

func filterByQueryHeader(request *http.Request, configuration domain.MockConfiguration) (domain.MockConfiguration, bool) {

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
