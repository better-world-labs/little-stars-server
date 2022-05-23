package emitter

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
	On(DomainEvent, ...DomainEventHandler)
	Off(DomainEvent, ...DomainEventHandler)
}
