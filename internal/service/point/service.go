package point

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/response"
	"encoding/json"
	"github.com/go-xorm/xorm"
	"time"
)

type FlowStatus int8

const (
	pointFlowTableName = "point_flow"
	FlowStatusInit     = 0
	FlowStatusPecked   = 1
	FlowCut            = 1
)

type Point struct {
	Id          int64
	AccountId   int64
	Points      float64
	Description string
	CreatedAt   time.Time `xorm:"create_at"`
	RefTable    string    `xorm:"class"`
	RefId       string    `xorm:"extra"`
}

type Flow struct {
	Id              int64                    `json:"id"`
	PointsEventType entities.PointsEventType `json:"pointEventType"`

	Status FlowStatus `json:"status"`
	UserId int64      `json:"userId"`

	Params        string    `json:"params"`
	PeckExpiredAt time.Time `json:"peckExpiredAt"`
	PeckedAt      time.Time `json:"peckedAt"`
	Points        int       `json:"points"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"createdAt"`
}

func init() {
	interfaces.S.Points = service{}
}

type service struct{}

// GetUserPointsEventTimes 读取用户对某个行为的积分次数
func (s service) GetUserPointsEventTimes(userId int64, eventType entities.PointsEventType) (int64, error) {
	return db.Table("point_flow").Where("user_id = ? and points_event_type = ?", userId, eventType).Count()
}

func insertPoints(userId int64, points int, pointsEventType entities.PointsEventType, eventParam interface{}, expiredAt time.Time) error {
	return db.Transaction(func(session *xorm.Session) error {
		eventParamJson, err := json.Marshal(eventParam)
		if err != nil {
			return err
		}
		_, err = session.Table(pointFlowTableName).Insert(Flow{
			Status:          FlowStatusInit,
			UserId:          userId,
			PointsEventType: pointsEventType,
			Params:          string(eventParamJson),
			PeckExpiredAt:   expiredAt,
			Points:          points,
			CreatedAt:       time.Now(),
		})
		if err != nil {
			return err
		}
		return nil
	})
}

func (service) ReceivePoints(userId int64, pointId int64) error {
	return db.Transaction(func(session *xorm.Session) error {
		_, err := session.Exec(`
			update point_flow 
			set 
				status = 1, 
				pecked_at = now() 
			where 
				id = ? 
				and user_id = ? 
				and peck_expired_at >now()
				and status = 0
		`, pointId, userId)

		if err != nil {
			return err
		}
		return nil
	})
}

func (service) GetUnReceivePoints(userId int64) ([]*entities.UserPointsFlow, error) {
	flows := make([]*entities.UserPointsFlow, 0)

	err := db.SQL(`
		select
			a.id,
			b.name,
			a.points,
			a.peck_expired_at as expired_at
		from point_flow as a
		left join point_event_define as b
			on b.points_event_type = a.points_event_type
		where
			a.user_id = ?
			and a.status = 0
			and a.peck_expired_at > now()
		order by a.peck_expired_at asc
	`, userId).Find(&flows)
	if err != nil {
		return nil, err
	}
	return flows, nil
}

func (service) GetUnReceivePointsCount(userId int64) (int, error) {
	count, err := db.Table("point_flow").
		Where("user_id = ? and status = 0 and peck_expired_at > now()", userId).
		Count()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (service) GetUserPeriodIncomePointsRecords(userId int64, begin, end time.Time) ([]*entities.UserPointsRecord, error) {
	records := make([]*entities.UserPointsRecord, 0)

	err := db.Table("point_flow").Alias("a").
		Select("a.id, b.description, a.points, a.pecked_at as time").
		Join("LEFT", []string{"point_event_define", "b"},
			"b.points_event_type = a.points_event_type").
		Where("a.user_id = ?"+
			" and a.status = 1"+
			" and a.points > 0"+
			" and a.pecked_at >= ? and a.pecked_at < ?", userId, begin, end).
		Find(&records)

	return records, err
}

func (service) GetUserPointsRecords(userId int64) ([]*entities.UserPointsRecord, error) {
	records := make([]*entities.UserPointsRecord, 0)
	err := db.SQL(`
		select
			a.id,
			if(a.description = '',b.description,a.description) as description,
			a.points,
			a.pecked_at as time
		from point_flow as a
		left join point_event_define as b
			on b.points_event_type = a.points_event_type
		where
			a.user_id = ?
			and a.status = 1
		order by id desc
	`, userId).Find(&records)

	if err != nil {
		return nil, err
	}
	return records, nil
}

func (service) GetUserTotalPointsForUpdate(session *xorm.Session, userId int64) (int, error) {
	var (
		err   error
		count int64
	)
	count, err = session.SQL(`
			select 
				ifnull(sum(points),0) as count 
			from point_flow
			where
				user_id = ? 
				and status = 1
		`, userId).ForUpdate().Count()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (s service) GetUserIncomePoints(userId int64) (int, error) {
	var (
		err   error
		count int64
	)
	err = db.Transaction(func(session *xorm.Session) error {
		count, err = session.SQL(`
			select 
				ifnull(sum(points),0) as count 
			from point_flow
			where
				user_id = ? 
				and status = 1
				and points > 0
		`, userId).Count()
		return err
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (service) GetUserTotalPoints(userId int64) (int, error) {
	var (
		err   error
		count int64
	)
	err = db.Transaction(func(session *xorm.Session) error {
		count, err = session.SQL(`
			select 
				ifnull(sum(points),0) as count 
			from point_flow
			where
				user_id = ? 
				and status = 1
		`, userId).Count()
		return err
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (service) GetUserDonatedPoints(userId int64) (int, error) {
	var (
		err   error
		count int64
	)
	err = db.Transaction(func(session *xorm.Session) error {
		count, err = session.SQL(`
			select 
				abs(ifnull(sum(points),0)) as count 
			from point_flow
			where
				user_id = ? 
				and status = 1
				and points_event_type = ?
		`, userId, entities.PointsEventTypeDonation).Count()
		return err
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func getUserTotalUnReceivedPoints(userId int64) (int, error) {
	var (
		err   error
		count int64
	)
	err = db.Transaction(func(session *xorm.Session) error {
		count, err = session.SQL(`
			select 
				ifnull(sum(points),0) as count
			from point_flow
			where 
				user_id = ? 
				and status = 0 
				and peck_expired_at > 0
		`, userId).Count()
		return err
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (service) GetUsersPeriodIncomePoints(userIds []int64, beginTime time.Time, endTime time.Time) (map[int64]int, error) {
	type Rst struct {
		UserId int64 `xorm:"user_id"`
		Total  int   `xorm:"total"`
	}

	rsts := make([]*Rst, 0)

	err := db.Table("point_flow").
		Select(`user_id,sum(points) as total`).
		In("user_id", userIds).
		And("status=1 and pecked_at >= ? and pecked_at <= ? and points > 0", beginTime, endTime).
		GroupBy("user_id").
		Find(&rsts)
	if err != nil {
		return nil, err
	}

	m := make(map[int64]int)

	for i := range rsts {
		rst := rsts[i]
		m[rst.UserId] = rst.Total
	}
	return m, nil
}

func (s service) AddPoint(userId int64, points int, description string, eventType entities.PointsEventType) error {
	return db.Transaction(func(session *xorm.Session) error {
		total, err := s.GetUserTotalPointsForUpdate(session, userId)
		if err != nil {
			return err
		}

		if total+points < 0 {
			return response.ErrorInsufficientBalance
		}

		_, err = session.Table(pointFlowTableName).Insert(Flow{
			Status:          FlowCut,
			UserId:          userId,
			PointsEventType: eventType,
			Params:          "{}",
			PeckExpiredAt:   time.Now(),
			PeckedAt:        time.Now(),
			Points:          points,
			Description:     description,
			CreatedAt:       time.Now(),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
