package user

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/db"
	"encoding/json"
	"github.com/go-xorm/xorm"
)

type config struct {
	service.UserConfigService
}

//go:inject-component
func NewUserConfigService() *config {
	return &config{}
}

func (config) GetConfig(userId int64, key string) (*entities.UserConfigDO, error) {
	var config entities.UserConfigDO
	existed, err := db.Table("user_config").Where("user_id = ? and `key` = ?", userId, key).Get(&config)
	if err != nil {
		return nil, err
	}
	if existed {
		return &config, nil
	}
	return nil, nil
}

func (config) PutConfig(userId int64, key string, config string) (updated bool, err error) {
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

func (s config) GetConfigToValue(userId int64, key string, value interface{}) error {
	config, err := s.GetConfig(userId, key)
	if err != nil {
		return err
	}

	if config == nil {
		return nil
	}

	if config.Value != "" {
		return json.Unmarshal([]byte(config.Value), &value)
	}
	return nil
}

func (s config) PutValueToConfig(userId int64, key string, value interface{}) (updated bool, err error) {
	jsonStr, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	return s.PutConfig(userId, key, string(jsonStr))
}
