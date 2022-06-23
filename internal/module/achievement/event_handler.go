package achievement

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/pkg/domain/emitter"
	"errors"
	log "github.com/sirupsen/logrus"
)

func InitEventHandler() {
	emitter.On(&events.ClockInEvent{}, onDeviceClockIn)
	emitter.On(&events.DeviceMarkedEvent{}, onDeviceMarked)
}

func onDeviceClockIn(event emitter.DomainEvent) error {
	if evt, ok := event.(*events.ClockInEvent); ok {
		err := interfaces.S.Medal.AwardMedalInspector(evt.CreatedBy)
		if err != nil {
			log.Error("onDeviceClockIn error:", err)
		}
	}

	return errors.New("onDeviceClockIn: invalid event type")
}

func onDeviceMarked(event emitter.DomainEvent) error {
	if evt, ok := event.(*events.DeviceMarkedEvent); ok {
		err := interfaces.S.Medal.AwardMedalInspector(evt.CreateBy)
		if err != nil {
			log.Error("onDeviceMarked error:", err)
		}
	}

	return errors.New("onDeviceMarked: invalid event type")
}
