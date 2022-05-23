package medal

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/module/achievement"
	"aed-api-server/internal/module/evidence/credential/claim"
	"aed-api-server/internal/pkg/base"
	"aed-api-server/internal/pkg/db"
	"errors"
	"github.com/go-xorm/xorm"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"strconv"
)

type Service struct {
}

func Init() {
	interfaces.S.Medal = &Service{}
}

func (m *Service) AwardMedalSaveLife(userId, helpInfoId int64) error {
	return m.AwardMedal(userId, achievement.MedalIdSaveLife, strconv.FormatInt(helpInfoId, 10))
}

func (m *Service) AwardMedalFirstDonation(userId, donationRecordId int64) error {
	_, exists, err := achievement.GetUserMedal(achievement.MedalIdFirstDonation, userId)
	if err != nil {
		return err
	}

	if !exists {
		return m.AwardMedal(userId, achievement.MedalIdFirstDonation, strconv.FormatInt(donationRecordId, 10))
	}

	return nil
}

func (m *Service) AwardMedal(userId, medalId int64, businessId string) error {
	medal, exists, err := achievement.GetById(medalId)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("medal not found")
	}

	_, exists, err = achievement.GetUserMedal(achievement.MedalIdInspector, userId)
	if exists {
		return nil
	}

	userMedal := &achievement.UserMedal{
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

		err = achievement.CreateUsersMedal(session, userMedal)
		if err != nil {
			log.DefaultLogger().Error("create user medal error: %v", err)
			return err
		}

		errChan := interfaces.S.Evidence.CreateEvidenceAsync(&claim.Medal{
			Mobile: account.Uid,
			Medal:  medal.Name,
		}, medal.Name, userId, entities.EvidenceCategoryMedal, strconv.FormatInt(userMedal.ID, 10))

		err = achievement.CreateUsersMedalToast(session, userMedal)
		if err != nil {
			log.DefaultLogger().Error("create user medal toast error: %v", err)
			return err
		}

		return <-errChan
	})
}

func (m *Service) AwardMedalInspector(userId int64) error {
	return m.AwardMedal(userId, achievement.MedalIdInspector, "device")
}
