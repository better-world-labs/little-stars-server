package game

import (
	"aed-api-server/internal/interfaces/domains"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/pkg/global"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"errors"
	"sort"
	"time"
)

type GameImpl struct {
	persistence       IGamePersistence
	playerPersistence IPlayerPersistence

	User service.UserService `inject:"-"`
}

//go:inject-component
func NewGameImp() service.IGame {
	return &GameImpl{
		persistence:       NewGamePersistence(),
		playerPersistence: NewPlayerPersistence(),
	}
}

func (g *GameImpl) updateSettled(gameId int64, settled bool) error {
	game, exists, err := g.GetGameById(gameId)
	if err != nil {
		return err
	}

	if !exists {
		return response.ErrorNotFound
	}

	game.Settled = settled
	return g.persistence.Update(&game.Game)
}
func (g *GameImpl) listEndUnsettledGames() ([]*domains.Game, error) {
	var res []*domains.Game

	games, err := g.listEndedGames()
	if err != nil {
		return nil, err
	}

	for _, game := range games {
		if !game.Settled {
			res = append(res, game)
		}
	}

	return res, nil
}

func (g *GameImpl) listEndedGames() ([]*domains.Game, error) {
	var res []*domains.Game

	games, err := g.persistence.ListGames()
	if err != nil {
		return nil, err
	}

	for _, game := range games {
		domain := domains.NewGame(game)
		if domain.Status() == entities.GameLifecycleEnded {
			res = append(res, domain)
		}
	}

	return res, nil
}

func (g *GameImpl) ListGames() ([]*domains.Game, error) {
	games, err := g.persistence.ListGames()
	if err != nil {
		return nil, err
	}

	return domains.NewGames(games), nil
}

func (g *GameImpl) ListGamesSorted(userId int64) ([]*domains.Game, error) {
	res, err := utils.PromiseAll(func() (interface{}, error) {
		return g.ListGames()
	}, func() (interface{}, error) {
		return g.playerPersistence.ListUserProcesses(userId)
	})

	if err != nil {
		return nil, err
	}
	games := res[0].([]*domains.Game)
	processes := res[1].([]*entities.GameProcess)

	gameMap := make(map[int64]*domains.Game, 0)
	for _, g := range games {
		gameMap[g.Id] = g
	}

	sortCollection := domains.NewGameSortCollections(games, processes)

	sort.Slice(sortCollection, func(i, j int) bool {
		return sortCollection[i].Order() < sortCollection[j].Order()
	})

	games = nil
	for _, s := range sortCollection {
		games = append(games, s.Game)
	}

	return games, err
}

func (g *GameImpl) GetGameById(id int64) (*domains.Game, bool, error) {
	game, exists, err := g.persistence.GetById(id)
	if err != nil {
		return nil, exists, err
	}

	if !exists {
		return nil, exists, response.ErrorNotFound
	}

	return domains.NewGame(game), exists, err
}

func (g *GameImpl) JoinGame(gameId int64, userId int64) error {
	has, err := g.HasUserStartedGames(userId)
	if err != nil {
		return err
	}

	if has {
		return errors.New("user has joined started games")
	}

	return g.doJoinGame(gameId, userId)
}

func (g *GameImpl) AwardSteps(gameId int64, userId int64, walks int) error {
	process, err := g.GetGameProcess(gameId, userId)
	if err != nil {
		return err
	}

	before := *process
	err = process.AddHistorySteps(walks)
	if err != nil {
		return err
	}

	err = g.CompareAndSetProcess(before, process)
	if err == response.ErrorConcurrentOperation {
		time.Sleep(300 * time.Millisecond)
		return g.AwardSteps(gameId, userId, walks)
	}

	return err
}

func (g *GameImpl) UpdateWechatSteps(gameId int64, userId int64, req *entities.WechatDataDecryptReq) (int, error) {
	walks, err := g.User.GetWalks(req)
	if err != nil {
		return 0, err
	}

	process, err := g.GetGameProcess(gameId, userId)
	if err != nil {
		return 0, err
	}

	before := *process
	err = process.ProcessStep(walks)
	if err != nil {
		return 0, err
	}

	err = g.CompareAndSetProcess(before, process)
	if err != nil {
		return 0, err
	}

	update := process.StepProcess() - before.StepProcess()
	if update < 0 {
		update = 0
	}

	return update, nil
}

func (g *GameImpl) CompareAndSetProcess(excepted domains.GameProcess, process *domains.GameProcess) error {
	process.UpdatedAt = global.FormattedTime(time.Now())
	updated, err := g.playerPersistence.CompareAndSet(excepted.GameProcess, &process.GameProcess)
	if err != nil {
		return err
	}

	if !updated {
		return response.ErrorConcurrentOperation
	}

	if process.Completed {
		return emitter.Emit(&events.GameProcessCompleted{
			UserId:      process.UserId,
			GameId:      process.GameId,
			CompletedAt: process.UpdatedAt.Time(),
		})
	}

	return nil
}

func (g *GameImpl) GetGameProcess(gameId, userId int64) (*domains.GameProcess, error) {
	res, err := utils.PromiseAll(func() (interface{}, error) {
		user, exists, err := g.User.GetUserById(userId)
		if err != nil {
			return nil, err
		}

		if !exists {
			return nil, response.ErrorNotFound
		}

		return user, nil
	}, func() (interface{}, error) {
		gamePlayer, exists, err := g.playerPersistence.GetByGameIdAndUserId(gameId, userId)
		if err != nil {
			return nil, err
		}

		if !exists {
			return nil, response.ErrorNotFound
		}

		return gamePlayer, nil
	})

	if err != nil {
		return nil, err
	}

	enhanced, err := g.enhanceGameProcess(res[1].(*entities.GameProcess))
	GameProcessAddUser(enhanced, res[0].(*entities.SimpleUser))
	return enhanced, nil
}

func (g *GameImpl) GetTopGameProcesses(gameId int64, top int) ([]*domains.GameProcess, error) {
	players, err := g.playerPersistence.ListTopBySteps(gameId, top)
	if err != nil {
		return nil, err
	}

	userIds := playersMapUserIds(players)
	users, err := g.User.GetMapUserByIDs(userIds)
	if err != nil {
		return nil, err
	}

	enhanced, err := g.enhanceGameProcesses(players)
	if err != nil {
		return nil, err
	}

	GameProcessesAddUsers(enhanced, users)
	return enhanced, nil
}

func (g *GameImpl) GetGameStat(gameId, userId int64) (*entities.GameStat, error) {
	res, err := utils.PromiseAll(func() (interface{}, error) {
		return g.playerPersistence.ListRankUserIds(gameId)
	}, func() (interface{}, error) {
		return g.playerPersistence.CountUserCompleted(gameId)
	}, func() (interface{}, error) {
		game, exists, err := g.GetGameById(gameId)
		if !exists && err == nil {
			err = response.ErrorNotFound
		}

		return game, err
	})
	if err != nil {
		return nil, err
	}

	rankIds := res[0].([]int64)
	completedCount := res[1].(int64)
	game := res[2].(*domains.Game)

	return &entities.GameStat{
		UsersTotal:     len(rankIds),
		UsersCompleted: int(completedCount),
		CaculatePoints: caculateGamePoints(int(completedCount), game.Points),
		RankPercent:    caculateUserRankPercents(userId, rankIds),
	}, nil
}

func (g *GameImpl) doJoinGame(gameId, userId int64) error {
	return g.playerPersistence.Create(&entities.GameProcess{
		GameId:           gameId,
		UserId:           userId,
		CreatedAt:        global.FormattedTime(time.Now()),
		UpdatedAt:        global.FormattedTime(time.Now()),
		HistoryStepsDate: time.Now().Add(-24 * time.Hour),
	})
}

func (g *GameImpl) HasUserStartedGames(userId int64) (bool, error) {
	joinedGames, err := g.ListJoinedStartedGames(userId)
	if err != nil {
		return false, nil
	}

	return len(joinedGames) > 0, nil
}

func (g *GameImpl) enhanceGameProcess(process *entities.GameProcess) (*domains.GameProcess, error) {
	p, err := g.enhanceGameProcesses([]*entities.GameProcess{process})
	return p[0], err
}

func (g *GameImpl) enhanceGameProcesses(processes []*entities.GameProcess) ([]*domains.GameProcess, error) {
	var gameIds []int64
	for _, p := range processes {
		gameIds = append(gameIds, p.GameId)
	}

	games, err := g.persistence.ListByIds(gameIds)
	if err != nil {
		return nil, err
	}

	gameMap := make(map[int64]*entities.Game)
	for _, game := range games {
		gameMap[game.Id] = game
	}

	return domains.NewGameProcesses(processes, gameMap), nil
}

func (g *GameImpl) ListUnCompletedGameProcessesDomain(userId int64) ([]*domains.GameProcess, error) {
	processes, err := g.playerPersistence.ListUserProcessesUncompleted(userId)
	if err != nil {
		return nil, err
	}

	return g.enhanceGameProcesses(processes)
}

func (g *GameImpl) ListJoinedStartedGames(userId int64) ([]*domains.Game, error) {
	games, err := g.ListJoinedGames(userId)
	if err != nil {
		return nil, err
	}

	var availableGames []*domains.Game
	for _, game := range games {
		gameDomain := domains.NewGame(game)
		if gameDomain.Status() == entities.GameLifecycleProcessing {
			availableGames = append(availableGames, gameDomain)
		}
	}

	return availableGames, nil
}

func (g *GameImpl) ListJoinedGames(userId int64) ([]*entities.Game, error) {
	gameIds, err := g.playerPersistence.ListJoinedGameIds(userId)
	if err != nil {
		return nil, err
	}

	return g.persistence.ListByIds(gameIds)
}
