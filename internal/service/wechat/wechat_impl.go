package wechat

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/config"
	"aed-api-server/internal/pkg/crypto"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"errors"
	"fmt"
	"github.com/imroc/req"
	log "github.com/sirupsen/logrus"
)

type wechatImpl struct {
	crypt *crypto.WXUserDataCrypt

	Conf        config.WechatOAuthConfig `conf:"wechat"`
	TokenServer config.Service           `conf:"wechat-token"`
}

//go:inject-component
func NewWechat() service.IWechat {
	return &wechatImpl{crypt: crypto.NewWXUserDataCrypt()}
}

func (w wechatImpl) CodeToSession(code string) (*entities.WechatCode2SessionRes, error) {
	log.Info("[wechat]", "code2Session: code=", code)
	param := req.QueryParam{"appid": w.Conf.AppID,
		"secret":     w.Conf.AppSecret,
		"grant_type": "authorization_code",
		"js_code":    code,
	}

	resp, err := req.Get(w.Conf.Server+"/sns/jscode2session", param)
	if err != nil {
		return nil, err
	}

	if str, err := resp.ToString(); err == nil {
		log.Debugf("[wechat] code2Session response: %s", str)
	}

	var res entities.WechatCode2SessionRes
	if err = resp.ToJSON(&res); err != nil {
		return nil, err
	}

	if res.ErrCode != 0 {
		return nil, errors.New(res.ErrMsg)
	}

	return &res, nil
}

func (w wechatImpl) Decrypt(encryptedData, iv, sessionKey string, dst interface{}) error {
	_, err := w.crypt.Decrypt(encryptedData, iv, sessionKey, dst)
	return err
}

func (w wechatImpl) GetWalks(req *entities.WechatDataDecryptReq) (*entities.WechatWalkData, error) {
	session, err := w.CodeToSession(req.Code)
	if err != nil {
		return nil, err
	}

	var data entities.WechatWalkData
	b, err := w.crypt.Decrypt(req.EncryptedData, req.Iv, session.SessionKey, &data)

	log.Info("GetWalks: walkData=%v", string(b))
	return &data, err
}

//GetMinaToken 获取小程序的accessToken
func (w wechatImpl) GetMinaToken() (string, error) {
	serverConf := interfaces.GetConfig().WechatToken
	appId := interfaces.GetConfig().Wechat.AppID
	resp, err := req.Get(fmt.Sprintf("http://%s:%d/api/wechat/token?appId=%s", serverConf.Host, serverConf.Port, appId))
	if err != nil {
		return "", err
	}

	var r response.Response
	err = resp.ToJSON(&r)
	if err != nil {
		return "", err
	}

	if token, ok := r.Data.(string); ok {
		return token, nil
	}

	return "", errors.New("invalid token type")
}

//SendSubscribeMsg 发送订阅消息
func (w wechatImpl) SendSubscribeMsg(msgKey entities.SubscribeMessageKey, openId, templateId, templateEl string, params interface{}) (*entities.WechatRst, error) {
	token, err := interfaces.S.Wx.GetMinaToken()
	if err != nil {
		log.Error("SendSubscribeMsg failed for get token error:", err)
		return nil, err
	}

	config := interfaces.GetConfig()

	state := "formal"
	if config.Server.Env == "dev" || config.Server.Env == "local" {
		state = "developer"
	}
	if config.Server.Env == "test" {
		state = "trial"
	}

	sprintf := fmt.Sprintf("/pages/index/index?templateId=%s&templateKey=%s&templateEl=%s", templateId, msgKey, templateEl)
	var body = map[string]interface{}{
		"touser":            openId,
		"template_id":       templateId,
		"page":              sprintf,
		"data":              params,
		"miniprogram_state": state,
	}

	var rst entities.WechatRst

	err = utils.Post("https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token="+token, body, &rst)
	if err != nil {
		log.Error("send subscribe msg error", err)
		return nil, err
	}

	log.Info("send subscribe msg rst:", rst)
	return &rst, nil
}

func (w *wechatImpl) MiniProgramCode2Session(code string, encryptPhone string, iv string, data *entities.WechatMiniProgramRes) error {
	log.Debugf("code2Session: code=%s", code)

	session, err := w.CodeToSession(code)
	if err != nil {
		return err
	}

	if encryptPhone != "" {
		var phone crypto.WxUserPhone
		_, err = w.crypt.Decrypt(encryptPhone, iv, session.SessionKey, &phone)
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
