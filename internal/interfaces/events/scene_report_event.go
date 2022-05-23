package events

import (
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// SceneReportEvent 上传现场播报
type SceneReportEvent struct {
	Id          string    `json:"id"`
	Time        time.Time `json:"time"`
	Aid         int64     `json:"aid"`
	Points      int       `json:"point"`
	UserId      int64     `json:"userId"`
	Description string    `json:"description"`
	Images      []string  `json:"images"`
}

func NewSceneReportEvent(aidId, userId int64, description string, images []string) *SceneReportEvent {
	return &SceneReportEvent{
		Id:          uuid.NewString(),
		Time:        time.Now(),
		Aid:         aidId,
		UserId:      userId,
		Description: description,
		Images:      images,
	}
}

func (*SceneReportEvent) Decode(bytes []byte) (emitter.DomainEvent, error) {
	var e SceneReportEvent
	err := json.Unmarshal(bytes, &e)
	return &e, err
}

func (h *SceneReportEvent) Encode() ([]byte, error) {
	return json.Marshal(h)
}
