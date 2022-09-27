package mock

import (
	"aed-api-server/internal/interfaces/entities"
	"errors"
	"time"
)

type WechatMock struct {
}

func NewWechatMock() *WechatMock {
	return &WechatMock{}
}

func (w WechatMock) GetMinaToken() (string, error) {
	return "token", nil
}

func (w WechatMock) SendSubscribeMsg(openId, templateId, page string, params interface{}) (*entities.WechatRst, error) {
	return nil, nil
}

func (w WechatMock) CodeToSession(code string) (*entities.WechatCode2SessionRes, error) {
	return &entities.WechatCode2SessionRes{
		WechatRst: entities.WechatRst{
			ErrCode: 0,
		},
		SessionKey: "z81Q1p9jPDumiyz2Vi33zQ==",
		Openid:     "oyL7e5eIlledFGNfu1XDLlSzdmxU",
	}, nil
}

func (w WechatMock) MiniProgramCode2Session(code string, encryptPhone string, iv string, data *entities.WechatMiniProgramRes) error {
	return errors.New("not implements")
}

func (w WechatMock) GetUserEncryptKey(openid, sessionKey string) ([]*entities.WechatEncryptKey, error) {
	return []*entities.WechatEncryptKey{
		{
			EncryptKey: sessionKey,
			Version:    3,
			ExpireIn:   3600,
			Iv:         "raUOWdl0H3/ORd9wSbKrRQ==",
			CreateTime: time.Now().UnixMilli() / 1000,
		},
		{
			EncryptKey: sessionKey,
			Version:    2,
			ExpireIn:   3600,
			Iv:         "raUOWdl0H3/ORd9wSbKrRQ==",
			CreateTime: time.Now().UnixMilli() / 1000,
		},
		{
			EncryptKey: sessionKey,
			Version:    1,
			ExpireIn:   3600,
			Iv:         "raUOWdl0H3/ORd9wSbKrRQ==",
			CreateTime: time.Now().UnixMilli() / 1000,
		},
	}, nil

}

func (w WechatMock) GetWalks(req *entities.WechatDataDecryptReq) (*entities.WechatWalkData, error) {
	panic("implement me")
}
