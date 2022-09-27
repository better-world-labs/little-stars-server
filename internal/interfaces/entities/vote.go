package entities

import (
	"time"
)

const (
	VoteProjectStatusPending = 0
	VoteProjectStatusStarted = 1
	VoteProjectStatusStop    = 2
)

const (
	VoteOptionTypeSingle = 0
	VoteOptionTypeMulti  = 1
)

type (
	Vote struct {
		Id         int64
		Name       string
		Image      string
		Text       string
		MaxTimes   int
		OptionType int
		CreatedAt  time.Time
		BeginAt    time.Time
		EndAt      time.Time
	}

	VoteOptionDetail struct {
		VoteOption

		Rank int `json:"rank"`
	}

	VoteOption struct {
		Id     int64  `json:"id"`
		VoteId int64  `json:"voteId"`
		Text   string `json:"text"`
		Vote   int    `json:"vote"`
	}

	VoteRecord struct {
		VoteId    int64
		UserId    int64
		Mode      int
		OptionIds []int64
	}
)

const (
	VoteRecordModeNormal = 1 //正常投票模式
	VoteRecordModePoints = 2 //积分投票模式
)

func (v Vote) Status() int {
	now := time.Now()
	if now.Before(v.BeginAt) {
		return VoteProjectStatusPending
	}

	if now.After(v.BeginAt) && now.Before(v.EndAt) {
		return VoteProjectStatusStarted
	}

	return VoteProjectStatusStop
}
