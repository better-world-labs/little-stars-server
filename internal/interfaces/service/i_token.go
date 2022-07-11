package service

type TokenService interface {
	Generate() string
	ValidateToken(token, value string) (bool, error)
	PutToken(token, value string) error
	RemoveToken(token string) (int64, error)
}
