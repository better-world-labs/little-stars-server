package controller

import (
	"aed-api-server/internal/pkg/global"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
	"time"
)

type SystemController struct {
}

//go:inject-component
func NewSystemController() *SystemController {
	return &SystemController{}
}

func (c SystemController) MountNoAuthRouter(r *route.Router) {
	g := r.Group("/system")
	g.GET("/time", c.GetTime)
}

func (SystemController) GetTime(c *gin.Context) (interface{}, error) {
	return map[string]interface{}{
		"time": global.FormattedTime(time.Now()),
	}, nil
}
