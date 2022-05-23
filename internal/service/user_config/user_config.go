package user_config

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/db"
	"encoding/json"
	"github.com/go-xorm/xorm"
	"time"
)

func init() {
	if interfaces.S.UserConfig == nil {
		interfaces.S.UserConfig = &Service{}
	}
}

type Service struct {
	service.UserConfigService
}

type UserConfigDO struct {
	UserId    int64
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Service) GetConfig(userId int64, key string) (string, error) {
	var config UserConfigDO
	existed, err := db.Table("user_config").Where("user_id = ? and `key` = ?", userId, key).Get(&config)
	if err != nil {
		return "", err
	}
	if existed {
		return config.Value, nil
	}
	return "", nil
}

func (Service) PutConfig(userId int64, key string, config string) (updated bool, err error) {
	err = db.Transaction(func(db *xorm.Session) error {
		rst, e := db.Exec("insert into user_config(user_id, `key`, `value`, updated_at, created_at)"+`
			values(?, ?, ?, now(), now())
			ON DUPLICATE KEY
			UPDATE value = ?
		`, userId, key, config, config)
		affected, e := rst.RowsAffected()
		if e == nil {
			updated = affected > 0
		}
		return e
	})
	return updated, err
}

func (s Service) GetConfigToValue(userId int64, key string, value interface{}) error {
	config, err := s.GetConfig(userId, key)
	if err != nil {
		return err
	}
	if config != "" {
		return json.Unmarshal([]byte(config), &value)
	}
	return nil
}

func (s Service) PutValueToConfig(userId int64, key string, value interface{}) (updated bool, err error) {
	jsonStr, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	return s.PutConfig(userId, key, string(jsonStr))
}
