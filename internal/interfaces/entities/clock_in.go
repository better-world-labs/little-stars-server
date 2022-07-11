package entities

import (
	"aed-api-server/internal/pkg/global"
	"time"
)

// ClockInBaseInfo 打卡基本信息
type ClockInBaseInfo struct {
	Id int64 `json:"id,string,omitempty" xorm:"id pk autoincr"`

	DeviceId        string   `json:"deviceId" binding:"required"`
	IsDeviceExisted bool     `json:"isDeviceExisted"` //设备是否存在
	Description     string   `json:"description"`     //打卡描述
	Images          []string `json:"images"`          //打卡图片

	CreatedBy int64     `json:"createdBy,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`

	OpenIn *TimeRange `json:"openIn" xorm:"-"`
}

// ClockIn 打卡信息
type ClockIn struct {
	ClockInBaseInfo

	CreatedBy *SimpleUser          `json:"createdBy"`
	CreatedAt global.FormattedTime `json:"createdAt"`
}

// DeviceClockIn 设备打卡信息
type DeviceClockIn struct {
	Device       *Device    `json:"device"`
	LastClockIns []*ClockIn `json:"lastClockIns"` //最近两次打卡信息

	IsLast2ClockInSame bool `json:"isLast2ClockInSame"` //最近两条打卡信息是否一致

	SupportExistedCount    int `json:"supportExistedCount"`    //设备存在的打卡数量
	SupportNotExistedCount int `json:"supportNotExistedCount"` //设备不存在的打卡数量

	LastSupportExistedUsers    []*SimpleUser `json:"lastSupportExistedUsers"`    //最近两个支持设备存在的人
	LastSupportNotExistedUsers []*SimpleUser `json:"lastSupportNotExistedUsers"` //最近两个支持设备不存在的人
}

// DeviceClockInStat 待打卡设备统计
type DeviceClockInStat struct {
	Total        int64 `json:"total"`        //总设备数
	Todo         int64 `json:"todo"`         //未打卡设备数
	ClockInCount int64 `json:"clockInCount"` //打卡人次
	UserCount    int64 `json:"userCount"`    //参与打卡的人数（去重）
}
