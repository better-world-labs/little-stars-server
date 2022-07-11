package medal

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/asserts"
	"encoding/json"
)

var (
	medalStore []*entities.Medal
	medalMap   map[int64]*entities.Medal
)

func InitMedalStore() {
	if jsonBytes, exists := asserts.GetResource("medal_data.json"); exists {
		err := json.Unmarshal(jsonBytes, &medalStore)
		if err != nil {
			panic(err)
		}

		medalMap = make(map[int64]*entities.Medal)
		for _, m := range medalStore {
			medalMap[m.ID] = m
		}
	} else {
		panic("invalid medal resource")
	}
}

func ListMedals() ([]*entities.Medal, error) {
	return medalStore, nil
}

func GetById(id int64) (*entities.Medal, bool, error) {
	medal, exists := medalMap[id]
	return medal, exists, nil
}
