package config

import (
	"aed-api-server/internal/module/oss"
	"aed-api-server/internal/pkg/cache"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/domain/config"
	"aed-api-server/internal/pkg/sms"
	"aed-api-server/internal/pkg/tencent"
	openapi "gitlab.openviewtech.com/openview-pub/gopkg/open-api"

	"gitlab.openviewtech.com/openview-pub/gopkg/log"
)

type (
	// AppConfig  应用配置结构体
	// 禁止依赖模块本身
	AppConfig struct {
		Server            ServerConfig             `properties:"server"`
		Domain            config.DomainEventConfig `properties:"domain-event"`
		Log               log.LogConfig            `properties:"log"`
		Database          db.MysqlConfig           `properties:"database"`
		Wechat            WechatOAuthConfig        `properties:"wechat"`
		AliOss            oss.Config               `properties:"alioss"`
		MapConfig         tencent.Config           `properties:"tencent_map"`
		Redis             cache.RedisConfig        `properties:"redis"`
		SmsClient         sms.Config               `properties:"aliyun-sms"`
		JwtConfig         JwtConfig                `properties:"jwt"`
		MiniProgramQrcode MiniProgramQrcodeConfig  `properties:"mini-program-qrcode"`
		CredentialConfig  openapi.Config           `properties:"credential"`
		EvidenceConfig    openapi.Config           `properties:"evidence"`
		Notifier          NotifierConfig           `properties:"notifier"`
		Exam              Exam                     `properties:"exam"`
		Backend           Backend                  `properties:"backend"`
		WechatToken       Service                  `properties:"wechat-token"`

		DonationApplyNotify     string `properties:"donation-apply-notify"`
		CptAedCert              int    `properties:"cpt-aed-cert"` //证书在凭证服务中CPT的编号
		CptMedal                int    `properties:"cpt-medal"`    //勋章在凭证服务中CPT的编号
		ImgBotService           string `properties:"img-bot-service"`
		WechatOffiaccountAppid  string `properties:"wechat-offiaccount-appid"`
		WechatOffiaccountSecret string `properties:"wechat-offiaccount-secret"`
	}

	// ServerConfig 服务配置结构体
	ServerConfig struct {
		Host string `properties:"host"`
		Port int    `properties:"port"`
		Env  string `properties:"env"`
	}

	Service struct {
		Host string `properties:"host"`
		Port int    `properties:"port"`
	}

	JwtConfig struct {
		ExpiresIn int64  `properties:"expiresIn"`
		Secret    string `properties:"secret"`
	}

	MiniProgramQrcodeConfig struct {
		ContentRootPath string `properties:"content-root-path"`
	}

	NotifierConfig struct {
		UserFinder string `properties:"user-finder"`
	}

	WechatOAuthConfig struct {
		Server    string `properties:"server"`
		AppID     string `properties:"app-key"`
		AppSecret string `properties:"app-secret"`
	}
	Exam struct {
		Debug bool `properties:"debug"`
	}

	Backend struct {
		Username string `properties:"username"`
		Password string `properties:"password"`
		Id       int64  `properties:"id"`
	}
)
