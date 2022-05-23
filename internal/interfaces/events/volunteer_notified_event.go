package events

import (
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// VolunteerNotifiedEvent 已通知志愿者
type VolunteerNotifiedEvent struct {
	Id    string
	Count int
	Aid   int64
	Time  time.Time
}

func NewVolunteerNotifiedEvent(aidId int64, count int) *VolunteerNotifiedEvent {
	return &VolunteerNotifiedEvent{
		Id:    uuid.NewString(),
		Time:  time.Now(),
		Aid:   aidId,
		Count: count,
	}
}

func (*VolunteerNotifiedEvent) Decode(bytes []byte) (emitter.DomainEvent, error) {
	var e VolunteerNotifiedEvent
	err := json.Unmarshal(bytes, &e)
	return &e, err
}

func (h *VolunteerNotifiedEvent) Encode() ([]byte, error) {
	return json.Marshal(h)
}
