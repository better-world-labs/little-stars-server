package game

import (
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/facility"
	"aed-api-server/internal/pkg/domain/emitter"
	"errors"
	"github.com/sirupsen/logrus"
)

var (
	ErrorInvalidEventType = errors.New("event assert failed, invalid event type")
)

func (g *GameImpl) Listen(on facility.OnEvent) {
	on(&events.GameProcessCompleted{}, g.handleGameProcessCompleted)
	on(&events.ClockInEvent{}, g.handleDeviceClockIn)
}

func (g *GameImpl) handleGameProcessCompleted(e emitter.DomainEvent) error {
	logrus.Info("handleGameProcessCompleted")

	if evt, ok := e.(*events.GameProcessCompleted); ok {
		game, exists, err := g.GetGameById(evt.GameId)
		if err != nil {
			return err
		}

		if !exists {
			return errors.New("game not found")
		}

		topExtraPoints := game.Settings.TopExtraPoints
		top3, err := g.playerPersistence.ListTopUserProcessesCompleted(game.Id, len(topExtraPoints))
		if err != nil {
			return err
		}

		for topIndex, process := range top3 {
			if process.UserId == evt.UserId {
				points := game.Settings.TopExtraPoints[topIndex]
				err := g.pointsAward(evt.GameId, evt.UserId, points)
				if err != nil {
					logrus.Errorf("pointsAward for user %d error: %v\n", evt.UserId, err)
					return err
				}

				err = g.ProgramNotifyWithPoints(evt.UserId, points)
				if err != nil {
					logrus.Errorf("ProgramNotifyFirst3 for user %d error: %v\n", evt.UserId, err)
					return err
				}

				return nil
			}
		}

		return g.ProgramNotifyWithoutPoints(evt.UserId)
	}

	return ErrorInvalidEventType
}

func (g *GameImpl) handleDeviceClockIn(e emitter.DomainEvent) error {
	logrus.Info("handleDeviceClockIn")

	if evt, ok := e.(*events.ClockInEvent); ok {
		joinedGames, err := g.ListJoinedStartedGames(evt.CreatedBy)
		if err != nil {
			return err
		}

		for _, game := range joinedGames {
			err := g.AwardSteps(game.Id, evt.CreatedBy, game.Settings.ClockInSteps)
			if err != nil {
				logrus.Errorf("AwardSteps for %d error: %v", evt.CreatedBy, err)
			}
		}

		return nil //TODO 注意一下，要 return
	}

	return ErrorInvalidEventType
}
