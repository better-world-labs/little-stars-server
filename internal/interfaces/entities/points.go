package entities

import (
	"aed-api-server/internal/pkg/global"
	"time"
)

type PointAddRst struct {
	Point       int    `json:"point"`
	Description string `json:"description"`
}

type UserPoints struct {
	TotalPoints      int `json:"totalPoints"`      //用户总积分
	CumulativePoints int `json:"cumulativePoints"` //用户累计积分
}

type UserPointsFlow struct {
	Id        int64                `json:"id" xorm:"id"`                //积分流水ID
	Name      string               `json:"name" xorm:"name"`            //积分名称
	Points    int                  `json:"points" xorm:"points"`        //积分数值
	ExpiredAt global.FormattedTime `json:"expiredAt" xorm:"expired_at"` //过期时间
}

type UserPointsRecord struct {
	Id          string     `json:"id"`
	Description string     `json:"description"`
	Points      int        `json:"points"`
	Time        *time.Time `json:"time"` // *global.FormattedTime 会在查询引发 panic
}
