package events

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
	"time"
)

type UserBaseEvent interface {
	emitter.DomainEvent

	//GetUserId 获取用户ID
	GetUserId() int64
}

type UserEvent struct {
	Id          int64 `xorm:"id pk autoincr"`
	UserId      int64
	EventType   entities.UserEventType
	EventParams []interface{}
	CreatedAt   time.Time
}

func (e *UserEvent) Decode(bt []byte) (emitter.DomainEvent, error) {
	var evt UserEvent
	err := json.Unmarshal(bt, &evt)
	if err != nil {
		return nil, err
	}
	return &evt, nil
}

func (e *UserEvent) Encode() ([]byte, error) {
	bt, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return bt, nil
}

func (e *UserEvent) GetUserId() int64 {
	return e.UserId
}
