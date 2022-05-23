package exception

import (
	"go.uber.org/zap/buffer"
)

type ServiceError struct {
	Code    int
	Message string
}

func NewServiceError(code int, message string) *ServiceError {
	return &ServiceError{Code: code, Message: message}
}

func (e *ServiceError) Error() string {
	b := buffer.Buffer{}
	b.AppendString(": ")
	b.AppendString(e.Message)
	return b.String()
}
