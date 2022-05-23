package events

import (
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// DeviceGotEvent 得到设备
type DeviceGotEvent struct {
	Id     string
	Time   time.Time
	Points int
	UserId int64
	Aid    int64
}

func NewDeviceGotEvent(aidId, userId int64) *DeviceGotEvent {
	return &DeviceGotEvent{
		Id:     uuid.NewString(),
		Time:   time.Now(),
		Aid:    aidId,
		UserId: userId,
	}
}

func (*DeviceGotEvent) Decode(bytes []byte) (emitter.DomainEvent, error) {
	var e DeviceGotEvent
	err := json.Unmarshal(bytes, &e)
	return &e, err
}

func (h *DeviceGotEvent) Encode() ([]byte, error) {
	return json.Marshal(h)
}
