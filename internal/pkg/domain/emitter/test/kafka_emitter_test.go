package test

import (
	"aed-api-server/internal/pkg/domain/emitter"
	"testing"
)

func Test_getValueTypeName(t *testing.T) {
	type XXEvent struct {
		emitter.DomainEvent
	}
	//name := getValueTypeName(XXEvent{})
	//assert.Equal(t, "XXEvent", name)
}
