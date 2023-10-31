package presenter

type MockType int

type Variable struct {
	Name      string `json:"name,omitempty"`
	ValueFrom string `json:"value_from,omitempty"`
}

type MockConfiguration struct {
	Id       int          `json:"id,omitempty"`
	Info     Info         `json:"info,omitempty"`
	Request  RequestMock  `json:"request,omitempty"`
	Response ResponseMock `json:"response,omitempty"`
}

type Info struct {
	TestName  string `json:"test_name,omitempty"`
	TestGroup string `json:"test_group,omitempty"`
}

type RequestMock struct {
	Method          string               `json:"method,omitempty"`
	URL             string               `json:"url,omitempty"`
	Headers         []Header             `json:"headers,omitempty"`
	QueryParameters []QueryParameter     `json:"query_parameters,omitempty"`
	Body            interface{}          `json:"body,omitempty"`
	Variables       []Variable           `json:"variables,omitempty"`
	Configuration   RequestConfiguration `json:"configuration,omitempty"`
	Regex           Regex                `json:"regex,omitempty"`
}

type Regex struct {
	URL             string           `json:"url,omitempty"`
	Headers         []Header         `json:"headers,omitempty"`
	QueryParameters []QueryParameter `json:"query_parameters,omitempty"`
	Body            interface{}      `json:"body,omitempty"`
}

type Header struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type QueryParameter struct {
	Key   string      `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

type RequestConfiguration struct {
	Forward    bool   `json:"forward,omitempty"`
	ForwardUrl string `json:"forward_url,omitempty"`
}

type ResponseMock struct {
	StatusCode    int                   `json:"status,omitempty"`
	Headers       []Header              `json:"headers,omitempty"`
	Body          interface{}           `json:"body,omitempty"`
	Configuration ResponseConfiguration `json:"configuration,omitempty"`
}

type ResponseConfiguration struct {
	ResponseDelay int `json:"response_delay,omitempty"`
}
