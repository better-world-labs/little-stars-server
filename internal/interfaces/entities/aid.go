package entities

import (
	"aed-api-server/internal/pkg/location"
	"time"
)

type HelpInfo struct {
	*location.Coordinate `xorm:"extends"`

	ID            int64     `xorm:"id pk autoincr" json:"id"`
	Address       string    `xorm:"address" json:"address"`
	DetailAddress string    `xorm:"detail_address" json:"detailAddress"`
	Publisher     int64     `xorm:"publisher" json:"publisher"`
	PublishTime   time.Time `xorm:"publish_time" json:"publishTime"`
	Images        []string  `xorm:"images" json:"images"`
	Exercise      bool      `xorm:"exercise" json:"exercise"`
}
