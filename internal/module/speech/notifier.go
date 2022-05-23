package speech

import (
	"aed-api-server/internal/module/user"
	"aed-api-server/internal/pkg/sms"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
)

type PhoneNotifier interface {
	Notify(target user.User, publisherName string, detailAddress string) error
}

type AliyunPhoneNotifier struct{}

func (a AliyunPhoneNotifier) Notify(target user.User, publisherName string, detailAddress string) error {
	t := make(map[string]string, 3)
	log.Info("Notify: user=", target.Nickname, "address: ", detailAddress)
	t["address"] = detailAddress
	return sms.SendSms(target.Mobile, "SMS_232163435", t)
}
