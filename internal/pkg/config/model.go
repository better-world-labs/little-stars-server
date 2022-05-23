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
		Server            ServerConfig             `yaml:"server"`
		Domain            config.DomainEventConfig `yaml:"domain-event"`
		Log               log.LogConfig            `yaml:"log"`
		Database          db.MysqlConfig           `yaml:"database"`
		Wechat            WechatOAuthConfig        `yaml:"wechat"`
		AliOss            oss.Config               `yaml:"alioss"`
		Kafka             KafkaConfig              `yaml:"kafka"`
		MapConfig         tencent.Config           `yaml:"tencent_map"`
		Redis             cache.RedisConfig        `yaml:"redis"`
		DidServer         openapi.Config           `yaml:"did-server"`
		SmsClient         sms.Config               `yaml:"aliyun-sms"`
		JwtConfig         JwtConfig                `yaml:"jwt"`
		MiniProgramQrcode MiniProgramQrcodeConfig  `yaml:"mini-program-qrcode"`
		CredentialConfig  openapi.Config           `yaml:"credential"`
		EvidenceConfig    openapi.Config           `yaml:"evidence"`
		//FiscoBcos         openapi.Config           `yaml:"fisco-bcos"`
		Notifier    NotifierConfig           `yaml:"notifier"`
		Exam        Exam                     `yaml:"exam"`
		DomainEvent config.DomainEventConfig `yaml:"domain.emitter"`

		Env                 string `yaml:"env-name"`
		Host                string `yaml:"env-host"`
		DonationApplyNotify string `yaml:"donation-apply-notify"`

		CptAedCert int `yaml:"cpt-aed-cert"` //证书在凭证服务中CPT的编号
		CptMedal   int `yaml:"cpt-medal"`    //勋章在凭证服务中CPT的编号

		Backend       Backend `yaml:"backend"`
		ImgBotService string  `yaml:"img-bot-service"`

		WechatOffiaccountAppid  string `yaml:"wechat-offiaccount-appid"`
		WechatOffiaccountSecret string `yaml:"wechat-offiaccount-secret"`
	}

	// ServerConfig 服务配置结构体
	ServerConfig struct {
		Port        int    `yaml:"port"`
		Environment string `yaml:"environment"`
	}

	KafkaConfig struct {
		BootStrapServers string `yaml:"boot-strap-servers"`
		Topic            string `yaml:"topic"`
	}

	JwtConfig struct {
		ExpiresIn int64  `yaml:"expiresIn"`
		Secret    string `yaml:"secret"`
	}

	MiniProgramQrcodeConfig struct {
		ContentRootPath string `yaml:"content-root-path"`
	}

	NotifierConfig struct {
		UserFinder string `yaml:"user-finder"`
	}

	WechatOAuthConfig struct {
		Server    string `yaml:"server"`
		AppID     string `yaml:"app-key"`
		AppSecret string `yaml:"app-secret"`
	}
	Exam struct {
		Debug bool `yaml:"debug"`
	}

	Backend struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Id       int64  `yaml:"id"`
	}
)
