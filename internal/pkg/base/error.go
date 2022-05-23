package base

import "fmt"

type Err interface {
	error

	Cause() error
}

func WrapError(module string, message string, cause error) Err {
	return &err{
		module:  module,
		message: message,
		cause:   cause,
	}
}

func NewError(module string, message string) Err {
	return &err{
		module:  module,
		message: message,
	}
}

type err struct {
	cause   error
	message string
	module  string
}

func (e err) Error() string {
	if e.cause == nil {
		return fmt.Sprintf("%s: %s", e.module, e.message)
	}

	return fmt.Sprintf("%s: %s, cause by %#v", e.module, e.message, e.cause)
}

func (e err) Cause() error {
	return e.cause
}
