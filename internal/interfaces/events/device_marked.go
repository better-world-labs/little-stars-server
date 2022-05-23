package events

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
)

// DeviceMarkedEvent 设备标记
type DeviceMarkedEvent struct {
	entities.Device
}

func NewDeviceMarkedEvent(d entities.Device) *DeviceMarkedEvent {
	return &DeviceMarkedEvent{
		Device: d,
	}
}

func (*DeviceMarkedEvent) Decode(bytes []byte) (emitter.DomainEvent, error) {
	var e DeviceMarkedEvent
	err := json.Unmarshal(bytes, &e)
	return &e, err
}

func (h *DeviceMarkedEvent) Encode() ([]byte, error) {
	return json.Marshal(h)
}
