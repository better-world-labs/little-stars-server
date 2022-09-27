package facility

import "aed-api-server/internal/pkg/domain/emitter"

type OnEvent func(event emitter.DomainEvent, handler emitter.DomainEventHandler) emitter.Emitter

type Listener interface {

	//Listen use: Listen(on facility.OnEvent)
	Listen(on OnEvent)
}
