package evidence

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/service/evidence/credential/claim"
	"errors"
	"github.com/sirupsen/logrus"
	"strconv"
)

var (
	ErrorInvalidEventType = errors.New("event assert failed, invalid event type")
)

func (s evidenceImpl) onMedalAwarded(e emitter.DomainEvent) error {
	logrus.Info("onMedalAwarded")

	if evt, ok := e.(*events.MedalAwarded); ok {
		medal, exists, err := s.Medal.GetById(evt.MedalID)
		if err != nil {
			return err
		}

		if !exists {
			return errors.New("medal not found")
		}

		user, err := s.AccountService.GetUserByID(evt.UserID)
		if err != nil {
			return err
		}

		if user == nil {
			return errors.New("user nod found")
		}

		return s.CreateEvidence(&claim.Medal{
			Mobile: user.Mobile,
			Medal:  medal.Name,
		}, medal.Name, user.ID, entities.EvidenceCategoryMedal, strconv.FormatInt(evt.ID, 10))
	}

	return ErrorInvalidEventType
}
