package user

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/db"
	"encoding/json"
)

type configV2 struct {
}

//go:inject-component
func NewUserConfigV2Service() *configV2 {
	return &configV2{}
}

func (configV2) Get(userId int64, key string) (*entities.UserConfig, error) {
	var configV2 entities.UserConfig
	existed, err := db.Table("user_config_v2").
		Where("user_id = ? and `key` = ?", userId, key).Get(&configV2)
	if err != nil {
		return nil, err
	}
	if existed {
		return &configV2, nil
	}
	return nil, nil
}

func (v configV2) GetOrDefault(userId int64, key string, defaultValue interface{}) (*entities.UserConfig, error) {
	value, err := v.Get(userId, key)
	if err != nil {
		return nil, err
	}

	if value == nil {
		_, err := v.Put(userId, key, map[string]interface{}{
			"value": defaultValue,
		})
		if err != nil {
			return nil, err
		}
	}

	return v.Get(userId, key)
}

func (configV2) Put(userId int64, key string, configV2 interface{}) (updated bool, err error) {
	json, err := json.Marshal(configV2)
	if err != nil {
		return false, err
	}

	rst, e := db.Exec("insert into user_config_v2(user_id, `key`, `value`, updated_at, created_at)"+`
			values(?, ?, ?, now(), now())
			ON DUPLICATE KEY
			UPDATE value = ?
		`, userId, key, json, json)
	affected, e := rst.RowsAffected()
	if e == nil {
		updated = affected > 0
	}

	return updated, err
}
