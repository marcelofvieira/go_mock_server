package domain

const (
	ConfigResponseDelay = "response_delay"
)

type MockConfiguration struct {
	Id            int                               `json:"id,omitempty"`
	Info          Info                              `json:"info,omitempty"`
	Request       RequestMock                       `json:"request,omitempty"`
	Response      ResponseMock                      `json:"response,omitempty"`
	MockVariables map[string]map[string]interface{} `json:"mock_variables,omitempty"`
	UserVariables map[string]map[string]interface{} `json:"variables,omitempty"`
}

type Info struct {
	TestName  string `json:"test_name,omitempty"`
	TestGroup string `json:"test_group,omitempty"`
}

type RequestMock struct {
	Method          string                 `json:"method,omitempty"`
	URL             string                 `json:"url,omitempty"`
	Regex           Regex                  `json:"regex,omitempty"`
	Headers         map[string]interface{} `json:"headers,omitempty"`
	QueryParameters map[string]interface{} `json:"query_parameters,omitempty"`
	Body            interface{}            `json:"body,omitempty"`
	Configuration   map[string]interface{} `json:"configuration,omitempty"`
}

type Regex struct {
	URL             string            `json:"url,omitempty"`
	Headers         map[string]string `json:"headers,omitempty"`
	QueryParameters map[string]string `json:"query_parameters"`
	Body            interface{}       `json:"body,omitempty"`
}

type ResponseMock struct {
	StatusCode    int                    `json:"status,omitempty"`
	Headers       map[string]interface{} `json:"headers,omitempty"`
	Body          interface{}            `json:"body,omitempty"`
	Configuration map[string]interface{} `json:"configuration,omitempty"`
}
