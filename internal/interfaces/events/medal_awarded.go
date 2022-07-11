package events

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/domain/emitter"
	"encoding/json"
)

// MedalAwarded 勋章已经颁发
type MedalAwarded struct {
	entities.UserMedal
}

func NewMedalAwarded(medal entities.UserMedal) *MedalAwarded {
	return &MedalAwarded{medal}
}

func (*MedalAwarded) Decode(bytes []byte) (emitter.DomainEvent, error) {
	var e MedalAwarded
	err := json.Unmarshal(bytes, &e)
	return &e, err
}

func (h *MedalAwarded) Encode() ([]byte, error) {
	return json.Marshal(h)
}
