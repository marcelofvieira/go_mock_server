package presenter

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
	Body                 string               `json:"body"`
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
