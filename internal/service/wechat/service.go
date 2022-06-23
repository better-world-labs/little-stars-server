package wechat

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"errors"
	"fmt"
	"github.com/imroc/req"
	log "github.com/sirupsen/logrus"
)

type srv struct{}

func NewWechatSrv() *srv {
	return &srv{}
}

//GetMinaToken 获取小程序的accessToken
func (*srv) GetMinaToken() (string, error) {
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
func (*srv) SendSubscribeMsg(openId string, templateId string, params interface{}) (*entities.WechatRst, error) {
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

	var body = map[string]interface{}{
		"touser":            openId,
		"template_id":       templateId,
		"page":              "/pages/index/index",
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
