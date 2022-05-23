package events

import (
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// GoingToGetDeviceEvent 已去取设备
type GoingToGetDeviceEvent struct {
	Id     string
	Time   time.Time
	Aid    int64
	UserId int64
}

func NewGoingToGetDeviceEvent(aidId, userId int64) *GoingToGetDeviceEvent {
	return &GoingToGetDeviceEvent{
		Id:     uuid.NewString(),
		Time:   time.Now(),
		Aid:    aidId,
		UserId: userId,
	}
}

func (*GoingToGetDeviceEvent) Decode(bytes []byte) (emitter.DomainEvent, error) {
	var e GoingToGetDeviceEvent
	err := json.Unmarshal(bytes, &e)
	return &e, err
}

func (h *GoingToGetDeviceEvent) Encode() ([]byte, error) {
	return json.Marshal(h)
}
