package caller

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/sms"
	log "github.com/sirupsen/logrus"
	"strings"
)

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

func (c manual120Caller) Call(helpInfo *entities.HelpInfo) error {
	log.Infof("Call 120 for helpInfo %d", helpInfo.ID)

	if helpInfo.Exercise {
		log.Infof("exercise. skip")
		return nil
	}

	conf := interfaces.GetConfig()
	phones := strings.Split(conf.Notifier.CallPhones, ",")
	for _, mobile := range phones {
		err := c.doCall(mobile, helpInfo.ID)
		if err != nil {
			log.Errorf("AidCall target phone %s error: %v", mobile, err)
			return err
		}

		log.Infof("AidCall target phone %s", mobile)
	}

	return nil
}

func (c *manual120Caller) doCall(mobile string, aid int64) error {
	log.Infof("doCall for mobile %s, aid=%d", mobile, aid)

	path, err := c.pathGenerator.GeneratePath(aid, mobile[7:])
	if err != nil {
		return err
	}

	err = sms.SendSms(mobile, "SMS_234030627", map[string]string{"env": c.envPrefix, "path": path})
	return err
}
