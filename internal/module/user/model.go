package user

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/location"
	"time"
)

type User struct {
	ID       int64     `xorm:"id pk autoincr"`
	Nickname string    `xorm:"nickname"`
	Uid      string    `xorm:"uid"`
	Avatar   string    `xorm:"avatar"`
	Mobile   string    `xorm:"mobile"`
	Unionid  string    `xorm:"unionid"`
	Openid   string    `xorm:"openid"`
	Created  time.Time `xorm:"created"`
}

func (a User) ToSimple() *entities.SimpleUser {
	return &entities.SimpleUser{
		ID:       a.ID,
		Nickname: a.Nickname,
		Avatar:   a.Avatar,
	}
}

type Position struct {
	*location.Coordinate `xorm:"extends"`

	ID        int64 `xorm:"id pk"`
	AccountID int64 `xorm:"account_id"`
}
