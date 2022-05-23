package events

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
)

type ClockInEvent struct {
	*entities.ClockInBaseInfo
}

func (e *ClockInEvent) Decode(bt []byte) (emitter.DomainEvent, error) {
	var evt ClockInEvent
	err := json.Unmarshal(bt, &evt)
	if err != nil {
		return nil, err
	}
	return &evt, nil
}

func (e *ClockInEvent) Encode() ([]byte, error) {
	bt, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return bt, nil
}
