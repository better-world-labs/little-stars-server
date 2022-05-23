package events

import (
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// GoingToSceneEvent 正在前往现场
type GoingToSceneEvent struct {
	Id     string
	Time   time.Time
	Aid    int64
	UserId int64
}

func NewGoingToSceneEvent(aidId, userId int64) *GoingToSceneEvent {
	return &GoingToSceneEvent{
		Id:     uuid.NewString(),
		Time:   time.Now(),
		Aid:    aidId,
		UserId: userId,
	}
}

func (*GoingToSceneEvent) Decode(bytes []byte) (emitter.DomainEvent, error) {
	var e GoingToSceneEvent
	err := json.Unmarshal(bytes, &e)
	return &e, err
}

func (h *GoingToSceneEvent) Encode() ([]byte, error) {
	return json.Marshal(h)
}
