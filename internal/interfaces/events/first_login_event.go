package events

import (
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
	"time"
)

type FirstLoginEvent struct {
	UserId  int64     `json:"userId"`
	Openid  string    `json:"openId"`
	LoginAt time.Time `json:"loginAt"`
}

func (FirstLoginEvent) Decode(bytes []byte) (emitter.DomainEvent, error) {
	var evt FirstLoginEvent
	err := json.Unmarshal(bytes, &evt)
	if err != nil {
		return nil, err
	}
	return &evt, nil
}

func (f *FirstLoginEvent) Encode() ([]byte, error) {
	return json.Marshal(f)
}

func (f *FirstLoginEvent) GetUserId() int64 {
	return f.UserId
}
