package service

type UserConfigService interface {
	GetConfig(userId int64, key string) (string, error)
	PutConfig(userId int64, key string, config string) (updated bool, err error)

	GetConfigToValue(userId int64, key string, value interface{}) error
	PutValueToConfig(userId int64, key string, value interface{}) (updated bool, err error)
}
