package point

import "aed-api-server/internal/pkg/global"

type Point struct {
	Id          int64                `json:"id,omitempty"`
	AccountId   int64                `json:"account,omitempty"`
	Points      float64              `json:"points,omitempty"`
	Description string               `json:"description,omitempty"`
	CreateAt    global.FormattedTime `json:"time,omitempty"`
	Class       string               `json:"-"`
	Extra       string               `json:"-"`
}
