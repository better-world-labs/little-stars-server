package service

import (
	"aed-api-server/internal/interfaces/entities"
	"time"
)

type SubscribeMsg interface {
	Report(userId int64, key string, templates []*entities.SubscribeTemplateSetting, setting *entities.SubscriptionsSetting) error
	GetLastReport(userId int64, key string) (templates []*entities.SubscribeTemplateSetting, setting *entities.SubscriptionsSetting, reportAt *time.Time, err error)
	Send(userId int64, openId string, msgKey entities.SubscribeMessageKey, params interface{}) (bool, error)
}
