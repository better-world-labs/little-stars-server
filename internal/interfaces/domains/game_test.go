package domains

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/global"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var (
	gameJson = `{"id": 1,
    "name": "哟西",
    "startAt": "2022-06-28 13:54:52",
    "endAt": "2030-06-29 13:55:00",
    "frontCover": "https://xxx/xxx",
    "background": "https://xxx/xxx",
    "pathData": "https://xxx/xxx",
    "points": 50000,
    "steps": 2000,
    "clockIns": [
      {
        "steps": 100,
        "name": "打卡点1",
        "images": [
          "https://xxx/xxx"
        ],
        "devices": [
          "xxx"
        ]
      }
    ]
  }`
	mockWalkData = &entities.WechatWalkData{
		StepInfoList: []*entities.WechatStepInfo{
			{
				TimeStamp: time.Now().Add(-72 * time.Hour).Unix(),
				Step:      1000,
			},
			{
				TimeStamp: time.Now().Add(-48 * time.Hour).Unix(),
				Step:      300,
			},
			{
				TimeStamp: time.Now().Add(-24 * time.Hour).Unix(),
				Step:      300,
			},
			{
				TimeStamp: time.Now().Unix(),
				Step:      300,
			},
		},
	}
)

// TestUpdateStepUnCompleted 正常加步数不完成
func TestProcessUnCompleted(t *testing.T) {
	var game entities.Game
	err := json.Unmarshal([]byte(gameJson), &game)
	require.Nil(t, err)
	tim, err := time.Parse("2006-01-02", "2021-07-04")
	process := NewGameProcess(&entities.GameProcess{
		GameId:           game.Id,
		UserId:           50,
		HistoryStepsDate: tim,
	}, &game)

	err = process.ProcessStep(mockWalkData)
	require.Nil(t, err)
	assert.Equal(t, 1600, process.HistorySteps)
	assert.Equal(t, 300, process.ActiveSteps)
}

// TestProcessHistoryComplete 历史步数完成
func TestProcessHistoryComplete(t *testing.T) {
	var game entities.Game
	err := json.Unmarshal([]byte(gameJson), &game)
	require.Nil(t, err)
	tim, err := time.Parse("2006-01-02", "2021-07-04")
	process := NewGameProcess(&entities.GameProcess{
		GameId:           game.Id,
		UserId:           50,
		HistorySteps:     400,
		HistoryStepsDate: tim,
	}, &game)

	err = process.ProcessStep(mockWalkData)
	require.Nil(t, err)
	require.Equal(t, 2000, process.HistorySteps)
	require.Equal(t, 0, process.ActiveSteps)
	require.Equal(t, true, process.Completed)
}

// TestProcessTodayComplete 今日步数完成
func TestProcessTodayComplete(t *testing.T) {
	var game entities.Game
	err := json.Unmarshal([]byte(gameJson), &game)
	require.Nil(t, err)
	tim, err := time.Parse("2006-01-02", "2021-07-04")
	process := NewGameProcess(&entities.GameProcess{
		GameId:           game.Id,
		UserId:           50,
		HistorySteps:     200,
		HistoryStepsDate: tim,
	}, &game)

	err = process.ProcessStep(mockWalkData)
	require.Nil(t, err)
	require.Equal(t, 1800, process.HistorySteps)
	require.Equal(t, 200, process.ActiveSteps)
	require.Equal(t, true, process.Completed)
}

// TestGameNotStarted 游戏未开始
func TestGameNotStarted(t *testing.T) {
	var game entities.Game
	err := json.Unmarshal([]byte(gameJson), &game)
	require.Nil(t, err)
	game.StartAt = global.FormattedTime(time.Now().Add(24 * time.Hour))
	tim, err := time.Parse("2006-01-02", "2022-07-04")
	process := NewGameProcess(&entities.GameProcess{
		GameId:           game.Id,
		UserId:           50,
		HistorySteps:     100,
		HistoryStepsDate: tim,
	}, &game)

	err = process.ProcessStep(mockWalkData)
	require.NotNil(t, err)
}

// TestGameEnded 游戏已经结束
func TestGameEnded(t *testing.T) {
	var game entities.Game
	err := json.Unmarshal([]byte(gameJson), &game)
	require.Nil(t, err)
	game.EndAt = global.FormattedTime(time.Now().Add(-24 * time.Hour))
	tim, err := time.Parse("2006-01-02", "2022-07-04")
	process := NewGameProcess(&entities.GameProcess{
		GameId:           game.Id,
		UserId:           50,
		HistorySteps:     100,
		HistoryStepsDate: tim,
	}, &game)

	err = process.ProcessStep(mockWalkData)
	require.NotNil(t, err)
}

// TestExtra 额外测试
func TestExtra(t *testing.T) {
	var game entities.Game
	err := json.Unmarshal([]byte(gameJson), &game)
	require.Nil(t, err)
	file, err := os.ReadFile("test_walk.json")
	require.Nil(t, err)
	var walk entities.WechatWalkData
	err = json.Unmarshal(file, &walk)
	require.Nil(t, err)

	game.EndAt = global.FormattedTime(time.Now().Add(24 * time.Hour))
	updateAt, err := time.Parse("2006-01-02", "2022-07-04")
	require.Nil(t, err)
	process := NewGameProcess(&entities.GameProcess{
		GameId:           game.Id,
		UserId:           50,
		HistorySteps:     500,
		ActiveSteps:      5475,
		HistoryStepsDate: updateAt,
	}, &game)

	err = process.ProcessStep(&walk)
	require.Nil(t, err)
	require.Equal(t, true, process.HistorySteps > 500)
}
