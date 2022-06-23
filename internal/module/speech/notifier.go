package speech

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/sms"
	log "github.com/sirupsen/logrus"
)

type PhoneNotifier interface {
	Notify(target entities.User, publisherName string, detailAddress string) error
}

type AliyunPhoneNotifier struct{}

func (a AliyunPhoneNotifier) Notify(target entities.User, publisherName string, detailAddress string) error {
	t := make(map[string]string, 3)
	log.Info("Notify: user=", target.Nickname, "address: ", detailAddress)
	t["address"] = detailAddress
	return sms.SendSms(target.Mobile, "SMS_232163435", t)
}
