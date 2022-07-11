package service

import (
	"aed-api-server/internal/pkg/global"
	"time"
)

type EarlyRecord struct {
	Id     int64                `json:"id"`
	UserId int64                `json:"userId"`
	Days   int                  `json:"days"`
	Time   global.FormattedTime `json:"time"`
}

func NewEarlyRecord(userId int64) *EarlyRecord {
	return &EarlyRecord{
		Days:   1,
		UserId: userId,
		Time:   global.FormattedTime(time.Now()),
	}
}

type EarlyService interface {
	GetLatestRecord(userId int64) (*EarlyRecord, bool, error)
	SignEarly(userId int64) (*EarlyRecord, error)
}
