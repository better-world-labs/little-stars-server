package test

import (
	"aed-api-server/internal/module/user"
	"aed-api-server/internal/pkg/crypto"
)

type mockWechatClient struct {
	crypto *crypto.WXUserDataCrypt
}

func NewMockedWechatClient() user.WechatClient {
	return &mockWechatClient{
		crypto: crypto.NewWXUserDataCrypt("xxx"),
	}
}

func (m mockWechatClient) Decrypt(encryptedData, iv, sessionKey string, dst interface{}) error {
	_, err := m.crypto.Decrypt(encryptedData, iv, sessionKey, dst)
	return err
}

func (m mockWechatClient) GetMinaToken() (string, error) {
	panic("implement me")
}

func (m mockWechatClient) SendSubscribeMsg(openId string, templateId string, params interface{}) {
	panic("implement me")
}

func (m mockWechatClient) CodeToSession(code string) (*user.Code2SessionRes, error) {
	return &user.Code2SessionRes{
		SessionKey: "SJV2hpDNbb2O6t3LrS7Jkg==",
		Openid:     "oyL7e5ekyKGOWapDiHxmk17vhJH8",
	}, nil
}

func (m mockWechatClient) MiniProgramCode2Session(code string, encryptPhone string, iv string, data *user.MiniProgramResponseDTO) error {
	if data == nil {
		panic("data is nil")
	}

	//appID := "wx4053ee583576b3c7"
	//sessionKey := "SJV2hpDNbb2O6t3LrS7Jkg=="
	//encryptedData := "6cG2ImIHml5WC8UtTrxjCB3IZNdrzPkxYMz6O69UGfFQpzybnD3mQUJS5CEbVCk6G7YiekHUkzKAyxt03BubLnLk+lYizTw/UlMCr7o8Wil4gVjGY1uFbPdF1HaroHxW1pir2QuLeQQoih/uAy5v0Nf43dDTNwFomlnJcBA+Tq6c9Y8wf/xxG8iDSX5wrOwS23ZVYWgV2zcwq2nLjfnJkQ=="
	//iv := "88YK8pkJBTUOjY0Iy8fvwQ=="
	data.ErrorCode = 0
	data.SessionKey = "SJV2hpDNbb2O6t3LrS7Jkg=="
	data.OpenID = "oyL7e5ekyKGOWapDiHxmk17vhJH8"
	//phone, err := m.crypto.Decrypt(encryptPhone, iv, data.SessionKey)
	return nil
}

func (m mockWechatClient) GetAccessToken(code string, data *user.OAuthAccessTokenDTO) error {
	panic("implement me")
}

func (m mockWechatClient) GetUserInfo(accessToken string, openID string, data *user.OAuthInfoDTO) error {
	panic("implement me")
}
