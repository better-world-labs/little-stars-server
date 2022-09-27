package service

import (
	"aed-api-server/internal/interfaces/domains"
	"aed-api-server/internal/interfaces/entities"
)

type IGame interface {
	// ListGamesSorted 读取排序过后的游戏列表
	ListGamesSorted(userId int64) ([]*domains.Game, error)

	// ListGames 读取游戏列表
	ListGames() ([]*domains.Game, error)

	// ListJoinedStartedGames 读取参与的进行中游戏列表
	ListJoinedStartedGames(userId int64) ([]*domains.Game, error)

	// ListUnCompletedGameProcessesDomain  读取未完成的进行中游戏进程
	ListUnCompletedGameProcessesDomain(userId int64) ([]*domains.GameProcess, error)

	GetGameById(id int64) (*domains.Game, bool, error)

	JoinGame(gameId int64, userId int64) error

	UpdateWechatSteps(gameId int64, userId int64, req *entities.WechatDataDecryptReq) (int, error)

	GetGameProcess(gameId, userId int64) (*domains.GameProcess, error)

	GetTopGameProcesses(gameId int64, top int) ([]*domains.GameProcess, error)

	GetGameStat(gameId, userId int64) (*entities.GameStat, error)
}
