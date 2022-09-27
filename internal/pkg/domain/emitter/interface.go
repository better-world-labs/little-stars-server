package emitter

import "time"

type Decoder interface {
	Decode([]byte) (DomainEvent, error)
}

type DomainEvent interface {
	Decoder
	Encode() ([]byte, error)
}

type DomainEventHandler func(event DomainEvent) error

type Emitter interface {
	Start()
	Close()
	Emit(event ...DomainEvent) error
	On(DomainEvent, ...DomainEventHandler) Emitter
	Off(DomainEvent, ...DomainEventHandler) Emitter
}

type DelayEmitter interface {
	DelayEmit(duration time.Duration, event ...DomainEvent) error
}
