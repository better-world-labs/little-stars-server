package sms

import (
	"aed-api-server/internal/pkg/base"
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	log "github.com/sirupsen/logrus"
	"gitlab.openviewtech.com/openview-pub/gopkg/error2"
)

var client *dysmsapi.Client
var enabled bool

const Module = "AliyunSms"

func InitSmsClient(config Config) {
	enabled = config.Enabled
	c, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", config.AccessKeyID, config.AccessKeySecret)
	error2.MustNil(err)
	client = c
}

func SendSms(mobile string, templateCode string, data map[string]string) error {
	if !enabled {
		return nil
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = "开源观科技"
	request.TemplateCode = templateCode
	request.PhoneNumbers = mobile

	if data != nil {
		param, err := json.Marshal(data)
		if err != nil {
			return base.WrapError(Module, "send smd error", err)
		}

		request.TemplateParam = string(param)
	}

	response, err := client.SendSms(request)
	error2.MustNil(err)
	if response.Code != "OK" {
		log.Warnf("respone incorrect: %s, msg: %s", response.Code, response.Message)
		if response.Code == "isv.BUSINESS_LIMIT_CONTROL" {
			return base.NewError(Module, "request overflow")
		}

		return base.NewError(Module, response.Message)
	}

	return nil
}
