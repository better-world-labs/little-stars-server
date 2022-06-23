package events

import (
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
	"time"
)

type UserOpenTreasureChest struct {
	UserId            int64     `json:"userId"`
	TreasureChestName string    `json:"treasureChestName"`
	TaskId            int64     `json:"taskId"`
	Points            int       `json:"points"`
	OpenTime          time.Time `json:"openTime"`
	Link              string    `json:"link"`
	LinkArgs          []string  `json:"linkArgs"`
}

func (e *UserOpenTreasureChest) Decode(bt []byte) (emitter.DomainEvent, error) {
	var evt UserOpenTreasureChest
	err := json.Unmarshal(bt, &evt)
	if err != nil {
		return nil, err
	}
	return &evt, nil
}

func (e *UserOpenTreasureChest) Encode() ([]byte, error) {
	bt, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return bt, nil
}

func (e *UserOpenTreasureChest) GetUserId() int64 {
	return e.UserId
}
