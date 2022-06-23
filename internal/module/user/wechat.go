package user

import (
	"aed-api-server/internal/pkg/config"
	"aed-api-server/internal/pkg/crypto"
	"errors"
	"github.com/imroc/req"
	log "github.com/sirupsen/logrus"
)

type WechatClient interface {
	MiniProgramCode2Session(code string, encryptPhone string, iv string, data *MiniProgramResponseDTO) error
	GetAccessToken(code string, data *OAuthAccessTokenDTO) error
	GetUserInfo(accessToken string, openID string, data *OAuthInfoDTO) error
	CodeToSession(code string) (*Code2SessionRes, error)
	Decrypt(encryptedData, iv, sessionKey string, dst interface{}) error
}

func NewWechatClient(config *config.WechatOAuthConfig) WechatClient {
	return &wechatClient{
		Config: config,
		crypt:  crypto.NewWXUserDataCrypt(config.AppID),
	}
}

type wechatClient struct {
	Config *config.WechatOAuthConfig
	crypt  *crypto.WXUserDataCrypt
}

// GetAccessToken 根据授权码获取 OAuth AccessToken
// @param code 授权码
// @param data 绑定结果的指针
// @return OAuthInfoDTo, error
func (c *wechatClient) GetAccessToken(code string, data *OAuthAccessTokenDTO) error {
	log.Debugf("get access token: code=%s", code)
	param := req.QueryParam{
		"appid":      c.Config.AppID,
		"secret":     c.Config.AppSecret,
		"grant_type": "authorization_code",
		"code":       code,
	}

	if data == nil {
		log.Errorf("bind data with nil pointer")
		return errors.New("bind data with nil pointer")
	}

	resp, err := req.Get(c.Config.Server+"/sns/oauth2/access_token", param)
	if err != nil {
		log.Errorf("get access token error: %v", err)
		return err
	}

	if str, err := resp.ToString(); err == nil {
		log.Debugf("response: %s", str)
	}

	if err = resp.ToJSON(data); err != nil {
		log.Errorf("get access token error: %v", err)
		return err
	}

	if data.ErrorCode != 0 {
		log.Errorf("get access token with error code: %d", data.ErrorCode)
		return errors.New("get access token with error code")
	}

	return nil
}

func (c *wechatClient) CodeToSession(code string) (*Code2SessionRes, error) {
	log.Info("[wechat]", "code2Session: code=", code)
	param := req.QueryParam{"appid": c.Config.AppID,
		"secret":     c.Config.AppSecret,
		"grant_type": "authorization_code",
		"js_code":    code,
	}

	resp, err := req.Get(c.Config.Server+"/sns/jscode2session", param)
	if err != nil {
		return nil, err
	}

	if str, err := resp.ToString(); err == nil {
		log.Debugf("[wechat] code2Session response: %s", str)
	}

	var res Code2SessionRes
	if err = resp.ToJSON(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *wechatClient) Decrypt(encryptedData, iv, sessionKey string, dst interface{}) error {
	_, err := c.crypt.Decrypt(encryptedData, iv, sessionKey, dst)
	return err
}

func (c *wechatClient) MiniProgramCode2Session(code string, encryptPhone string, iv string, data *MiniProgramResponseDTO) error {
	log.Debugf("code2Session: code=%s", code)

	session, err := c.CodeToSession(code)
	if err != nil {
		return err
	}

	if encryptPhone != "" {
		var phone crypto.WxUserPhone
		_, err = c.crypt.Decrypt(encryptPhone, iv, session.SessionKey, &phone)
		if err != nil {
			return err
		}
		data.DecryptedPhone = phone.PhoneNumber
	}

	data.OpenID = session.Openid
	data.SessionKey = session.SessionKey
	data.UnionID = session.UnionId

	return nil
}

// GetUserInfo 读取用户信息
// @param accessToken OAuth Access Token
// @param openID OpenID
// @param data 绑定结果的指针
// @return error
func (c *wechatClient) GetUserInfo(accessToken string, openID string, data *OAuthInfoDTO) error {
	log.Debugf("get user info: openID=%s, accessToken=%s", openID, accessToken)

	if data == nil {
		log.Errorf("bind data with nil pointer")
		return errors.New("bind data with nil pointer")
	}

	resp, err := req.Get(c.Config.Server+"/sns/userinfo", req.QueryParam{
		"access_token": accessToken,
		"openid":       openID,
	})

	if err != nil {
		log.Errorf("get user info error: %v", err)
		return err
	}

	if str, err := resp.ToString(); err == nil {
		log.Debugf("response: %s", str)
	}

	if err = resp.ToJSON(data); err != nil {
		log.Errorf("get user info error: %v", err)
		return err
	}

	if data.ErrorCode != 0 {
		log.Errorf("get user info with error code: %d", data.ErrorCode)
		return errors.New("get user info with error code")
	}

	return nil
}
