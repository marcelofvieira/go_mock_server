package apperrors

import (
	"errors"
	"fmt"
)

var (
	DefaultError = Define("default_error")
)

type Error struct {
	code       string
	message    string
	cause      error
	data       interface{}
	stacktrace frame
}

func (e Error) Error() string {
	return e.message
}

// Define define a new error base model.
func Define(code string) Error {
	return Error{code: code}
}

// New creates a new error.
func New(theType Error, cause error, message string) error {
	return newError(theType, cause, message, nil)
}

// Newf creates a new error with message formatting.
func Newf(theType Error, cause error, message string, args ...interface{}) error {
	return newError(theType, cause, fmt.Sprintf(message, args...), nil)
}

// NewWithData creates a new error with data.
func NewWithData(theType Error, cause error, message string, data interface{}) error {
	return newError(theType, cause, message, data)
}

// Wrap wraps an error with a message.
func Wrap(e error, message string) error {
	return wrap(e, message, nil)
}

// Wrapf wraps an error with a formatting message.
func Wrapf(e error, message string, args ...interface{}) error {
	return wrap(e, fmt.Sprintf(message, args...), nil)
}

// WrapWithData wraps an error and add extra data
func WrapWithData(e error, message string, data interface{}) error {
	return wrap(e, message, data)
}

// Is verify if a given error has the same time of the given target error.
// The target parameter should be an error previously defined with the Define function.
func Is(e error, target error) bool {
	var err Error
	if errors.As(e, &err) {
		var targetErr Error
		if errors.As(target, &targetErr) {
			return err.code == targetErr.code
		}
	}

	return errors.Is(target, e)
}

// Code retrieves the error internal code of a given error.
func Code(e error) string {
	var err Error
	if errors.As(e, &err) {
		return err.code
	}

	return ""
}

// Data retrieves the data of a given error or nil if it do not have such data.
func Data(e error) interface{} {
	var err Error
	if errors.As(e, &err) {
		return err.data
	}

	return nil
}

// String returns a string containing all the subbasement information about the given error.
func String(e error) string {
	if e == nil {
		return ""
	}

	var err Error
	if errors.As(e, &err) {
		r := err.message + " | [err_code: " + err.code + "]"

		r = fmt.Sprintf("%v | SRC: %v:%v", r, err.stacktrace.file(), err.stacktrace.line())

		if err.cause != nil {
			r = fmt.Sprintf("%v | CAUSE: {%v}", r, String(err.cause))
		}

		return r
	}

	return e.Error()
}

func newError(theType Error, cause error, message string, data interface{}) error {
	return Error{code: theType.code, message: message, cause: cause, data: data, stacktrace: stackLevel(2)}
}

func wrap(e error, message string, data interface{}) error {
	var err Error
	if errors.As(e, &err) {
		return Error{code: err.code, message: message, cause: e, data: data, stacktrace: stackLevel(2)}
	}

	return Error{code: DefaultError.code, message: message, cause: e, data: data, stacktrace: stackLevel(2)}
}
