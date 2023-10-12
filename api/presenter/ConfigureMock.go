package presenter

type MockType int

type MockConfiguration struct {
	Id       int          `json:"id"`
	Info     Info         `json:"info"`
	Request  RequestMock  `json:"request"`
	Response ResponseMock `json:"response"`
}

type Info struct {
	TestName  string `json:"test_name"`
	TestGroup string `json:"test_group"`
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
