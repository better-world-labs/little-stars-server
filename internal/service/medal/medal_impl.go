package medal

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/facility"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/base"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/pkg/utils"
	"errors"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type medalImpl struct {
}

//go:inject-component
func NewMedalService() service.IMedal {
	s := &medalImpl{}
	return s
}

func (m *medalImpl) Listen(on facility.OnEvent) {
	on(&events.ClockInEvent{}, m.onDeviceClockIn)
	on(&events.DeviceMarkedEvent{}, m.onDeviceMarked)
}

func (m *medalImpl) AwardMedalSaveLife(userId, helpInfoId int64) error {
	return m.AwardMedal(userId, entities.MedalIdSaveLife, strconv.FormatInt(helpInfoId, 10))
}

func (m *medalImpl) AwardMedalFirstDonation(userId, donationRecordId int64) error {
	defer utils.TimeStat("AwardMedalFirstDonation")()
	_, exists, err := GetUserMedal(entities.MedalIdFirstDonation, userId)
	if err != nil {
		return err
	}

	if !exists {
		return m.AwardMedal(userId, entities.MedalIdFirstDonation, strconv.FormatInt(donationRecordId, 10))
	}

	return nil
}

func (m *medalImpl) AwardMedal(userId, medalId int64, businessId string) error {
	_, exists, err := m.GetById(medalId)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("medal not found")
	}

	_, exists, err = GetUserMedal(entities.MedalIdInspector, userId)
	if exists {
		return nil
	}

	userMedal := &entities.UserMedal{
		MedalID:    medalId,
		UserID:     userId,
		BusinessId: businessId,
	}

	return db.Begin(func(session *xorm.Session) error {
		account, exists, err := interfaces.S.User.GetUserById(userId)
		if err != nil {
			return err
		}

		if !exists {
			return errors.New("user not found")
		}

		if account == nil {
			return base.NewError("achievement.doIssue", "account not found")
		}

		err = CreateUsersMedal(session, userMedal)
		if err != nil {
			log.Errorf("create user medal error: %v", err)
			return err
		}

		err = CreateUsersMedalToast(session, userMedal)
		if err != nil {
			log.Errorf("create user medal toast error: %v", err)
			return err
		}

		return emitter.Emit(events.NewMedalAwarded(*userMedal))
	})
}

func (m *medalImpl) AwardMedalInspector(userId int64) error {
	return m.AwardMedal(userId, entities.MedalIdInspector, "device")
}

func (m *medalImpl) ListMedals() ([]*entities.Medal, error) {
	return ListMedals()
}

func (m *medalImpl) GetById(id int64) (*entities.Medal, bool, error) {
	return GetById(id)
}
