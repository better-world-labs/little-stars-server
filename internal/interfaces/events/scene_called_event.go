package events

import (
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// SceneCalledEvent 联系现场
type SceneCalledEvent struct {
	Id     string
	Time   time.Time
	UserId int64 `json:"userId"`
	Aid    int64 `json:"aid"`
}

func NewSceneCalledEvent(aidId, userId int64) *SceneCalledEvent {
	return &SceneCalledEvent{
		Id:     uuid.NewString(),
		Time:   time.Now(),
		Aid:    aidId,
		UserId: userId,
	}
}

func (*SceneCalledEvent) Decode(bytes []byte) (emitter.DomainEvent, error) {
	var e SceneCalledEvent
	err := json.Unmarshal(bytes, &e)
	return &e, err
}

func (h *SceneCalledEvent) Encode() ([]byte, error) {
	return json.Marshal(h)
}
