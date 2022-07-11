package point

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/facility"
	"aed-api-server/internal/pkg/domain/emitter"
	"errors"
	"github.com/sirupsen/logrus"
)

var (
	ErrorInvalidEventType = errors.New("event assert failed, invalid event type")
)

func (*service) Listen(on facility.OnEvent) {
	on(&events.PointsEvent{}, handlePointsEvent)
	on(&events.AddPoints{}, handleAddPoints)
}

func handlePointsEvent(evt emitter.DomainEvent) error {
	logrus.Info("[point.EventHandler]", "handlePointEvent")
	event := evt.(*events.PointsEvent)
	_, err := interfaces.S.PointsScheduler.DealPointsEvent(event)
	if err != nil {
		logrus.Error("[point.EventHandler]", "error: ", err.Error())
	}

	return err
}

func handleAddPoints(e emitter.DomainEvent) error {
	logrus.Info("handleAddPoints")

	if evt, ok := e.(*events.AddPoints); ok {
		return interfaces.S.Points.AddPoint(
			evt.UserId,
			evt.Points,
			evt.Description,
			evt.EventType,
		)
	}

	return ErrorInvalidEventType
}
