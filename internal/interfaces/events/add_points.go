package events

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
)

type AddPoints struct {
	UserId      int64                    `json:"userId"`
	Points      int                      `json:"points"`
	EventType   entities.PointsEventType `json:"eventType"`
	Description string                   `json:"description"`
}

func (*AddPoints) Decode(bytes []byte) (emitter.DomainEvent, error) {
	var e AddPoints
	err := json.Unmarshal(bytes, &e)
	return &e, err
}

func (h *AddPoints) Encode() ([]byte, error) {
	return json.Marshal(h)
}
