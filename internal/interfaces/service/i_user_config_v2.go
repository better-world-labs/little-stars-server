package service

import (
	"aed-api-server/internal/interfaces/entities"
)

type IUserConfigV2 interface {
	Get(userId int64, key string) (*entities.UserConfig, error)
	GetOrDefault(userId int64, key string, defaultValue interface{}) (*entities.UserConfig, error)
	Put(userId int64, key string, value interface{}) (updated bool, err error)
}
