package controller

import (
	"aed-api-server/internal/controller/dto"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type FriendsController struct {
	Service service.FriendsService `inject:"-"`
}

//go:inject-component
func NewFriendsController() *FriendsController {
	return &FriendsController{}
}

func (con *FriendsController) MountAuthRouter(r *route.Router) {
	friendsGroup := r.Group("/friends")
	friendsGroup.GET("/add-points", con.ListFriendsPoints)
}

func (con *FriendsController) ListFriendsPoints(ctx *gin.Context) (interface{}, error) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)
	points, err := con.Service.ListFriendsPoints(userId)
	if err != nil {
		return nil, err
	}

	return dto.NewPointsDto(points), nil
}
