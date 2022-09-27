package service

import (
	"aed-api-server/internal/interfaces/entities"
	"io"
)

type IWechat interface {
	// GetMinaToken 获取小程序的accessToken
	GetMinaToken() (string, error)

	// SendSubscribeMsg 发送订阅消息
	SendSubscribeMsg(openId, templateId, page string, params interface{}) (*entities.WechatRst, error)

	// CodeToSession 授权码换取 SessionKey
	CodeToSession(code string) (*entities.WechatCode2SessionRes, error)

	// MiniProgramCode2Session  小程序登陆换取 SessionKey
	MiniProgramCode2Session(code string, encryptPhone string, iv string, data *entities.WechatMiniProgramRes) error

	// GetUserEncryptKey 读取用户最近三个的密钥
	GetUserEncryptKey(openid, sessionKey string) ([]*entities.WechatEncryptKey, error)

	// GenericUrlLink 获取小程序链接
	GenericUrlLink(path, query string) (string, error)

	//GenMinaCode 生成小程序码
	GenMinaCode(path string, width int, autoColor, isHyaline bool, lineColor string) (reader io.ReadCloser, contentType string, err error)
}
