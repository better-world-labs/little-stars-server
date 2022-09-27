package entities

import (
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"time"
)

const (
	SourceWxMenu      = "wx-menu"      // 小程序菜单来源，sharer 为 用户 ID 或者 Openid
	SourceWxMini      = "wx-mini"      // 小程序来源，sharer 为 用户 ID 或者 Openid
	SourceOfficial    = "official"     // 官方来源，sharer 为 oxxx
	SourceWxP         = "wx-p"         // 公众号来源，sharer 为 oxxx
	SourcePlacard     = "placard"      // 海报来源，sharer 为 用户 ID 或者 Openid
	SourceCmsPlacard  = "cms-placard"  // 海报来源，sharer 为 用户 ID 或者 Openid
	SourceWxCircle    = "wx-circle"    // 朋友圈来源，sharer 为 用户 ID 或者 Openid
	SourceCharityCard = "charity-card" // 公益名牌来源，sharer 为 用户 ID 或者 Openid
)

type (
	UserIdTag  int64
	OpenIdTag  string
	NotUserTag string

	Trace struct {
		From     string    `json:"from"     xorm:"from"`
		To       string    `json:"to"       xorm:"to"`
		Sharer   string    `json:"sharer"   xorm:"sharer"`
		OpenID   string    `json:"openid"   xorm:"open_id"`
		DeviceID string    `json:"deviceId" xorm:"device_id"`
		Source   string    `json:"source"   xorm:"source"`
		CreateAt time.Time `json:"createAt" xorm:"created_at"`
	}

	CreateQrCodeReq struct {
		Source   string `json:"source"`
		Sharer   string `json:"sharer"`
		PagePath string `json:"pagePath"`
	}

	CreateQrCodeRes struct {
		Image string `json:"image"`
	}
)

func (t *Trace) GetSharerTag() (interface{}, error) {
	switch t.Source {
	case SourceOfficial, SourceWxP:
		return NotUserTag(t.Sharer), nil

	case SourcePlacard,
		SourceCmsPlacard,
		SourceWxCircle,
		SourceWxMenu,
		SourceCharityCard,
		SourceWxMini:
		match, err := regexp.Match("^[0-9]+$", []byte(t.Sharer))
		if err != nil {
			return nil, err
		}

		if match {
			intSharer, err := strconv.ParseInt(t.Sharer, 10, 64)
			if err != nil {
				return nil, err
			}

			return UserIdTag(intSharer), nil
		}

		return OpenIdTag(t.Sharer), nil

	default:
		return nil, errors.New("invalid source")
	}
}

func (*Trace) Decode(b []byte) (emitter.DomainEvent, error) {
	var t Trace
	return &t, json.Unmarshal(b, &t)
}

func (t *Trace) Encode() ([]byte, error) {
	return json.Marshal(t)
}
