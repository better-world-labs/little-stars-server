package controller

import (
	"aed-api-server/internal/interfaces/domains"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type Game struct {
	Svc service.IGame `inject:"-"`
}

//go:inject-component
func NewGame() *Game {
	return &Game{}
}

func (g *Game) MountAuthRouter(r *route.Router) {
	r.GET("/games", g.ListGames)
	r.GET("/games/:id", g.GetGameById)
	r.POST("/join/games/:id", g.JoinGame)
	r.PUT("/games/:id/steps", g.UpdateStep)
	r.GET("/games/:id/my-player", g.GetMyPlayer)
	r.GET("/game-stat/:id/", g.GetGameStat)
	r.GET("/games/:id/players", g.TopPlayers)
	r.GET("/joined/game-ids", g.ProcessingGameId)
}

func (g *Game) MountAdminRouter(r *route.Router) {
	//TODO admin
}

func (g *Game) ListGames(ctx *gin.Context) (interface{}, error) {
	userId := utils.GetContextUserId(ctx)
	games, err := g.Svc.ListGamesSorted(userId)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"games": NewGameDtos(games),
	}, nil
}

func (g *Game) GetGameById(ctx *gin.Context) (interface{}, error) {
	gameId, err := utils.GetContextPathParamInt64(ctx, "id")
	if err != nil {
		return nil, response.NewIllegalArgumentError("invalid path param id")
	}

	game, exists, err := g.Svc.GetGameById(gameId)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("game not found")
	}

	return NewGameDto(game), nil
}

func (g *Game) JoinGame(ctx *gin.Context) (interface{}, error) {
	gameId, err := utils.GetContextPathParamInt64(ctx, "id")
	userId := utils.GetContextUserId(ctx)
	if err != nil {
		return nil, response.NewIllegalArgumentError("invalid path param id")
	}

	err = g.Svc.JoinGame(gameId, userId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (g *Game) UpdateStep(ctx *gin.Context) (interface{}, error) {
	var req entities.WechatDataDecryptReq

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		return nil, response.NewIllegalArgumentError("invalid json param")
	}

	userId := utils.GetContextUserId(ctx)
	gameId, err := utils.GetContextPathParamInt64(ctx, "id")
	if err != nil {
		return nil, response.NewIllegalArgumentError("invalid path param id")
	}

	update, err := g.Svc.UpdateWechatSteps(gameId, userId, &req)
	if err != nil {
		logrus.Errorf("UpdateWechatSteps error: %v", err)
	}

	return map[string]interface{}{
		"updated": update > 0,
		"update":  update,
	}, nil
}

func (g *Game) GetMyPlayer(ctx *gin.Context) (interface{}, error) {
	userId := utils.GetContextUserId(ctx)
	gameId, err := utils.GetContextPathParamInt64(ctx, "id")
	if err != nil {
		return nil, response.NewIllegalArgumentError("invalid gameId")
	}

	process, err := g.Svc.GetGameProcess(gameId, userId)
	if err != nil {
		return nil, err
	}

	return NewGameProcessDto(process), nil
}

func (g *Game) GetGameStat(ctx *gin.Context) (interface{}, error) {
	userId := utils.GetContextUserId(ctx)
	gameId, err := utils.GetContextPathParamInt64(ctx, "id")
	if err != nil {
		return nil, response.NewIllegalArgumentError("invalid gameId")
	}

	stat, err := g.Svc.GetGameStat(gameId, userId)
	if err != nil {
		return nil, err
	}

	return stat, nil

}

func (g *Game) TopPlayers(ctx *gin.Context) (interface{}, error) {
	gameId, err := utils.GetContextPathParamInt64(ctx, "id")
	if err != nil {
		return nil, response.NewIllegalArgumentError("invalid gameId")
	}

	var top struct {
		Top int `form:"top" binding:"required"`
	}

	err = ctx.ShouldBindQuery(&top)
	if err != nil {
		return nil, response.NewIllegalArgumentError("invalid gameId")
	}

	processes, err := g.Svc.GetTopGameProcesses(gameId, top.Top)
	if err != nil {
		return nil, err
	}

	return NewGameProcessDtos(processes), nil
}

func (g *Game) ProcessingGameId(ctx *gin.Context) (interface{}, error) {
	userId := utils.GetContextUserId(ctx)

	processes, err := g.Svc.ListJoinedStartedGames(userId)
	if err != nil {
		return nil, err
	}

	gameIds := make([]int64, 0)
	for _, p := range processes {
		gameIds = append(gameIds, p.Id)
	}

	return map[string]interface{}{
		"games": gameIds,
	}, nil
}

type (
	GameDto struct {
		*domains.Game

		Status entities.GameLifecycle `json:"status"`
	}

	GameProcessDto struct {
		*entities.GameProcess

		User            *entities.SimpleUser `json:"user"`
		UnlockedClockIn []bool               `json:"unlockedClockIn"`
	}
)

func NewGameDto(g *domains.Game) *GameDto {
	dto := &GameDto{Game: g}
	dto.Status = g.Status()
	return dto
}

func NewGameDtos(gs []*domains.Game) []*GameDto {
	dtos := make([]*GameDto, 0)

	for _, g := range gs {
		dtos = append(dtos, NewGameDto(g))
	}

	return dtos
}

func NewGameProcessDto(p *domains.GameProcess) *GameProcessDto {
	return &GameProcessDto{
		GameProcess:     &p.GameProcess,
		User:            p.User,
		UnlockedClockIn: p.UnlockedClockInIndex(),
	}
}

func NewGameProcessDtos(ps []*domains.GameProcess) []*GameProcessDto {
	dtos := make([]*GameProcessDto, 0)

	for _, p := range ps {
		dtos = append(dtos, NewGameProcessDto(p))
	}

	return dtos
}
