package speech

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/base"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/pkg/location"
	"errors"
	log "github.com/sirupsen/logrus"
)

func userService() service.UserServiceOld {
	return interfaces.S.UserOld
}

var n PhoneNotifier = AliyunPhoneNotifier{}

const Module = "Speech"

type UserFinder interface {
	FindUser(position location.Coordinate) ([]*entities.User, error)
}

func onAidPublished(event emitter.DomainEvent) error {
	re, ok := event.(*events.HelpInfoPublishedEvent)

	if !ok {
		return errors.New("invalid event type")
	}

	err := caller.Call(&re.HelpInfo)
	if err != nil {
		log.Errorf("Notifier: onAidPublished error: %v", err)
		return err
	}

	err = DoNotify(&re.HelpInfo)
	if err != nil {
		log.Errorf("Notifier: onAidPublished error: %v", err)
	} else {
		log.Info("Notifier: onAidPublished succeed")
	}

	return nil
}

func DoNotify(helpInfo *entities.HelpInfo) error {
	publisher, err := userService().GetUserByID(helpInfo.Publisher)
	if err != nil {
		return base.WrapError(Module, "find users to notify error", err)
	}

	if publisher == nil {
		return nil
	}

	accounts, err := finder.FindUser(*helpInfo.Coordinate)
	if err != nil {
		return base.WrapError(Module, "find users to notify error", err)
	}

	var count int
	for _, account := range accounts {
		if account.ID == publisher.ID {
			continue
		}

		//address := fmt.Sprintf("%s %s", r.Address, r.DetailAddress)
		err = n.Notify(*account, publisher.Nickname, helpInfo.Address)
		if err != nil {
			log.Errorf("do notify for mobile %s error: %v", account.Mobile, err)
			continue
		}

		log.Infof("notify user %s succeed", account.Nickname)
		count++
	}

	return emitter.Emit(events.NewVolunteerNotifiedEvent(helpInfo.ID, count))
}
