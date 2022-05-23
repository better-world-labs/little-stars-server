package response

import "aed-api-server/internal/pkg/exception"

type aedError struct {
	message string
	code    int
}

func (e *aedError) Error() string {
	return e.message
}

func (e *aedError) Message() string {
	return e.message
}

func (e *aedError) Code() int {
	return e.code
}

type ApiError struct {
	aedError

	httpStatus int
}

func (e *ApiError) HttpStatus() int {
	return e.httpStatus
}

func NewApiError(code int, message string, httpStatus int) *ApiError {
	return &ApiError{
		aedError: aedError{
			message: message,
			code:    code,
		},
		httpStatus: httpStatus,
	}
}

type HTTPServiceError struct {
	*exception.ServiceError

	HttpStatus int
}

func NewHTTPServiceError(httpStatus int, code int, message string) *HTTPServiceError {
	return &HTTPServiceError{
		ServiceError: exception.NewServiceError(code, message),
		HttpStatus:   httpStatus,
	}
}
