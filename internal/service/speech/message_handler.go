package speech

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/facility"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/base"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/pkg/sms"
	"aed-api-server/internal/pkg/star"
	"aed-api-server/internal/service/speech/caller"
	"aed-api-server/internal/service/speech/finder"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

const AidMiniProgramPagePath = "/subcontract/rescue/details"

func userService() service.UserServiceOld {
	return interfaces.S.UserOld
}

type AidProcessor struct {
	caller caller.AidCenterCaller
	finder finder.UserFinder

	Wechat service.IWechat `inject:"-"`
}

//go:inject-component
func NewAidProcessor() *AidProcessor {
	conf := interfaces.GetConfig()
	prefix := star.GetDomainPrefix(star.Env(conf.Server.Env))

	return &AidProcessor{
		caller: caller.NewAidCaller(prefix, caller.NewPathGenerator(caller.NewTokenService())),
		finder: finder.NewUserFinder(conf.Notifier.UserFinder),
	}
}

func (p *AidProcessor) Listen(on facility.OnEvent) {
	on(&events.HelpInfoPublishedEvent{}, p.onAidPublished)
}

func (p *AidProcessor) onAidPublished(event emitter.DomainEvent) error {
	re, ok := event.(*events.HelpInfoPublishedEvent)

	if !ok {
		return errors.New("invalid event type")
	}

	log.Infof("AidProcessor: onAidPublished, id=%d", re.ID)

	err := p.caller.Call(&re.HelpInfo)
	if err != nil {
		log.Errorf("Notifier: onAidPublished error: %v", err)
		return err
	}

	err = p.notify(&re.HelpInfo)
	if err != nil {
		log.Errorf("Notifier: onAidPublished error: %v", err)
	} else {
		log.Info("Notifier: onAidPublished succeed")
	}

	return nil
}

func (p *AidProcessor) notify(helpInfo *entities.HelpInfo) error {
	publisher, err := userService().GetUserByID(helpInfo.Publisher)
	if err != nil {
		return base.WrapError("AidProcessor", "find users to notify error", err)
	}

	if publisher == nil {
		return nil
	}

	accounts, err := p.finder.FindUser(helpInfo.Coordinate)
	if err != nil {
		return base.WrapError("AidProcessor", "find users to notify error", err)
	}

	var count int
	for _, account := range accounts {
		if account.ID == publisher.ID {
			continue
		}

		if !helpInfo.Exercise {
			err = p.doNotify(helpInfo, *account, publisher.Nickname, helpInfo.Address)
			if err != nil {
				log.Errorf("do notify for mobile %s error: %v", account.Mobile, err)
				continue
			}
		}

		log.Infof("notify user %s succeed", account.Nickname)
		count++
	}

	return emitter.Emit(events.NewVolunteerNotifiedEvent(helpInfo.ID, count))
}

func (p *AidProcessor) genericUrlLinkCode(path, query string) (string, error) {
	link, err := p.Wechat.GenericUrlLink(path, query)
	if err != nil {
		return "", err
	}

	return CreateLink(link)
}

func (p *AidProcessor) doNotify(info *entities.HelpInfo, target entities.User, publisherName string, detailAddress string) error {
	linkCode, err := p.genericUrlLinkCode(AidMiniProgramPagePath, fmt.Sprintf("id=%d", info.ID))
	if err != nil {
		log.Errorf("genericUrlLinkCode error: %v", err)
		return nil
	}

	env := star.Env(interfaces.GetConfig().Server.Env)
	titleSuffix := star.GetEnvText(env)
	if info.Practice {
		titleSuffix = fmt.Sprintf("演习%s", titleSuffix)
	}

	return sms.SendSms(target.Mobile, "SMS_250390117", map[string]string{
		"envText":   titleSuffix,
		"address":   detailAddress,
		"envDomain": star.GetDomainPrefix(env),
		"linkCode":  linkCode,
	})
}
