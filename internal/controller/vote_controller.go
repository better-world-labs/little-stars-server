package controller

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/global"
	"errors"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
	"strconv"
)

type VoteController struct {
}

func (c VoteController) GetById(ctx *gin.Context) (interface{}, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	vote, exists, err := interfaces.S.Vote.GetVoteById(id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, err
	}

	return NewVoteDto(vote), nil
}

func (c VoteController) GetVoteOptions(ctx *gin.Context) (interface{}, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	options, err := interfaces.S.Vote.ListVoteOptionsRank(id)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"options": options,
	}, nil
}

func (c VoteController) VotePoints(ctx *gin.Context) (interface{}, error) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)
	idStr := ctx.Param("id")
	var options struct {
		OptionIds []int64 `json:"optionIds" binding:"required"`
	}

	err := ctx.ShouldBindJSON(&options)
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	err = interfaces.S.Vote.VotePoints(id, userId, options.OptionIds)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c VoteController) VoteNormal(ctx *gin.Context) (interface{}, error) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)
	idStr := ctx.Param("id")
	var options struct {
		OptionIds []int64 `json:"optionIds" binding:"required"`
	}

	err := ctx.ShouldBindJSON(&options)
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	err = interfaces.S.Vote.VoteNormal(id, userId, options.OptionIds)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c VoteController) GetUserRemainTimes(ctx *gin.Context) (interface{}, error) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	session := db.GetSession()
	defer session.Close()
	vote, exists, err := interfaces.S.Vote.GetVoteById(id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("not found")
	}

	times, err := interfaces.S.Vote.GetUserRemainTimes(session, vote, userId)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"remainTimes": times,
	}, nil
}

func (c VoteController) GetOptionById(ctx *gin.Context) (interface{}, error) {
	idStr := ctx.Param("optionId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	option, exists, err := interfaces.S.Vote.GetVoteOptionById(id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, err
	}

	return option, nil
}

//go:inject-component
func NewVoteController() *VoteController {
	return &VoteController{}
}

func (c VoteController) MountAuthRouter(r *route.Router) {
	voteGroup := r.Group("/votes")
	voteGroup.GET("/:id/remain-times", c.GetUserRemainTimes)
	voteGroup.POST("/:id/do-vote", c.VoteNormal)
	voteGroup.POST("/:id/do-points-vote", c.VotePoints)
}

func (c VoteController) MountNoAuthRouter(r *route.Router) {
	voteGroup := r.Group("/votes")
	voteGroup.GET("/:id", c.GetById)
	voteGroup.GET("/:id/options", c.GetVoteOptions)
	voteGroup.GET("/options/:optionId", c.GetOptionById)
}

type (
	VoteDto struct {
		Id         int64                `json:"id"`
		Name       string               `json:"name"`
		Image      string               `json:"image"`
		Text       string               `json:"text"`
		MaxTimes   int                  `json:"maxTimes"`
		OptionType int                  `json:"optionType"`
		Status     int                  `json:"status"`
		CreatedAt  global.FormattedTime `json:"createdAt"`
		BeginAt    global.FormattedTime `json:"beginAt"`
		EndAt      global.FormattedTime `json:"endAt"`
	}
)

func NewVoteDto(vote *entities.Vote) *VoteDto {
	return &VoteDto{
		Id:         vote.Id,
		Name:       vote.Name,
		Image:      vote.Image,
		Text:       vote.Text,
		MaxTimes:   vote.MaxTimes,
		OptionType: vote.OptionType,
		Status:     vote.Status(),
		CreatedAt:  global.FormattedTime(vote.CreatedAt),
		BeginAt:    global.FormattedTime(vote.BeginAt),
		EndAt:      global.FormattedTime(vote.EndAt),
	}
}
