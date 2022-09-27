package entities

import (
	"aed-api-server/internal/pkg/global"
	"encoding/json"
	"time"
)

const (
	GameLifecycleNotStart   GameLifecycle = 0
	GameLifecycleProcessing GameLifecycle = 1
	GameLifecycleEnded      GameLifecycle = 2
)

type (
	GameLifecycle uint8

	GameSetting struct {
		ClockInSteps   int   `json:"clockInSteps"`
		TopExtraPoints []int `json:"topExtraPoints"`
	}

	Game struct {
		Id           int64                `json:"id"`
		Name         string               `json:"name"`
		StartAt      global.FormattedTime `json:"startAt"`
		EndAt        global.FormattedTime `json:"endAt"`
		FrontCover   string               `json:"frontCover"`
		Background   string               `json:"background"`
		LineResource string               `json:"lineResource"`
		LineTemp     string               `json:"lineTemp"`
		Points       int                  `json:"points"`
		Steps        int                  `json:"steps"`
		Distance     int                  `json:"distance"`
		Settled      bool                 `json:"-"`
		ClockIns     []*ClockInDefinition `json:"clockIns"`
		Settings     *GameSetting         `json:"settings"`
	}

	ClockInDefinition struct {
		Steps   int      `json:"steps"`
		Name    string   `json:"name"`
		Images  []string `json:"images"`
		Devices []string `json:"devices"`
	}

	GameProcess struct {
		Id               int64                `json:"id"`
		GameId           int64                `json:"gameId"`
		UserId           int64                `json:"userId"`
		HistorySteps     int                  `json:"historySteps"`
		ActiveSteps      int                  `json:"activeSteps"`
		Completed        bool                 `json:"completed"`
		UnlockedClockIn  []int                `xorm:"-" json:"unlockedClockIn"`
		CreatedAt        global.FormattedTime `json:"createdAt"`
		UpdatedAt        global.FormattedTime `json:"updatedAt"`
		HistoryStepsDate time.Time            `json:"-"`
	}

	GameProcessUser struct {
		GameProcess

		User *SimpleUser `json:"user"`
	}

	GameStat struct {
		UsersTotal     int `json:"usersTotal"`
		UsersCompleted int `json:"usersCompleted"`
		CaculatePoints int `json:"caculatePoints"`
		RankPercent    int `json:"rankPercent"`
	}
)

func (e *ClockInDefinition) FromDB(b []byte) error {
	return json.Unmarshal(b, e)
}

func (e *ClockInDefinition) ToDB() ([]byte, error) {
	return json.Marshal(e)
}

func (e *GameSetting) FromDB(b []byte) error {
	return json.Unmarshal(b, e)
}

func (e *GameSetting) ToDB() ([]byte, error) {
	return json.Marshal(e)
}
