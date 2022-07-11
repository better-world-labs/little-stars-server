package domains

import (
	"aed-api-server/internal/interfaces/entities"
	"time"
)

type Game struct {
	entities.Game
}

func NewGame(game *entities.Game) *Game {
	return &Game{Game: *game}
}

func NewGames(games []*entities.Game) []*Game {
	var dg []*Game

	for _, g := range games {
		dg = append(dg, NewGame(g))
	}

	return dg
}

func (g *Game) DistanceForEveryStep() float64 {
	return float64(g.Game.Distance * 1000 / g.Steps)
}

func (g *Game) Status() entities.GameLifecycle {
	now := time.Now()
	if now.Before(g.StartAt.Time()) {
		return entities.GameLifecycleNotStart
	}

	if now.After(g.EndAt.Time()) {
		return entities.GameLifecycleEnded
	}

	return entities.GameLifecycleProcessing
}
