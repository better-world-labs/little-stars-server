package emitter

type Sender struct{}

func (*Sender) Send(events ...DomainEvent) error {
	return Emit(events...)
}
