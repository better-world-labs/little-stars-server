package speech

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/sms"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
)

//var targetAidCallPhones = []string{"15548720906"}

var targetAidCallPhones = []string{"18516597958", "13693520907", "15548720906"}

type AidCenterCaller interface {
	Call(*entities.HelpInfo) error
}

type manual120Caller struct {
	envPrefix     string
	pathGenerator PathGenerator
}

func NewAidCaller(envPrefix string, generator PathGenerator) AidCenterCaller {
	return &manual120Caller{
		envPrefix:     envPrefix,
		pathGenerator: generator,
	}
}

//func isPublisherInBlackList(mobile string) bool {
//	_, exists := triggerBlackList[mobile]
//	return exists
//}

func (c manual120Caller) Call(helpInfo *entities.HelpInfo) error {
	log.DefaultLogger().Infof("Call 120 for helpInfo %d", helpInfo.ID)

	//publisher, err := userService.GetAccountByID(r.Record.Publisher)
	//if err != nil {
	//	return err
	//}

	//if !isPublisherInBlackList(publisher.Mobile) {
	for _, mobile := range targetAidCallPhones {
		err := c.doCall(mobile, helpInfo.ID)
		if err != nil {
			log.DefaultLogger().Errorf("AidCall target phone %s error: %v", mobile, err)
			return err
		}

		log.DefaultLogger().Infof("AidCall target phone %s", mobile)
	}
	//}

	return nil
}

func (c *manual120Caller) doCall(mobile string, aid int64) error {
	log.DefaultLogger().Infof("doCall for mobile %s, aid=%d", mobile, aid)

	path, err := c.pathGenerator.GeneratePath(aid, mobile[7:])
	if err != nil {
		return err
	}

	err = sms.SendSms(mobile, "SMS_234030627", map[string]string{"env": c.envPrefix, "path": path})
	//fmt.Println(path)
	return err
}
