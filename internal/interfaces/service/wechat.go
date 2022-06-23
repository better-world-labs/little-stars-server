package service

import "aed-api-server/internal/interfaces/entities"

type WechatClient interface {
	//GetMinaToken 获取小程序的accessToken
	GetMinaToken() (string, error)

	//SendSubscribeMsg 发送订阅消息
	SendSubscribeMsg(openId string, templateId string, params interface{}) (*entities.WechatRst, error)
}
