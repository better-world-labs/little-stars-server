package entities

import "aed-api-server/internal/pkg/global"

type Medal struct {
	ID              int64  `json:"id,string"`
	Name            string `json:"name"`
	Order           int    `json:"order"`
	ActiveIcon      string `json:"activeIcon"`
	InactiveIcon    string `json:"inactiveIcon"`
	Description     string `json:"description"`
	ShareBackground string `json:"shareBackground"`
}

// TODO json
type UserMedal struct {
	ID         int64                `xorm:"id pk autoincr" json:"id,string"`
	MedalID    int64                `xorm:"medal_id" json:"medalId,string"`
	UserID     int64                `xorm:"user_id" json:"userId"`
	BusinessId string               `xorm:"business_id" json:"businessId"`
	Created    global.FormattedTime `xorm:"created" json:"created"`
}

const (
	MedalIdSaveLife      = 1
	MedalIdFirstDonation = 2
	MedalIdInspector     = 3
)
