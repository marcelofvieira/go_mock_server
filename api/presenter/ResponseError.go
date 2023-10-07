package presenter

type ResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
