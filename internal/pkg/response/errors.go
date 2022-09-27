package response

type AuthorizationError struct {
	Message string
}

func (e AuthorizationError) Error() string {
	return e.Message
}

var (
	ErrorInvalidToken = &AuthorizationError{Message: "invalid token"}
)
