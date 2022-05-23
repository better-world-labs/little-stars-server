package achievement

import (
	"aed-api-server/internal/pkg/asserts"
	"encoding/json"
)

type Medal struct {
	ID              int64  `json:"id,string"`
	Name            string `json:"name"`
	Order           int    `json:"order"`
	ActiveIcon      string `json:"activeIcon"`
	InactiveIcon    string `json:"inactiveIcon"`
	Description     string `json:"description"`
	ShareBackground string `json:"shareBackground"`
}

const (
	MedalIdSaveLife      = 1
	MedalIdFirstDonation = 2
	MedalIdInspector     = 3
)

var (
	medalStore []*Medal
	medalMap   map[int64]*Medal
)

func InitMedalStore() {
	if jsonBytes, exists := asserts.GetResource("medal_data.json"); exists {
		err := json.Unmarshal(jsonBytes, &medalStore)
		if err != nil {
			panic(err)
		}

		medalMap = make(map[int64]*Medal)
		for _, m := range medalStore {
			medalMap[m.ID] = m
		}
	} else {
		panic("invalid medal resource")
	}
}

func ListMedals() ([]*Medal, error) {
	return medalStore, nil
}

func GetById(id int64) (*Medal, bool, error) {
	medal, exists := medalMap[id]
	return medal, exists, nil
}
