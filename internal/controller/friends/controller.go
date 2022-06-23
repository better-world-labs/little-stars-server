package friends

import (
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type Controller struct {
	Service service.FriendsService `inject:"-"`
}

func NewController() *Controller {
	return &Controller{}
}

func (con *Controller) MountAuthRouter(r *route.Router) {
	friendsGroup := r.Group("/friends")
	friendsGroup.GET("/add-points", con.ListFriendsPoints)
}

func (con *Controller) ListFriendsPoints(ctx *gin.Context) (interface{}, error) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)
	points, err := con.Service.ListFriendsPoints(userId)
	if err != nil {
		return nil, err
	}

	return NewPointsDto(points), nil
}
