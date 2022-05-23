package server

import (
	"aed-api-server/internal/controller"
	"aed-api-server/internal/controller/donation"
	"aed-api-server/internal/controller/essay"
	"aed-api-server/internal/middleware"
	"aed-api-server/internal/module/user"
	"aed-api-server/internal/pkg/config"
	"github.com/gin-gonic/gin"
)

func InitAdminRoutes(g *gin.RouterGroup, appConfig *config.AppConfig, wc user.WechatClient) {
	userController := user.NewController(appConfig.Backend, user.NewService(wc))
	g.POST("/user/login", userController.Login)

	g.Use(middleware.NewAuthorizationAdmin(appConfig.Backend).AuthorizeAdmin)
	essayController := &essay.Controller{}
	essayGroup := g.Group("/essays")
	essayGroup.POST("", essayController.Create)
	essayGroup.DELETE("/:id", essayController.Delete)
	essayGroup.GET("/:id", essayController.GetOne)
	essayGroup.PUT("/:id", essayController.Update)
	g.POST("/essays-sorts", essayController.Sort)
	essayGroup.GET("", essayController.List)

	// 积分捐献
	donationController := donation.NewController(appConfig.MiniProgramQrcode.ContentRootPath)
	donationGroup := g.Group("/donations")
	donationGroup.POST("", donationController.AdminCreateDonation)
	donationGroup.GET("/:id", donationController.AdminGetDonation)

	configController := controller.ConfigController{}
	g.POST("/system/configs", configController.PutConfig)

	feedbackController := controller.FeedbackController{}
	g.GET("/user-feedbacks/excel", feedbackController.ExportFeedback)

	g.GET("/stat/kpi", controller.StatController{}.KpiStat)
}
