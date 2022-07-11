package controller

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
	"time"
)

type TraceController struct {
	Service service.TraceService `inject:"-"`
}

//go:inject-component
func NewTraceController() *TraceController {
	return &TraceController{}
}

type Dto struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Sharer   string `json:"sharer"`
	Code     string `json:"code"`
	DeviceID string `json:"deviceId"`
	Source   string `json:"source"`
}

func (t TraceController) MountNoAuthRouter(r *route.Router) {
	g := r.Group("/traces")
	g.POST("/official-accounts", t.Create)
	g.POST("/normal", t.Create)
}

func (t TraceController) Create(ctx *gin.Context) (interface{}, error) {
	param := Dto{}
	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, err
	}

	trace, err := t.Service.Create(param.Code, entities.Trace{
		From:     param.From,
		To:       param.To,
		Sharer:   param.Sharer,
		DeviceID: param.DeviceID,
		Source:   param.Source,
		CreateAt: time.Now(),
	})

	if err != nil {
		return nil, err
	}

	return trace, nil
}
