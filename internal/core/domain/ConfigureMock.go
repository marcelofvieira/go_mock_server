package domain

type MockType int

type Variable struct {
	Name        string
	ValueBefore string
	ValueAfter  string
}

const (
	Literal MockType = iota
	Regex   MockType = iota
)

type MockConfiguration struct {
	Id       int          `json:"id"`
	Info     Info         `json:"info"`
	Request  RequestMock  `json:"request"`
	Response ResponseMock `json:"response"`
}

type Info struct {
	TestName  string   `json:"test_name"`
	TestGroup string   `json:"test_group"`
	MockType  MockType `json:"mock_type"`
}

type RequestMock struct {
	Method          string               `json:"method"`
	URL             string               `json:"url"`
	Headers         []Header             `json:"headers"`
	QueryParameters []QueryParameter     `json:"query_parameters"`
	Body            interface{}          `json:"body"`
	Variables       []Variable           `json:"variables"`
	Configuration   RequestConfiguration `json:"configuration"`
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
	StatusCode    int                   `json:"status"`
	Headers       []Header              `json:"headers"`
	Body          interface{}           `json:"body"`
	Configuration ResponseConfiguration `json:"configuration"`
}

type ResponseConfiguration struct {
	ResponseDelay int `json:"response_delay"`
}
