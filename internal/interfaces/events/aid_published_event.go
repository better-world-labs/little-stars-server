package events

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
)

// HelpInfoPublishedEvent  发布求助信息
type HelpInfoPublishedEvent struct {
	entities.HelpInfo
}

func NewHelpInfoPublishedEvent(info entities.HelpInfo) *HelpInfoPublishedEvent {
	return &HelpInfoPublishedEvent{
		HelpInfo: info,
	}
}

func (*HelpInfoPublishedEvent) Decode(bytes []byte) (emitter.DomainEvent, error) {
	var e HelpInfoPublishedEvent
	err := json.Unmarshal(bytes, &e)
	return &e, err
}

func (h *HelpInfoPublishedEvent) Encode() ([]byte, error) {
	return json.Marshal(h)
}
