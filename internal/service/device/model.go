package device

import "aed-api-server/internal/pkg/global"

type DeviceGuide struct {
	Id        int64 `json:"-"`
	Uid       string
	AccountId int64
	DeviceId  string
	Remark    string
	Desc      string
	Pic       string
	Created   global.FormattedTime `json:"time,omitempty" xorm:"created"`
}

const (
	GalleryTypeDevice = 1
	GalleryTypePicket = 2
)
