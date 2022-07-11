package command

import (
	"aed-api-server/internal/interfaces/facility"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/pkg/redis"
	"gitlab.openviewtech.com/openview-pub/gopkg/inject"
)

func newSender() facility.Sender {
	return &emitter.Sender{}
}

func loadFacility(component *inject.Component) {
	component.Load(newSender(), "sender")
	component.Load(redis.NewCache("little-star"), "cache")
	component.Load(redis.NewLocker("little-star"), "locker")
}
