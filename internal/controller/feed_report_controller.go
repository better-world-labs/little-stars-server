package controller

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type FeedReport struct {
	Svc service.IFeedReport `inject:"-"`
}

//go:inject-component
func NewFeedReport() *FeedReport {
	return &FeedReport{}
}

func (f FeedReport) MountAuthRouter(r *route.Router) {
	r.POST("/feeds/:id/reports", f.Create)
}

func (f FeedReport) Create(ctx *gin.Context) (interface{}, error) {
	userId := utils.GetContextUserId(ctx)
	var report entities.FeedReport
	feedId, err := utils.GetContextPathParamInt64(ctx, "id")
	if err != nil {
		return nil, err
	}

	err = ctx.ShouldBindJSON(&report)
	if err != nil {
		return nil, err
	}

	report.CreatedBy = userId
	report.FeedId = feedId
	return nil, f.Svc.Create(&report)
}

func (f FeedReport) MountAdminRouter(r *route.Router) {
	//TODO admin route
}
