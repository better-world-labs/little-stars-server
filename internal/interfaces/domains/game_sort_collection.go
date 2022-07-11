package domains

import "aed-api-server/internal/interfaces/entities"

type GameSortCollection struct {
	*Game
	*entities.GameProcess
}

func NewGameSortCollection(game *Game, process *entities.GameProcess) *GameSortCollection {
	return &GameSortCollection{Game: game, GameProcess: process}
}

func NewGameSortCollections(games []*Game, processes []*entities.GameProcess) []*GameSortCollection {
	var res []*GameSortCollection
	processMap := make(map[int64]*entities.GameProcess, 0)

	for _, p := range processes {
		processMap[p.GameId] = p
	}

	for _, g := range games {
		res = append(res, NewGameSortCollection(g, processMap[g.Id]))
	}

	return res
}

func (g GameSortCollection) Order() int {
	if g.GameProcess == nil && g.Game.Status() == entities.GameLifecycleProcessing {
		return 0
	}

	if g.Game.Status() == entities.GameLifecycleNotStart {
		return 1
	}

	return 3
}
