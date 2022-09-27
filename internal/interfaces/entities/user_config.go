package entities

import (
	"encoding/json"
	"time"
)

type UserConfigDO struct {
	UserId    int64
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserConfig struct {
	UserId    int64
	Key       string
	Value     map[string]interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c UserConfig) ParseValue(dst interface{}) error {
	value := c.Value["value"]
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, dst)
}
