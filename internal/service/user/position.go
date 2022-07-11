package user

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/geo"
	"aed-api-server/internal/pkg/location"
	"github.com/go-xorm/xorm"
	"time"
)

const HeatStatPeriod = 30

type PositionHeat struct {
	UserId        int64
	HashId        uint64
	Heat          uint
	HeatUpdatedAt time.Time
}

type PositionRecord struct {
	UserId    int64
	HashId    uint64
	CreatedAt time.Time
	Lon       float64
	Lat       float64
}

type position interface {
	RecordPosition(userId int64, Lat, Lon float64) error
	UpdatePositionHeat(userId int64, hashId uint64) error
	UpdateAllPositionHeat() error
}

type positionService struct{}

func (p *positionService) RecordPosition(userId int64, Lat, Lon float64) error {
	hashId := geo.Hash(Lat, Lon)
	return db.Transaction(func(db *xorm.Session) error {
		//1. 记录位置
		err := writeUserPosition(userId, hashId, Lat, Lon)
		if err != nil {
			return err
		}

		//2. 更新位置
		err = updateUserPosition(userId, Lat, Lon)
		if err != nil {
			return err
		}

		//3. 更新热点
		return p.UpdatePositionHeat(userId, hashId)
	})
}

func writeUserPosition(userId int64, hashId uint64, lat float64, lon float64) error {
	record := PositionRecord{
		UserId:    userId,
		HashId:    hashId,
		Lat:       lat,
		Lon:       lon,
		CreatedAt: time.Now(),
	}

	_, err := db.Table("user_position_record").Insert(record)
	return err
}

func updateUserPosition(userId int64, lat float64, lon float64) error {
	userPosition := entities.Position{
		Coordinate: &location.Coordinate{
			Latitude:  lat,
			Longitude: lon,
		},
		AccountID: userId,
	}
	var current entities.Position
	exists, err := db.Table("account_position").
		Where("account_id = ?", userId).
		Get(&current)
	if err != nil {
		return err
	}

	if !exists {
		_, err := db.Table("account_position").Insert(userPosition)
		return err
	} else {
		_, err := db.Table("account_position").ID(current.ID).Update(userPosition)
		return err
	}
}

func (p *positionService) UpdatePositionHeat(userId int64, hashId uint64) error {
	return db.Transaction(func(db *xorm.Session) error {
		var heat PositionHeat
		_, err := db.SQL(`
			select
				user_id,
				hash_id,
				count(1) as heat,
				now() as heat_updated_at
			from user_position_record
			WHERE
				created_at > TIMESTAMPADD(day,?,CURRENT_DATE)
				and user_id = ?
				and hash_id = ?
    	`, -HeatStatPeriod, userId, hashId).Get(&heat)

		if err != nil {
			return err
		}

		_, err = db.Exec(`
			insert into user_position_heat(user_id, hash_id, heat, heat_updated_at)
			values(?,?,?,now())
			on duplicate key update
			heat = ?,
			heat_updated_at = now()
		`, userId, hashId, heat.Heat, heat.Heat)

		return err
	})
}

func (p *positionService) UpdateAllPositionHeat() error {
	return db.Transaction(func(db *xorm.Session) error {
		_, err := db.SQL(`delete * from user_position_heat`).Exec()
		if err != nil {
			return err
		}

		_, err = db.SQL(`
			insert into user_position_heat (user_id, hash_id, heat, heat_updated_at)
			select
				user_id,
				hash_id,
				count(1) as heat,
				now() as heat_updated_at
			from user_position_record
			WHERE
				created_at > TIMESTAMPADD(day,?,CURRENT_DATE)
			group by user_id, hash_id
		`, -HeatStatPeriod).Exec()
		return err
	})
}
