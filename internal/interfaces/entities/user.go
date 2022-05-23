package entities

const (
	UserConfigKeyFirstEnterAEDMap = "first-enter-aed-map"
)

type UserDTO struct {
	ID       int64  `json:"id,string"`
	Uid      string `json:"uid"`
	Nickname string `json:"nickname"`
	Token    string `json:"token,omitempty"`
	Avatar   string `json:"avatarUrl"`
	Mobile   string `json:"mobile"`
	Openid   string `json:"openid"`
}

type SimpleUser struct {
	ID       int64  `json:"id,string" xorm:"id"`
	Nickname string `json:"nickname" xorm:"nickname"`
	Avatar   string `json:"avatarUrl" xorm:"avatar"`
	Uid      string `json:"uid" xorm:"uid"`
}
