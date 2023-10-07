package mocks

import (
	"mock_server_mux/internal/core/domain"
	"net/http"
)

func GetConfigMock1() domain.MockConfiguration {
	return domain.MockConfiguration{
		Info: domain.Info{
			TestGroup: "Default",
			TestName:  "",
		},
		Request: domain.RequestMock{
			Method: "GET",
			URL:    "/example/[0-9]+/t1+$",
			Headers: []domain.Header{
				{
					Name:  "Content-Type",
					Value: "application/json",
				},
			},
			QueryParameters: []domain.QueryParameter{},
		},
		Response: domain.ResponseMock{
			Configurations: domain.ResponseConfiguration{
				ResponseDelay: 0,
			},
			StatusCode: 200,
			Headers: []domain.Header{
				{
					Name:  "Content-Type",
					Value: "application/json",
				},
			},
			Body: "{'test': 'ok'}",
		},
	}
}

func GetConfigMock2() domain.MockConfiguration {
	return domain.MockConfiguration{
		Info: domain.Info{
			TestGroup: "Default",
			TestName:  "",
		},
		Request: domain.RequestMock{
			Method: "GET",
			URL:    "/example/[A-Z][a-z]+/t2+$",
			Headers: []domain.Header{
				{
					Name:  "Content-Type",
					Value: "application/json",
				},
			},
			QueryParameters: []domain.QueryParameter{},
		},

		Response: domain.ResponseMock{
			Configurations: domain.ResponseConfiguration{
				ResponseDelay: 0,
			},
			StatusCode: http.StatusNotFound,
			Headers: []domain.Header{
				{
					Name:  "Content-Type",
					Value: "application/json",
				},
			},
			Body: "{'error': 'not found'}",
		},
	}
}
