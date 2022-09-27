package wechat

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/crypto"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"aed-api-server/internal/pkg/wx_crypto"
	"aed-api-server/internal/server/config"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/imroc/req"
	log "github.com/sirupsen/logrus"
	"io"
	"regexp"
	"strings"
)

type wechatImpl struct {
	Conf        config.WechatOAuthConfig `conf:"wechat"`
	TokenServer config.Service           `conf:"wechat-token"`
}

//go:inject-component
func NewWechat() service.IWechat {
	return &wechatImpl{}
}

func (w *wechatImpl) CodeToSession(code string) (*entities.WechatCode2SessionRes, error) {
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

func (w *wechatImpl) GetUserEncryptKey(openid, sessionKey string) ([]*entities.WechatEncryptKey, error) {
	token, err := w.GetMinaToken()
	if err != nil {
		return nil, err
	}

	signature := crypto.HmacSHA256([]byte(sessionKey), []byte(""))
	sig := strings.ToUpper(hex.EncodeToString(signature))
	log.Infof("GetuserEncryptKey: openid=%s, signature=%s, token=%s, sessionKey=%s", openid, sig, token, sessionKey)

	resp, err := req.Post(fmt.Sprintf("https://api.weixin.qq.com/wxa/business/getuserencryptkey?access_token=%s&openid=%s&signature=%s&sig_method=hmac_sha256", token, openid, sig))
	if err != nil {
		return nil, err
	}

	var r struct {
		Errcode     int                          `json:"errcode"`
		Errmsg      string                       `json:"errmsg"`
		KeyInfoList []*entities.WechatEncryptKey `json:"key_info_list"`
	}

	fmt.Printf("%s\n", resp)
	err = resp.ToJSON(&r)
	if err != nil {
		return nil, err
	}

	if r.Errcode != 0 {
		return nil, errors.New(r.Errmsg)
	}

	return r.KeyInfoList, nil
}

//GetMinaToken 获取小程序的accessToken
func (w *wechatImpl) GetMinaToken() (string, error) {
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

func (w *wechatImpl) GenericUrlLink(path, query string) (string, error) {
	log.Infof("GenericUrlLink: path=%s, query=%s", path, query)

	token, err := interfaces.S.Wx.GetMinaToken()
	if err != nil {
		log.Error("SendSubscribeMsg failed for get token error:", err)
		return "", err
	}

	envVersion := w.GetMinaVersion()

	param := map[string]interface{}{
		"path":            path,
		"query":           query,
		"env_version":     envVersion,
		"expire_type":     1,
		"expire_interval": 30,
	}

	resp, err := req.Post(fmt.Sprintf("https://api.weixin.qq.com/wxa/generate_urllink?access_token=%s", token), req.BodyJSON(param))
	if err != nil {
		return "", err
	}

	var r struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
		UrlLink string `json:"url_link"`
	}

	if resp.Response().StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("request error with status code %d", resp.Response().StatusCode))
	}

	err = resp.ToJSON(&r)
	if err != nil {
		return "", err
	}

	if r.Errcode != 0 {
		return "", errors.New(fmt.Sprintf("failed with errcode=%d, errmsg=%s", r.Errcode, r.Errmsg))
	}

	return r.UrlLink, nil
}

//SendSubscribeMsg 发送订阅消息
func (w *wechatImpl) SendSubscribeMsg(openId, templateId, page string, params interface{}) (*entities.WechatRst, error) {
	token, err := interfaces.S.Wx.GetMinaToken()
	if err != nil {
		log.Error("SendSubscribeMsg failed for get token error:", err)
		return nil, err
	}

	conf := interfaces.GetConfig()

	state := "formal"
	if conf.Server.Env == "dev" || conf.Server.Env == "local" {
		state = "developer"
	}
	if conf.Server.Env == "test" {
		state = "trial"
	}

	var body = map[string]interface{}{
		"touser":            openId,
		"template_id":       templateId,
		"page":              page,
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
		var phone wx_crypto.WxUserPhone
		_, err = wx_crypto.Decrypt(encryptPhone, iv, session.SessionKey, &phone)
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

var json = regexp.MustCompile(`json`)

//GenMinaCode 生成小程序码
func (w *wechatImpl) GenMinaCode(path string, width int, autoColor, isHyaline bool, lineColor string) (reader io.ReadCloser, contentType string, err error) {
	token, err := w.GetMinaToken()
	if err != nil {
		return
	}

	url := "https://api.weixin.qq.com/wxa/getwxacode?access_token=" + token
	envVersion := w.GetMinaVersion()

	param := map[string]interface{}{
		"path":        path,
		"env_version": envVersion,
		"width":       width,
		"auto_color":  autoColor,
		"is_hyaline":  isHyaline,
	}
	if lineColor != "" {
		param["line_color"] = lineColor
	}

	request := req.New()
	request.SetJSONEscapeHTML(false)
	resp, err := request.Post(url, req.BodyJSON(param))
	if err != nil {
		return
	}
	r := resp.Response()

	contentType = r.Header.Get("Content-Type")
	if json.MatchString(contentType) {
		type Err struct {
			ErrCode int    `json:"errCode"`
			ErrMsg  string `json:"errMsg"`
		}

		var e Err

		err = resp.ToJSON(&e)
		if err == nil {
			err = errors.New(e.ErrMsg)
		}
		return
	}
	reader = r.Body
	return
}

func (w *wechatImpl) GetMinaVersion() string {
	conf := interfaces.GetConfig()

	envVersion := "release"
	if conf.Server.Env == "dev" || conf.Server.Env == "local" {
		envVersion = "develop"
	}
	if conf.Server.Env == "test" {
		envVersion = "trial"
	}
	return envVersion
}
