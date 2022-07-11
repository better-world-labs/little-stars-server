package service

import "aed-api-server/internal/interfaces/entities"

type IWechat interface {
	// GetMinaToken 获取小程序的accessToken
	GetMinaToken() (string, error)

	// SendSubscribeMsg 发送订阅消息
	SendSubscribeMsg(msgKey entities.SubscribeMessageKey, openId, templateId, templateEl string, params interface{}) (*entities.WechatRst, error)

	// CodeToSession 授权码换取 SessionKey
	CodeToSession(code string) (*entities.WechatCode2SessionRes, error)

	// MiniProgramCode2Session  小程序登陆换取 SessionKey
	MiniProgramCode2Session(code string, encryptPhone string, iv string, data *entities.WechatMiniProgramRes) error

	// Decrypt 解密数据
	Decrypt(encryptedData, iv, sessionKey string, dst interface{}) error

	// GetWalks  读取近一个月月的步数
	GetWalks(req *entities.WechatDataDecryptReq) (*entities.WechatWalkData, error)
}
