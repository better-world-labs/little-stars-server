package game

import (
	"aed-api-server/internal/interfaces/entities"
)

type (
	IGamePersistence interface {
		GetById(id int64) (*entities.Game, bool, error)

		ListGames() ([]*entities.Game, error)

		Create(game *entities.Game) error

		Update(game *entities.Game) error

		Delete(id int64) error

		ListByIds(ids []int64) ([]*entities.Game, error)
	}

	IPlayerPersistence interface {
		GetById(id int64) (*entities.GameProcess, bool, error)

		CountUserCompleted(gameId int64) (int64, error)

		ListRankUserIds(gameId int64) ([]int64, error)

		GetByGameIdAndUserId(gameId, userId int64) (*entities.GameProcess, bool, error)

		ListUserProcesses(userId int64) ([]*entities.GameProcess, error)

		ListUserProcessesUncompleted(userId int64) ([]*entities.GameProcess, error)

		ListUserProcessesCompleted(userId int64) ([]*entities.GameProcess, error)

		ListTopUserProcessesCompleted(gameId int64, top int) ([]*entities.GameProcess, error)

		ListJoinedGameIds(userId int64) ([]int64, error)

		ListTopBySteps(gameId int64, top int) (games []*entities.GameProcess, err error)

		Create(game *entities.GameProcess) error

		Update(game *entities.GameProcess) error

		CompareAndSet(excepted entities.GameProcess, player *entities.GameProcess) (bool, error)

		Delete(id int64) error
	}
)
