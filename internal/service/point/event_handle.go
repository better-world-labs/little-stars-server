package point

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/pkg/domain/emitter"
	log "github.com/sirupsen/logrus"
)

func InitEventHandler() {
	emitter.On(&events.PointsEvent{}, pointsEventHandler)
}

func pointsEventHandler(evt emitter.DomainEvent) error {
	log.Info("[point.EventHandler]", "handlePointEvent")
	event := evt.(*events.PointsEvent)
	_, err := interfaces.S.PointsScheduler.DealPointsEvent(event)
	if err != nil {
		log.Error("[point.EventHandler]", "error: ", err.Error())
	}

	return err
}
