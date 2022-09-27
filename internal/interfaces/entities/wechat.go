package entities

import (
	"aed-api-server/internal/pkg/utils"
	"time"
)

type (
	WechatDataDecryptReq struct {
		Code          string `json:"code"`
		Iv            string `json:"iv"`
		EncryptedData string `json:"encryptedData"`
	}

	WechatOAuthInfoRes struct {
		WechatRst

		UnionID    string `json:"unionid"`
		Nickname   string `json:"nickname"`
		HeadImgURL string `json:"headimgurl"`
	}

	WechatOAuthAccessTokenRes struct {
		WechatRst

		AccessToken string `json:"access_token"`
		OpenID      string `json:"openid"`
	}

	WechatCode2SessionRes struct {
		WechatRst

		SessionKey string `json:"session_key"`
		Openid     string `json:"openid"`
		UnionId    string `json:"unionid"`
	}

	WechatMiniProgramRes struct {
		WechatRst

		OpenID         string `json:"openid"`
		UnionID        string `json:"unionid"`
		SessionKey     string `json:"session_key"`
		DecryptedPhone string `json:"-"`
	}

	WechatRst struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}

	WechatStepInfo struct {
		TimeStamp int64 `json:"timeStamp"`
		Step      int   `json:"step"`
	}

	WechatWalkData struct {
		StepInfoList []*WechatStepInfo `json:"stepInfoList"`
	}

	WechatEncryptKey struct {
		EncryptKey string `json:"encrypt_key"`
		Version    int    `json:"version"`
		ExpireIn   int64  `json:"expire_in"`
		Iv         string `json:"iv"`
		CreateTime int64  `json:"create_time"`
	}
)

func (w WechatWalkData) GetTodayWalks() int {
	for _, info := range w.StepInfoList {
		t := time.UnixMilli(info.TimeStamp * 1000)
		if utils.IsToday(t) {
			return info.Step
		}
	}

	return 0
}

func (w WechatEncryptKey) ExpiredInSecond() time.Duration {
	return time.Duration(w.ExpireIn)
}

func (w WechatEncryptKey) Expired() bool {
	created := time.UnixMilli(w.CreateTime * 1000)
	return time.Now().After(created.Add(time.Duration(w.ExpireIn) * time.Second))
}

func GetEncryptKey(keys []*WechatEncryptKey, version int) *WechatEncryptKey {
	for _, k := range keys {
		if k.Version == version {
			return k
		}
	}

	return nil
}
