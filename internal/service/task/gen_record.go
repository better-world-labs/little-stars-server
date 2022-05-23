package task

import (
	"aed-api-server/internal/pkg/db"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"time"
)

const GenRecordTable = "task_gen_record"

type GenRecord struct {
	Id        int64     `json:"id" xorm:"id"`
	UserId    int64     `json:"userId" xorm:"user_id"`
	Latitude  string    `json:"latitude" xorm:"latitude"`
	Longitude string    `json:"longitude" xorm:"longitude"`
	CreatedAt time.Time `json:"createdAt" xorm:"created_at"`
}

func createRecord(record *GenRecord) (int64, error) {
	record.CreatedAt = time.Now()
	return db.GetSession().Table(GenRecordTable).Insert(record)
}

func findLastedByUserId(userId int64, after time.Time) (*GenRecord, error) {
	record := GenRecord{}
	has, err := db.SQL(`select * from task_gen_record where user_id = ? and created_at > ?`, userId, after).Get(&record)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	log.Info("record:", record)
	return &record, nil
}
