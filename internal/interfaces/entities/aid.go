package entities

import (
	"aed-api-server/internal/pkg/global"
	"aed-api-server/internal/pkg/location"
	"time"
)

type HelpInfo struct {
	location.Coordinate `xorm:"extends"`

	ID            int64     `xorm:"id pk autoincr" json:"id"`
	Address       string    `xorm:"address" json:"address"`
	DetailAddress string    `xorm:"detail_address" json:"detailAddress"`
	Publisher     int64     `xorm:"publisher" json:"publisher"`
	PublishTime   time.Time `xorm:"publish_time" json:"publishTime"`
	Images        []string  `xorm:"images" json:"images"`
	Exercise      bool      `xorm:"exercise" json:"exercise"`
	Practice      bool      `xorm:"practice" json:"practice"`
	NpcId         *int64    `xorm:"npc_id" json:"npcId,omitempty"`
}

type PublishDTO struct {
	Address       string   `json:"address"`
	DetailAddress string   `json:"detailAddress"`
	Longitude     float64  `json:"longitude"`
	Latitude      float64  `json:"latitude"`
	Images        []string `json:"images"`
}

type ActionDTO struct {
	*location.Coordinate

	AidID int64 `json:"aidId"`
}

type DTO struct {
	PublishDTO
}

type HelpInfoComposedDTO struct {
	HelpInfoDTO

	Distance          int64              `json:"distance"`
	NewestActivity    *NewestActivityDTO `json:"newestActivity"`
	FirstDeviceGetter string             `json:"firstDeviceGetter"`
	DeviceGetterCount int                `json:"deviceGetterCount"`
}

type NewestActivityDTO struct {
	ID       int64                `json:"id"`
	UserName string               `json:"userName"`
	Class    string               `json:"class"`
	Record   interface{}          `json:"record,omitempty"`
	Created  global.FormattedTime `json:"created"`
}

type HelpInfoDTO struct {
	location.Coordinate

	ID              int64    `json:"id"`
	Address         string   `json:"address"`
	DetailAddress   string   `json:"detailAddress"`
	PublishTime     string   `json:"publishTime"`
	PublisherName   string   `json:"publisherName"`
	PublisherMobile string   `json:"publisherMobile"`
	Images          []string `json:"images"`
	Exercise        bool     `json:"exercise"`
	Practice        bool     `json:"practice"`
	NpcId           *int64   `json:"npcId,omitempty"`
}

type Call120RequestDto struct {
	MobileLast4 string `json:"mobileLast4"`
}
