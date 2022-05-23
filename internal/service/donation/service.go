package donation

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	service2 "aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/domain/emitter"
	page "aed-api-server/internal/pkg/query"
	"aed-api-server/internal/pkg/utils"
	"errors"
	"fmt"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"strconv"
	"time"

	"github.com/go-xorm/xorm"
)

const (
	TableNameDonation = "donation"
	TableNameRecord   = "donation_record"
)

type Service struct {
}

func init() {
	interfaces.S.Donation = NewService()
}

func NewService() service2.DonationService {
	return &Service{}
}

func (s *Service) CreateDonation(donation *entities.Donation) error {
	donation.CreatedAt = time.Now()
	_, err := db.Table(TableNameDonation).Insert(donation)
	return err
}

func (s *Service) UpdateDonation(donation *entities.Donation) error {
	_, err := db.Table(TableNameDonation).ID(donation.Id).Update(donation)
	return err
}

func (s *Service) UpdateDonationWithSession(
	session *xorm.Session,
	donation *entities.Donation,
) error {
	_, err := session.Table(TableNameDonation).ID(donation.Id).Update(donation)
	return err
}

func (s *Service) DeleteDonation(id int64) error {
	_, err := db.Table(TableNameDonation).Delete(&entities.Donation{Id: id})
	return err
}

func (s *Service) ListDonatorsDonation(userId int64) ([]*entities.DonationWithUserDonated, error) {
	resulsts, err := utils.PromiseAll(func() (interface{}, error) {
		return s.ListDonationByUserId(userId)
	}, func() (interface{}, error) {
		return s.ListUserPointsForDonations(userId)
	})

	if err != nil {
		return nil, err
	}

	donations := resulsts[0].([]*entities.Donation)
	userPonints := resulsts[1].(map[int64]int)
	donationsWithUserPoints := make([]*entities.DonationWithUserDonated, len(donations))

	for i := range donations {
		donation := donations[i]
		points := userPonints[donation.Id]
		donationsWithUserPoints[i] = &entities.DonationWithUserDonated{
			Donation:      *donation,
			DonatedPoints: points,
		}
	}

	return donationsWithUserPoints, nil
}

func (s *Service) ListDonationByUserId(userId int64) ([]*entities.Donation, error) {
	var donations []*entities.Donation
	sql := fmt.Sprintf(`
            select distinct(d.id) id,
			d.title title,
			d.images images,
			d.description description,
			d.target_points target_points,
			d.actual_points actual_points,
			d.start_at start_at,
			d.complete_at complete_at,
			d.expired_at expired_at,
			d.executor executor,
			d.executor_number executor_number,
			d.feedback feedback,
			d.plan plan,
			d.plan_image plan_image,
			d.budget budget
			from %s r left join %s d
     		on r.donation_id = d.id
			where r.user_id = ?	
           `, TableNameRecord, TableNameDonation)

	err := db.SQL(sql, userId).Find(&donations)
	fixStatuses(donations)
	return donations, err
}

func (s *Service) ListUserPointsForDonations(userId int64) (map[int64]int, error) {
	donationPointsMap := make(map[int64]int, 0)

	var records []*entities.DonationRecord

	sql := fmt.Sprintf(`
        select donation_id , sum(points) points, user_id from %s 
		where user_id = ?
        group by donation_id
    `, TableNameRecord)

	err := db.SQL(sql, userId).Find(&records)

	for _, e := range records {
		donationPointsMap[e.DonationId] = e.Points
	}

	return donationPointsMap, err
}

func (s *Service) ListDonation(p page.Query) ([]*entities.Donation, error) {
	var donations []*entities.Donation
	t := db.Table(TableNameDonation)
	if p.Page > 0 {
		t.Limit(p.Size, (p.Page-1)*p.Size)
	}

	err := t.Asc("sort").Find(&donations)
	fixStatuses(donations)
	return donations, err
}

func (s *Service) GetDonationDetail(id int64) (*entities.Donation, bool, error) {
	donation, exists, err := s.GetDonationById(id)
	if err != nil {
		return nil, exists, err
	}

	if !exists {
		return nil, false, nil
	}

	recordCount, err := s.GetDonationRecordCount(id)
	if err != nil {
		return nil, exists, err
	}

	i := int(recordCount)
	donation.RecordsCount = &i
	return donation, true, nil

}

func fixStatus(donation *entities.Donation) {
	if donation.StartAt.Before(time.Now()) {
		donation.Status = entities.StatusIng
	} else {
		donation.Status = entities.StatusNotStarted
	}

	if donation.CompleteAt != nil {
		donation.Status = entities.StatusCompleted
		return
	}

	if donation.ExpiredAt.Before(time.Now()) {
		donation.Status = entities.StatusExpired
	}
}

func fixStatuses(donations []*entities.Donation) {
	for _, d := range donations {
		fixStatus(d)
	}
}

func (s *Service) GetDonationByIdForUpdate(
	session *xorm.Session,
	id int64,
) (*entities.Donation, bool, error) {
	var one entities.Donation
	exists, err := session.Table(TableNameDonation).Where("id = ?", id).ForUpdate().Get(&one)
	fixStatus(&one)
	return &one, exists, err
}

func (s *Service) GetRecordById(recordId int64) (*entities.DonationRecord, bool, error) {
	var one entities.DonationRecord
	exists, err := db.Table(TableNameRecord).Where("id = ?", recordId).Get(&one)
	return &one, exists, err
}

func (s *Service) GetDonationById(id int64) (*entities.Donation, bool, error) {
	var one entities.Donation
	exists, err := db.Table(TableNameDonation).Where("id = ?", id).Get(&one)
	fixStatus(&one)
	return &one, exists, err
}

func (s *Service) GetDonationRecordCount(id int64) (int64, error) {
	count, err := db.Table(TableNameRecord).Where("donation_id = ?", id).Count()
	return count, err
}

func (s *Service) Donate(record *entities.DonationRecord) (*entities.Donation, error) {
	err := db.Begin(func(session *xorm.Session) error {
		donation, exists, err := s.GetDonationByIdForUpdate(session, record.DonationId)
		if err != nil {
			return err
		}
		if !exists {
			return errors.New("donation not found")
		}

		donatedPoints, err := donation.Donate(record.Points)
		if err != nil {
			return err
		}
		record.Points = donatedPoints

		//_, err = utils.PromiseAll(func() (interface{}, error) {
		err = s.UpdateDonationWithSession(session, donation)
		if err != nil {
			return err
		}
		//}, func() (interface{}, error) {
		err = s.CreateRecordWithSession(session, record)
		if err != nil {
			return err
		}
		err = interfaces.S.Medal.AwardMedalFirstDonation(record.UserId, record.Id)
		if err != nil {
			return err
		}
		err = interfaces.S.Points.AddPoint(
			record.UserId,
			-record.Points,
			fmt.Sprintf("参与项目：%s", donation.Title),
			entities.PointsEventTypeDonation,
		)
		if err != nil {
			return err
		}
		//}, func() (interface{}, error) {
		err = interfaces.S.Evidence.CreateTextEvidence(
			"项目捐献",
			record.UserId,
			entities.EvidenceCategoryDonation,
			strconv.FormatInt(record.Id, 10),
			map[string]interface{}{
				"donationId":        record.DonationId,
				"user_id":           record.UserId,
				"lastTransactionId": "", //TODO
			},
		)
		//})

		count, err := s.CountUserRecord(record.UserId)
		if err != nil {
			log.Error("CountUserRecord err", err)
			return err
		}
		if count < 1 {
			//触发积分奖励
			err = emitter.Emit(&events.PointsEvent{
				PointsEventType: entities.PointsEventTypeDonationAward,
				UserId:          record.UserId,
				Params: entities.PointsEventParams{
					RefTable:   "donation_record#",
					RefTableId: record.Id,
				},
			})
		}
		return err
	})

	if err != nil {
		return nil, err
	}

	donation, _, err := s.GetDonationDetail(record.DonationId)
	return donation, err
}

func (s *Service) CountUserRecord(userId int64) (int, error) {
	var count int64
	var err error
	err = db.Transaction(func(session *xorm.Session) error {
		count, err = session.SQL(`select 1 from donation_record where user_id = ?`, userId).Count()
		return err
	})
	return int(count), err
}

func (s *Service) CreateRecordWithSession(
	session *xorm.Session,
	record *entities.DonationRecord,
) error {
	_, err := session.Table(TableNameRecord).Insert(record)
	return err
}

func (s *Service) ListRecords(donationId int64, latest int) ([]*entities.DonationRecord, error) {
	var records []*entities.DonationRecord
	session := db.Table(TableNameRecord).Where("donation_id = ?", donationId).Desc("created_at")
	if latest > 0 {
		session.Limit(latest, 0)
	}

	err := session.Find(&records)
	return records, err
}

func (s *Service) ListUsersRecordsTop(
	donationId int64,
	top int,
) ([]*entities.DonationRecord, error) {
	var records []*entities.DonationRecord

	sql := fmt.Sprintf(`
		select donation_id, user_id, sum(points) as points
		from %s 
		where donation_id = ?
		group by user_id
		order by points desc
		limit ?
    `, TableNameRecord)
	err := db.SQL(sql, donationId, top).Find(&records)

	return records, err
}

func (s *Service) StatDonationByUserId(userId int64) (stat entities.DonationStat, err error) {
	_, err = db.SQL(`
		select
			sum(points) as donation_total_points,
			count(distinct donation_id) as donation_project_count,
			count(1) as donation_count
		from donation_record
		where
			user_id = ?
	`, userId).Get(&stat)

	return stat, err
}
