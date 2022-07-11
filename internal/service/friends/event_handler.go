package friends

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/facility"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/domain/emitter"
	"errors"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
)

func (f *Service) Listen(on facility.OnEvent) {
	on(&events.FirstLoginEvent{}, handleNewUserLogin)
	on(&entities.Trace{}, handleNewTraceCreated)
}

func handleNewUserLogin(event emitter.DomainEvent) error {
	_tag := "[friends.handleNewUserLogin]"

	if evt, ok := event.(*events.FirstLoginEvent); ok {
		log.Info(_tag, "handleNewUserLogin: userId=", evt.UserId)
		prospective, exists, err := getProspective(evt.Openid)
		if err != nil {
			return err
		}

		if !exists {
			log.Info(_tag, "no sharer found, skip")
			return nil
		}

		return db.Begin(func(session *xorm.Session) error {
			err := createRelationship(&Relationship{
				UserId:   evt.UserId,
				ParentId: prospective.ParentId,
				CreateAt: evt.LoginAt,
			})
			if err != nil {
				return err
			}

			err = emitter.Emit(&events.PointsEvent{
				PointsEventType: entities.PointsEventTypeBeInvited,
				UserId:          evt.UserId,
				Params: map[string]interface{}{
					"userId": evt.UserId,
				},
			})

			return err
		})
	}

	return errors.New("invalid event type")
}

func handleNewTraceCreated(event emitter.DomainEvent) error {
	_tag := "[friends.handleNewTraceCreated]"

	if trace, ok := event.(*entities.Trace); ok {
		log.Info(_tag, "sharer=%s, openid=%s", trace.Sharer, trace.OpenID)
		_, isOldUser, err := interfaces.S.User.GetUserByOpenid(trace.OpenID)
		if err != nil {
			return err
		}

		if isOldUser {
			log.Info(_tag, "old user, skip")
			return nil
		}

		tag, err := trace.GetSharerTag()
		if err != nil {
			// 数据有问题，跳过
			return nil
		}

		var parent *entities.SimpleUser
		var exists bool

		switch tag.(type) {
		case entities.OpenIdTag:
			tag := tag.(entities.OpenIdTag)
			parent, exists, err = interfaces.S.User.GetUserByOpenid(string(tag))
			if err != nil {
				return err
			}

		case entities.UserIdTag:
			tag := tag.(entities.UserIdTag)
			parent, exists, err = interfaces.S.User.GetUserById(int64(tag))
			if err != nil {
				return err
			}

		default:
			log.Info("skip tag type")
		}

		if exists {
			err := createProspective(&Prospective{
				OpenId:   trace.OpenID,
				ParentId: parent.ID,
			})
			return err
		}
	}

	return errors.New("invalid event type")
}
