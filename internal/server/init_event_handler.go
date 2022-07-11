package server

import (
	"aed-api-server/internal/interfaces/facility"
	"aed-api-server/internal/pkg/domain/emitter"
)

func initEventHandler() {
	for _, instances := range component.GetInstances() {
		if o, ok := instances.Obj.(facility.Listener); ok {
			o.Listen(func(event emitter.DomainEvent, handler emitter.DomainEventHandler) {
				emitter.On(event, handler)
			})
		}
	}
}
