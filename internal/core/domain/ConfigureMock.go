package domain

import (
	"mock_server_mux/api/presenter"
)

type MockConfiguration struct {
	Id       int          `json:"id"`
	Name     string       `json:"name"`
	Info     Info         `json:"info"`
	Request  RequestMock  `json:"request"`
	Response ResponseMock `json:"response"`
}

type Info struct {
	TestGroup string `json:"test_group"`
	TestName  string `json:"test_name"`
}

type RequestMock struct {
	Method               string               `json:"method"`
	URL                  string               `json:"url"`
	Headers              []Header             `json:"headers"`
	QueryParameters      []QueryParameter     `json:"query_parameters"`
	Body                 interface{}          `json:"body"`
	RequestConfiguration RequestConfiguration `json:"request_configuration"`
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type QueryParameter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type RequestConfiguration struct {
	Forward    bool   `json:"forward"`
	ForwardUrl string `json:"forward_url"`
}

type ResponseMock struct {
	Configurations ResponseConfiguration `json:"configurations"`
	StatusCode     int                   `json:"status"`
	Headers        []Header              `json:"headers"`
	Body           interface{}           `json:"body"`
}

type ResponseConfiguration struct {
	ResponseDelay int `json:"response_delay"`
}

func (mc MockConfiguration) ToPresenter() presenter.MockConfiguration {
	return presenter.MockConfiguration{
		Id:       mc.Id,
		Name:     mc.Name,
		Info:     mc.Info.ToPresenter(),
		Request:  mc.Request.ToPresenter(),
		Response: mc.Response.ToPresenter(),
	}
}

func (i Info) ToPresenter() presenter.Info {
	return presenter.Info{
		TestGroup: i.TestGroup,
		TestName:  i.TestName,
	}
}

func (rm RequestMock) ToPresenter() presenter.RequestMock {
	requestMock := presenter.RequestMock{
		Method:               rm.Method,
		URL:                  rm.URL,
		Body:                 rm.Body,
		RequestConfiguration: rm.RequestConfiguration.ToPresenter(),
	}

	for _, header := range rm.Headers {
		requestMock.Headers = append(requestMock.Headers, header.ToPresenter())
	}

	for _, queryParameter := range rm.QueryParameters {
		requestMock.QueryParameters = append(requestMock.QueryParameters, queryParameter.ToPresenter())
	}

	return requestMock
}

func (rc RequestConfiguration) ToPresenter() presenter.RequestConfiguration {
	return presenter.RequestConfiguration{
		Forward:    rc.Forward,
		ForwardUrl: rc.ForwardUrl,
	}
}

func (h Header) ToPresenter() presenter.Header {
	return presenter.Header{
		Key:   h.Key,
		Value: h.Value,
	}
}

func (qp QueryParameter) ToPresenter() presenter.QueryParameter {
	return presenter.QueryParameter{
		Key:   qp.Key,
		Value: qp.Value,
	}
}

func (rm ResponseMock) ToPresenter() presenter.ResponseMock {
	responseMock := presenter.ResponseMock{
		Configurations: rm.Configurations.ToPresenter(),
		StatusCode:     rm.StatusCode,
		Body:           rm.Body,
	}

	for _, header := range rm.Headers {
		responseMock.Headers = append(responseMock.Headers, header.ToPresenter())
	}

	return responseMock
}

func (rc ResponseConfiguration) ToPresenter() presenter.ResponseConfiguration {
	return presenter.ResponseConfiguration{
		ResponseDelay: rc.ResponseDelay,
	}
}
