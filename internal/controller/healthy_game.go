package controller

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

//go:inject-component
func NewHealthyGame() *HealthyGame {
	return &HealthyGame{}
}

type HealthyGame struct {
	Game service.IHealthyGame `inject:"-"`
	User service.UserService  `inject:"-"`
}

func (ctr *HealthyGame) MountNoAuthRouter(r *route.Router) {
	group := r.Group("/healthy-games")
	group.GET("/info", ctr.getGameInfo)
}
func (ctr *HealthyGame) MountAuthRouter(r *route.Router) {
	group := r.Group("/healthy-games")
	group.POST("/answers", ctr.commitAnswers)
	group.GET("/result", ctr.getResult)
}

func (ctr *HealthyGame) getGameInfo(context *gin.Context) (interface{}, error) {
	authorization := context.GetHeader(pkg.AuthorizationHeaderKey)
	var userId int64
	if authorization != "" {
		user, err := ctr.User.ParseInfoFromJwtToken(authorization)
		if err != nil {
			log.Infof("ParseInfoFromJwtToken error: %v", err)
		} else {
			userId = user.ID
			context.Set(pkg.AccountIDKey, userId)
			context.Set(pkg.AccountKey, user)
		}
	}

	shareFromUserUid := context.Query("shareFrom")

	return ctr.Game.GetInfo(userId, shareFromUserUid)
}

func (ctr *HealthyGame) commitAnswers(context *gin.Context) (interface{}, error) {
	userId := utils.GetContextUserId(context)

	type Req struct {
		Answers []*entities.Answer `json:"answers"`
	}

	var req Req
	if err := context.BindJSON(&req); err != nil {
		return nil, err
	}

	return ctr.Game.CommitAnswers(userId, req.Answers)
}

func (ctr *HealthyGame) getResult(context *gin.Context) (interface{}, error) {
	userId := utils.GetContextUserId(context)
	return ctr.Game.GetResult(userId)
}
