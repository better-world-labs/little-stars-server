package entities

import (
	"aed-api-server/internal/pkg/location"
	"fmt"
	"time"
)

const (
	UserConfigKeyFirstEnterAEDMap          = "first-enter-aed-map"
	UserConfigKeyFirstEnterCommunity       = "first-enter-community"
	UserConfigKeySubscribeOfficialAccounts = "subscribe-official-accounts"
)

type LoginCommandV2 struct {
	Code     string `json:"code"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatarUrl"`
}

type MobileCommand struct {
	EncryptPhone string `json:"encryptedMobile"`
	Code         string `json:"code"`
	Iv           string `json:"iv"`
}

// v1.10 废弃
type LoginCommand struct {
	MobileCode   string `json:"mobileCode"`
	Code         string `json:"code"`
	EncryptPhone string `json:"encryptedMobile"`
	Iv           string `json:"iv"`
	Nickname     string `json:"nickname"`
	Avatar       string `json:"avatarUrl"`
}

type SimpleLoginCommand struct {
	Code string `json:"code" binding:"required"`
}

type AccountDTOWithSessionKey struct {
	UserDTO

	SessionKey string `json:"sessionKey"`
}

type User struct {
	ID         int64     `xorm:"id pk autoincr"`
	Nickname   string    `xorm:"nickname"`
	Uid        string    `xorm:"uid"`
	Status     int       `xorm:"status"`
	Avatar     string    `xorm:"avatar"`
	Mobile     string    `xorm:"mobile"`
	Unionid    string    `xorm:"unionid"`
	SessionKey string    `xorm:"session_key"`
	Openid     string    `xorm:"openid"`
	Created    time.Time `xorm:"created"`
}

func (a User) ToSimple() *SimpleUser {
	return &SimpleUser{
		ID:       a.ID,
		Nickname: a.Nickname,
		Avatar:   a.Avatar,
	}
}

func (a User) Available() bool {
	return a.Avatar != "" || a.Nickname != ""
}

type Position struct {
	*location.Coordinate `xorm:"extends"`

	ID        int64 `xorm:"id pk"`
	AccountID int64 `xorm:"account_id"`
}
type UserStat struct {
	TotalCount  int64
	MobileCount int64
}

type UserDTO struct {
	ID       int64  `xorm:"id pk autoincr" json:"id"`
	Uid      string `json:"uid"`
	Nickname string `json:"nickname"`
	Token    string `xorm:"-" json:"token,omitempty"`
	Avatar   string `json:"avatarUrl"`
	Mobile   string `json:"mobile"`
	Openid   string `json:"openid"`
}

type SimpleUser struct {
	ID       int64  `json:"id" xorm:"id"`
	Nickname string `json:"nickname" xorm:"nickname"`
	Avatar   string `json:"avatarUrl" xorm:"avatar"`
	Uid      string `json:"uid" xorm:"uid"`
	Mobile   string `json:"mobile"`
}

type UserAboutStat struct {
	UserAboutDonationsCount int `json:"userAboutDonationsCount"`
	UserAboutFeedsCount     int `json:"userAboutFeedsCount"`
	UserAboutSosCount       int `json:"userAboutSosCount"`
}

type SubscribeTemplateSetting struct {
	TemplateId string `json:"templateId"`
	Status     string `json:"status"`
}

type SubscriptionsSetting struct {
	MainSwitch bool                        `json:"mainSwitch"`
	Templates  []*SubscribeTemplateSetting `json:"templates"`
}

type UserEventType string

const (
	UserEventTypeGetTreeInfo UserEventType = "get-tree-info"
	UserEventTypeGetWalkStep UserEventType = "get-walk-step"
)

func GetUserEventTypeOfReport(key string) UserEventType {
	return UserEventType(fmt.Sprintf("report:%s", key))
}

const (
	Report_enterAedMap         = "enter-aed-map"
	Report_readNews            = "read-news"
	Report_showSubscribeQrCode = "show-subscribe-qr"
	Report_openTreasureChest   = "open-treasure-chest"
	Report_scanPage            = "scan-page"
	Report_scanVideo           = "scan-video"
	Report_feedLike            = "feed-like"
	Report_feedLikeCancel      = "feed-like-cancel"
	Report_feedShare           = "feed-shared"
	Report_viewFeedPage        = "view-feed-page"
)
