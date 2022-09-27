package entities

import (
	"aed-api-server/internal/pkg/global"
	"time"
)

const (
	FlowDirectionIn  FlowDirection = "in"
	FlowDirectionOut FlowDirection = "out"
)

type (
	FlowDirection string
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

type AwardFlowQueryCommand struct {
	UserId    int64  `form:"userId"`
	Keyword   string `form:"keyword"`
	Direction string `form:"direction"`
}

type AwardPointFlow struct {
	Id            int64         `json:"id"` //积分流水ID
	UserId        int64         `json:"-"`
	Points        int           `json:"points"` //积分数值
	Direction     FlowDirection `json:"direction"`
	Description   string        `json:"description"`
	PeckExpiredAt time.Time     `json:"peckExpiredAt"` //过期时间
	CreatedAt     time.Time     `json:"createdAt"`     //创建时间
	PeckedAt      *time.Time    `json:"peckedAt"`      //创建时间
	Status        int           `json:"status"`
}

type UserPointsRecord struct {
	Id          string     `json:"id"`
	Description string     `json:"description"`
	Points      int        `json:"points"`
	Time        *time.Time `json:"time"` // *global.FormattedTime 会在查询引发 panic
}

type Point struct {
	Id          int64                `json:"id,omitempty"`
	AccountId   int64                `json:"account,omitempty"`
	Points      float64              `json:"points,omitempty"`
	Description string               `json:"description,omitempty"`
	CreateAt    global.FormattedTime `json:"time,omitempty"`
	Class       string               `json:"-"`
	Extra       string               `json:"-"`
}
