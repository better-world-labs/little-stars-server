package entities

import "time"

type Essay struct {
	ID         int64     `xorm:"id pk autoincr" json:"id"`
	Title      string    `json:"title" binding:"required"`
	Type       int       `json:"type" binding:"required,min=1,max=2"`
	FrontCover []string  `json:"frontCover" binding:"required"`
	Content    string    `json:"content" binding:"required"`
	CreateAt   time.Time `json:"createAt"`
	Extra      string    `json:"extra"`
	Sort       int64     `xorm:"sort"`
}
