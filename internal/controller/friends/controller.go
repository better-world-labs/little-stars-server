package friends

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service service.FriendsService
}

func NewController() *Controller {
	return &Controller{interfaces.S.Friends}
}

func (con *Controller) ListFriendsPoints(ctx *gin.Context) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)
	points, err := con.service.ListFriendsPoints(userId)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	response.ReplyOK(ctx, NewPointsDto(points))
}
