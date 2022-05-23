package events

import (
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// AidCalledEvent 已拨打急救电话
type AidCalledEvent struct {
	Id   string
	Time time.Time
	Aid  int64
}

func NewAidCalled(aidId int64) *AidCalledEvent {
	return &AidCalledEvent{
		Id:   uuid.NewString(),
		Time: time.Now(),
		Aid:  aidId,
	}
}

func (*AidCalledEvent) Decode(bytes []byte) (emitter.DomainEvent, error) {
	var e AidCalledEvent
	err := json.Unmarshal(bytes, &e)
	return &e, err
}

func (h *AidCalledEvent) Encode() ([]byte, error) {
	return json.Marshal(h)
}
