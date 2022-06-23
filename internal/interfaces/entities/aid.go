package entities

import (
	"aed-api-server/internal/pkg/global"
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

type PublishDTO struct {
	Address       string   `json:"address"`
	DetailAddress string   `json:"detailAddress"`
	Longitude     float64  `json:"longitude"`
	Latitude      float64  `json:"latitude"`
	Images        []string `json:"images"`
}

type ActionDTO struct {
	*location.Coordinate

	AidID string `json:"aidId"`
}

func FromImageModel(model *HelpImage) string {
	return model.Origin
}

func FromImageModels(models []*HelpImage) []string {
	dtos := make([]string, len(models))

	for i := range models {
		dtos[i] = FromImageModel(models[i])
	}

	return dtos
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
	ID       int64                `json:"id,string"`
	UserName string               `json:"userName"`
	Class    string               `json:"class"`
	Record   interface{}          `json:"record,omitempty"`
	Created  global.FormattedTime `json:"created"`
}

type HelpInfoDTO struct {
	location.Coordinate

	ID              int64    `json:"id,string"`
	Address         string   `json:"address"`
	DetailAddress   string   `json:"detailAddress"`
	PublishTime     string   `json:"publishTime"`
	PublisherName   string   `json:"publisherName"`
	PublisherMobile string   `json:"publisherMobile"`
	Images          []string `json:"images"`
	Exercise        bool     `json:"exercise"`
}

type Call120RequestDto struct {
	MobileLast4 string `json:"mobileLast4"`
}

type HelpImage struct {
	ID         int64     `xorm:"id pk autoincr"`
	HelpInfoID int64     `xorm:"help_info_id index"`
	Origin     string    `xorm:"origin"`
	Thumbnail  string    `xorm:"thumbnail"`
	Created    time.Time `xorm:"created"`
}

type HelpInfoImaged struct {
	*HelpInfo

	images []*HelpImage
}
