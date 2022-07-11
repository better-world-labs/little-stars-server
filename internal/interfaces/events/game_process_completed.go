package events

import (
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
	"time"
)

type GameProcessCompleted struct {
	GameId      int64     `json:"gameId"`
	UserId      int64     `json:"userId"`
	CompletedAt time.Time `json:"completedAt"`
}

func (*GameProcessCompleted) Decode(bytes []byte) (emitter.DomainEvent, error) {
	var e GameProcessCompleted
	err := json.Unmarshal(bytes, &e)
	return &e, err
}

func (h *GameProcessCompleted) Encode() ([]byte, error) {
	return json.Marshal(h)
}
