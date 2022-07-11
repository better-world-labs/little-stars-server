package entities

import "time"

type UserConfigDO struct {
	UserId    int64
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserConfigDOV2 struct {
	UserId    int64
	Key       string
	Value     map[string]interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
}
