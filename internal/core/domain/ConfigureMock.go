package domain

type MockType int

type Variable struct {
	Name      string `json:"name,omitempty"`
	ValueFrom string `json:"value_from,omitempty"`
}

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
	Method          string               `json:"method"`
	URL             string               `json:"url"`
	Regex           Regex                `json:"regex"`
	Headers         []Header             `json:"headers"`
	QueryParameters []QueryParameter     `json:"query_parameters"`
	Body            interface{}          `json:"body"`
	Variables       []Variable           `json:"variables"`
	Configuration   RequestConfiguration `json:"configuration"`
}

type Regex struct {
	URL             string           `json:"url,omitempty"`
	Headers         []Header         `json:"headers,omitempty"`
	QueryParameters []QueryParameter `json:"query_parameters,omitempty"`
	Body            interface{}      `json:"body,omitempty"`
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type QueryParameter struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type RequestConfiguration struct {
	Forward    bool   `json:"forward"`
	ForwardUrl string `json:"forward_url"`
}

type ResponseMock struct {
	StatusCode    int                   `json:"status"`
	Headers       []Header              `json:"headers"`
	Body          interface{}           `json:"body"`
	Configuration ResponseConfiguration `json:"configuration"`
}

type ResponseConfiguration struct {
	ResponseDelay int `json:"response_delay"`
}
