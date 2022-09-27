package domains

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/utils"
	"errors"
	"math"
	"time"
)

type GameProcess struct {
	entities.GameProcess
	User *entities.SimpleUser
	Game *Game
}

var (
	ErrorInvalidGameLifecycle = errors.New("invalid operate for current game lifecycle")
	ErrorGameAlreadyCompleted = errors.New("already completed")
)

func NewGameProcesses(processes []*entities.GameProcess, games map[int64]*entities.Game) []*GameProcess {
	var res []*GameProcess

	for _, p := range processes {
		res = append(res, NewGameProcess(p, games[p.GameId]))
	}

	return res
}

func NewGameProcess(process *entities.GameProcess, game *entities.Game) *GameProcess {
	return &GameProcess{GameProcess: *process, Game: NewGame(game)}
}

func (e *GameProcess) StepProcess() int {
	return e.HistorySteps + e.ActiveSteps
}

func (e *GameProcess) UnlockedClockInIndex() []bool {
	index := make([]bool, 0)

	for _, clockIn := range e.Game.ClockIns {
		if e.StepProcess() >= clockIn.Steps {
			index = append(index, true)
		}
	}

	return index
}

func (e *GameProcess) AddHistorySteps(walk int) error {
	if e.Game.Status() != entities.GameLifecycleProcessing {
		return ErrorInvalidGameLifecycle
	}

	if e.Completed {
		return ErrorGameAlreadyCompleted
	}

	e.HistorySteps += walk
	if e.StepProcess() >= e.Game.Steps {
		e.HistorySteps = e.Game.Steps - e.ActiveSteps
		e.Completed = true
	}

	return nil
}

func (e *GameProcess) ProcessStep(data *entities.WechatWalkData) error {
	if e.Game == nil {
		return errors.New("error for nil game")
	}

	if e.Game.Status() != entities.GameLifecycleProcessing {
		return ErrorInvalidGameLifecycle
	}

	if e.Completed {
		return ErrorGameAlreadyCompleted
	}

	// 清空老的今日步数
	e.ActiveSteps = 0
	for _, s := range data.StepInfoList {
		t := time.UnixMilli(s.TimeStamp * 1000)

		// 叠加历史步数
		if !utils.IsToday(t) && t.After(e.HistoryStepsDate) {
			err := e.AddHistorySteps(s.Step)
			if err == ErrorGameAlreadyCompleted {
				return nil
			}

			e.HistoryStepsDate = t
			continue
		}

		// 更新今日步数
		if utils.IsToday(t) {
			err := e.UpdateToday(s.Step)
			if err == ErrorGameAlreadyCompleted {
				return nil
			}
		}
	}

	return nil
}

func (e *GameProcess) UpdateToday(steps int) error {
	if e.Game.Status() != entities.GameLifecycleProcessing {
		return ErrorInvalidGameLifecycle
	}

	if e.Completed {
		return ErrorGameAlreadyCompleted
	}

	e.ActiveSteps = steps
	if e.StepProcess() >= e.Game.Steps {
		e.ActiveSteps = e.Game.Steps - e.HistorySteps
		e.Completed = true
	}

	return nil
}

func (e *GameProcess) Distance() int {
	round := math.Round(e.Game.DistanceForEveryStep() * float64(e.StepProcess()))
	return int(round)
}
