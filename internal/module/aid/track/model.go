package track

import (
	"time"
)

type Model struct {
	ID               int64     `xorm:"id pk autoincr" json:"id"`
	HelpInfoID       int64     `xorm:"help_info_id"`
	UserID           int64     `xorm:"user_id"`
	DeviceGot        bool      `xorm:"device_got"`
	DeviceGotTime    time.Time `xorm:"device_got_time"`
	SceneArrived     bool      `xorm:"scene_arrived"`
	SceneArrivedTime time.Time `xorm:"scene_arrived_time"`
}
