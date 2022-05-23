package test

import (
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
	"time"
)

type TimeTick struct {
	Id        string        `json:"id"`
	Timestamp time.Duration `json:"timestamp"`
	Tick      int           `json:"tick"`
	Time      time.Time     `json:"time"`
}

func (t *TimeTick) Encode() ([]byte, error) {
	return json.Marshal(&t)
}

func (t *TimeTick) Decode(bytes []byte) (emitter.DomainEvent, error) {
	var evt TimeTick
	err := json.Unmarshal(bytes, &evt)
	return &evt, err
}
