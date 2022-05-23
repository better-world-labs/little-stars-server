package response

type IllegalArgumentError struct {
	message string
}

func NewIllegalArgumentError(message string) error {
	return &IllegalArgumentError{}
}

func (e IllegalArgumentError) Error() string {
	return e.message
}
