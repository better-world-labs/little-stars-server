package trace

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"time"
)

type Controller struct {
	service service.TraceService
}

func NewTraceController(service service.TraceService) *Controller {
	return &Controller{service: service}
}

type Dto struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Sharer   string `json:"sharer"`
	Code     string `json:"code"`
	DeviceID string `json:"deviceId"`
	Source   string `json:"source"`
}

func (t Controller) Create(ctx *gin.Context) {
	param := Dto{}
	err := ctx.BindJSON(&param)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	trace, err := t.service.Create(param.Code, entities.Trace{
		From:     param.From,
		To:       param.To,
		Sharer:   param.Sharer,
		DeviceID: param.DeviceID,
		Source:   param.Source,
		CreateAt: time.Now(),
	})

	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	response.ReplyOK(ctx, trace)
}
