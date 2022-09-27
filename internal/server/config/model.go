package config

import (
	"aed-api-server/internal/pkg/cache"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/location"
	"aed-api-server/internal/pkg/sms"
	"aed-api-server/internal/service/oss"
	openapi "gitlab.openviewtech.com/openview-pub/gopkg/open-api"

	"gitlab.openviewtech.com/openview-pub/gopkg/log"
)

type (
	// AppConfig  应用配置结构体
	// 禁止依赖模块本身
	AppConfig struct {
		Server ServerConfig `properties:"server"`
		//Domain                  config.DomainEventConfig `properties:"domain-event"`
		Log                     log.LogConfig     `properties:"log"`
		Database                db.MysqlConfig    `properties:"database"`
		Wechat                  WechatOAuthConfig `properties:"wechat"`
		AliOss                  oss.Config        `properties:"alioss"`
		Redis                   cache.RedisConfig `properties:"redis"`
		SmsClient               sms.Config        `properties:"aliyun-sms"`
		JwtConfig               JwtConfig         `properties:"jwt"`
		CredentialConfig        openapi.Config    `properties:"credential"`
		EvidenceConfig          openapi.Config    `properties:"evidence"`
		Notifier                NotifierConfig    `properties:"notifier"`
		Exam                    Exam              `properties:"exam"`
		Backend                 Backend           `properties:"backend"`
		WechatToken             Service           `properties:"wechat-token"`
		Tencent                 location.Config   `properties:"tencent_map"`
		DonationApplyNotify     string            `properties:"donation-apply-notify"`
		CptAedCert              int               `properties:"cpt-aed-cert"` //证书在凭证服务中CPT的编号
		CptMedal                int               `properties:"cpt-medal"`    //勋章在凭证服务中CPT的编号
		ImgBotService           string            `properties:"img-bot-service"`
		WechatOffiaccountAppid  string            `properties:"wechat-offiaccount-appid"`
		WechatOffiaccountSecret string            `properties:"wechat-offiaccount-secret"`
		ClockInRangeCheck       bool              `properties:"clock-in-range-check"`
		PracticeOrganizer       int64             `properties:"practice-organizer"`
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

	NotifierConfig struct {
		UserFinder string `properties:"user-finder"`
		CallPhones string `properties:"call-phones"`
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
