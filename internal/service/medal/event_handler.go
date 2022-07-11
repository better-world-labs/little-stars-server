package medal

import (
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/pkg/domain/emitter"
	"errors"
	log "github.com/sirupsen/logrus"
)

func (m *medalImpl) onDeviceClockIn(event emitter.DomainEvent) error {
	log.Info("onDeviceClockIn")

	if evt, ok := event.(*events.ClockInEvent); ok {
		err := m.AwardMedalInspector(evt.CreatedBy)
		if err != nil {
			log.Error("onDeviceClockIn error:", err)
		}

		return nil
	}

	return errors.New("onDeviceClockIn: invalid event type")
}

func (m *medalImpl) onDeviceMarked(event emitter.DomainEvent) error {
	log.Info("onDeviceMarked")

	if evt, ok := event.(*events.DeviceMarkedEvent); ok {
		err := m.AwardMedalInspector(evt.CreateBy)
		if err != nil {
			log.Error("onDeviceMarked error:", err)
		}

		return nil
	}

	return errors.New("onDeviceMarked: invalid event type")
}
