package merit_tree

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/pkg/global"
	"aed-api-server/internal/pkg/response"
	"github.com/go-xorm/xorm"
	"time"
)

type EarlyService struct{}

//go:inject-component
func NewEarlyService() service.EarlyService {
	return &EarlyService{}
}

func (e EarlyService) GetLatestRecord(userId int64) (*service.EarlyRecord, bool, error) {
	var early service.EarlyRecord
	exists, err := db.Table("sign_early").Where("user_id = ?", userId).Desc("time").Get(&early)
	if err != nil {
		return nil, exists, err
	}

	return &early, exists, nil
}

func (e EarlyService) SignEarly(userId int64) (*service.EarlyRecord, error) {
	r := service.EarlyRecord{
		UserId: userId,
		Time:   global.FormattedTime(time.Now()),
	}

	if !checkEarlyTimeOk(r.Time.Time()) {
		return nil, response.ErrorSignEarlyTimeNotAllowed
	}

	record, exists, err := e.GetLatestRecord(userId)
	if err != nil {
		return nil, err
	}

	if !exists {
		return e.Save(service.NewEarlyRecord(userId))
	}

	if checkSameDay(record.Time.Time(), r.Time.Time()) {
		return nil, response.ErrorSignEarlyTodayAlreadySignedYet
	}

	if checkSameDay(record.Time.Time(), r.Time.Time().Add(-time.Hour*24)) {
		r.Days = record.Days + 1
	} else {
		r.Days = 1
	}

	return e.Save(&r)
}

func (e EarlyService) Save(record *service.EarlyRecord) (r *service.EarlyRecord, err error) {
	err = db.Begin(func(session *xorm.Session) error {
		_, err := session.Table("sign_early").Insert(record)
		if err != nil {
			return err
		}

		r = record
		event := interfaces.S.PointsScheduler.BuildPointsEventTypeSignEarly(record.UserId, record.Id, record.Days)
		return emitter.Emit(event)
	})

	return
}

func checkEarlyTimeOk(t time.Time) bool {
	return t.Hour() >= 5 && t.Hour() <= 9
}

func checkSameDay(t1 time.Time, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 &&
		m1 == m2 &&
		d1 == d2
}
