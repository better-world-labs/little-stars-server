package service

type ConfigService interface {
	GetConfig(key string) (string, error)
	PutConfig(key string, config string) error
	GetAllConfig() (string, error)
}
