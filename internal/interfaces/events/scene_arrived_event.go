package events

import (
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// SceneArrivedEvent 到达现场
type SceneArrivedEvent struct {
	Id     string
	Time   time.Time
	Points int
	UserId int64
	Aid    int64
}

func NewSceneArrivedEvent(aidId, userId int64, point int) *SceneArrivedEvent {
	return &SceneArrivedEvent{
		Id:     uuid.NewString(),
		Time:   time.Now(),
		Aid:    aidId,
		UserId: userId,
		Points: point,
	}
}

func (*SceneArrivedEvent) Decode(bytes []byte) (emitter.DomainEvent, error) {
	var e SceneArrivedEvent
	err := json.Unmarshal(bytes, &e)
	return &e, err
}

func (h *SceneArrivedEvent) Encode() ([]byte, error) {
	return json.Marshal(h)
}
