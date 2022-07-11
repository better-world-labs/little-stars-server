package speech

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/sms"
	log "github.com/sirupsen/logrus"
)

const targetAidCallPhone = "18816570787"

var triggerBlackList = map[string]interface{}{
	"18349162361": nil,
	"18512303122": nil,
	"18628138065": nil,
	"13880559010": nil,
}

type manual120CallerOld struct {
}

func NewAidCallerOld() AidCenterCaller {
	return &manual120CallerOld{}
}

func isPublisherInBlackList(mobile string) bool {
	_, exists := triggerBlackList[mobile]
	return exists
}

func (c manual120CallerOld) Call(helpInfo *entities.HelpInfo) error {
	log.Infof("Call 120 for helpInfo %d", helpInfo.ID)

	publisher, err := userService().GetUserByID(helpInfo.ID)
	if err != nil {
		return err
	}

	if !isPublisherInBlackList(publisher.Mobile) {
		return sms.SendSms(targetAidCallPhone, "SMS_232178406", nil)
	}

	log.Infof("publisher %s in black list, skip", publisher.Nickname)
	return nil
}
