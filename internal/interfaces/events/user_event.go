package events

import (
	"aed-api-server/internal/pkg/domain/emitter"
)

type UserEvent interface {
	emitter.DomainEvent

	//GetUserId 获取用户ID
	GetUserId() int64
}
