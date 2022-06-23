package entities

import "encoding/json"

const (
	DeviceSourceLocal    = 0
	DeviceSourceImported = 1
)

type (
	AddDevice struct {
		Longitude        float64     `json:"longitude,omitempty"`
		Latitude         float64     `json:"latitude,omitempty"`
		Address          string      `json:"address,omitempty"`
		Contract         string      `json:"contract,omitempty"`
		Title            string      `json:"detailAddress,omitempty"`
		DeviceImage      string      `json:"deviceImage,omitempty"`
		EnvironmentImage string      `json:"environmentImage,omitempty"`
		State            int         `json:"state,omitempty"`
		Guide            string      `json:"guide,omitempty"`
		GuideInfo        []GuideInfo `json:"info"`
		OpenIn           *TimeRange  `json:"openIn" binding:"required"`
	}

	DeviceGuideList struct {
		DeviceId string                `json:"deviceId" form:"deviceId"`
		Info     []DeviceGuideListItem `json:"info"`
	}

	DeviceGuideListItem struct {
		Uid       string      `json:"uid,omitempty"`
		AccountId int64       `json:"accountId,omitempty"`
		UserName  string      `json:"userName,omitempty"`
		Avatar    string      `json:"avatar,omitempty"`
		Info      []GuideInfo `json:"info,omitempty"`
		Time      string      `json:"time"`
	}

	GuideInfo struct {
		Desc   string   `json:"desc"`
		Pic    []string `json:"pic"`
		Remark string   `json:"remark"`
	}

	Gallery struct {
		Type int    `xorm:"type" json:"type"`
		Url  string `xorm:"url" json:"url"`
	}

	Device struct {
		Id               string        `json:"id" xorm:"id pk"`
		Address          string        `json:"address"`
		Title            string        `json:"detailAddress"` // 详细地址
		Icon             string        `xorm:"-" json:"icon"`
		ClockInImage     string        `xorm:"clock_in_image" json:"-"`
		Longitude        float64       `json:"longitude"`
		Latitude         float64       `json:"latitude"`
		Distance         int64         `json:"distance" xorm:"-"`
		Tel              string        `json:"contract"`
		DeviceImage      string        `json:"deviceImage,omitempty" xorm:"origin"`
		EnvironmentImage string        `json:"environmentImage,omitempty" xorm:"env_origin"`
		State            int           `json:"state,omitempty"`
		CredibleState    int           `json:"credibleState" xorm:"credible_state"`
		Created          int64         `json:"created" xorm:"created"`
		Source           int           `json:"-" xorm:"source"`
		CreateBy         int64         `json:"createBy" xorm:"create_by"`
		OpenIn           TimeRange     `json:"openIn" xorm:"open_in"`
		Inspector        []*SimpleUser `json:"inspector" xorm:"-"`
	}

	PicketedDeviceCount struct {
		CredibleState int
		Count         int
	}

	TimeRange struct {
		Week      []int8 `json:"week" binding:"required"`
		BeginTime string `json:"beginTime" binding:"required"`
		EndTime   string `json:"endTime" binding:"required"`
	}
)

func (t *TimeRange) FromDB(b []byte) error {
	return json.Unmarshal(b, t)
}

func (t *TimeRange) ToDB() ([]byte, error) {
	return json.Marshal(t)
}
