package user

// TODO 整理实体关系，让组件更通用
type OAuthInfoDTO struct {
	OAuthAccessTokenDTO

	UnionID    string `json:"unionid"`
	Nickname   string `json:"nickname"`
	HeadImgURL string `json:"headimgurl"`
}

type OAuthAccessTokenDTO struct {
	OAuthResponseBaseDTO

	AccessToken string `json:"access_token"`
	OpenID      string `json:"openid"`
}
type Code2SessionRes struct {
	SessionKey string `json:"session_key"`
	Openid     string `json:"openid"`
	UnionId    string `json:"unionid"`
}

type MiniProgramResponseDTO struct {
	OAuthResponseBaseDTO

	OpenID         string `json:"openid"`
	UnionID        string `json:"unionid"`
	SessionKey     string `json:"session_key"`
	DecryptedPhone string `json:"-"`
}

type OAuthResponseBaseDTO struct {
	ErrorCode int `json:"errcode"`
}
