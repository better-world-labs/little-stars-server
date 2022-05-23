package device

import (
	"aed-api-server/internal/interfaces/entities"
	page "aed-api-server/internal/pkg/query"
)

type DelDevice struct {
	UdId string `json:"id,omitempty"`
}

type ListDevice struct {
	Longitude float64 `form:"longitude,omitempty"`
	Latitude  float64 `form:"latitude,omitempty"`
	Distance  float64 `form:"distance,omitempty"`
	page.Query
}

type InfoDevice struct {
	UdId      string  `form:"id,omitempty"`
	Longitude float64 `form:"longitude,omitempty"`
	Latitude  float64 `form:"latitude,omitempty"`
}

type UpdateDevice struct {
	UdId             string  `json:"id,omitempty"`
	Address          string  `json:"address"`
	Title            string  `json:"detailAddress"`
	Longitude        float64 `json:"longitude"`
	Latitude         float64 `json:"latitude"`
	DeviceImage      string  `json:"deviceImage"`
	EnvironmentImage string  `json:"environmentImage"`
	State            int     `json:"state"`
	Guide            string  `json:"guide,omitempty"`
}

type AddDeviceGuideDto struct {
	DeviceId string               `json:"deviceId"`
	Info     []entities.GuideInfo `json:"info"`
}

type DeviceGuideDto struct {
	Uid string `json:"uid" form:"uid"`
}

type CorrectDeviceGuideDto struct {
	DeviceId string               `json:"deviceId"`
	Info     []entities.GuideInfo `json:"info"`
}
