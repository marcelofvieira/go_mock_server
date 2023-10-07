package apperrors

var (
	NotFound            = Define("not_found")
	IllegalOperation    = Define("illegal_operation")
	InvalidInput        = Define("invalid_input")
	InternalServerError = Define("internal")
	BadRequest          = Define("bad_request")
	NotImplemented      = Define("not_implemented")
)
