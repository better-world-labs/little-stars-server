package server

import (
	"aed-api-server/internal/controller"
	"aed-api-server/internal/controller/clock_in"
	"aed-api-server/internal/controller/donation"
	"aed-api-server/internal/controller/essay"
	"aed-api-server/internal/controller/project"
	"aed-api-server/internal/module/activity"
	"aed-api-server/internal/module/aid"
	"aed-api-server/internal/module/imageprocessing"
	"aed-api-server/internal/module/skill"
	"aed-api-server/internal/module/speech"
	trace2 "aed-api-server/internal/module/trace"
	"aed-api-server/internal/module/user"
	"aed-api-server/internal/pkg/config"
	"aed-api-server/internal/pkg/environment"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/sms"
	"github.com/gin-gonic/gin"
)

func InitRoutes(e *gin.RouterGroup, c *config.AppConfig, client user.WechatClient) {
	prefix := environment.GetDomainPrefix(c.Server.Environment)
	speech.SetAidCaller(speech.NewAidCaller(prefix, speech.NewPathGenerator(speech.NewTokenService())))
	speech.SetUserFinder(speech.NewUserFinder(c.Notifier.UserFinder))
	imageprocessing.Init()
	sms.InitSmsClient(c.SmsClient)

	us := user.NewService(client)
	userController := user.NewController(c.Backend, us)
	aidController := aid.NewController(aid.NewService(us))
	activityController := activity.NewController()

	// open
	e.GET("/health-check", func(c *gin.Context) {
		c.Writer.WriteHeader(200)
	})

	traceService := trace2.NewTraceService(client)
	traceController := trace2.NewTraceController(traceService)
	traceGroup := e.Group("/traces")
	traceGroup.POST("/official-accounts", traceController.Create)
	traceGroup.POST("/normal", traceController.Create)

	skillService := skill.NewService(&c.AliOss)
	achievementController := controller.NewController()
	e.GET("/system/settings", func(ctx *gin.Context) {
		response.ReplyOK(ctx, map[string]interface{}{
			"release": false,
			"tag":     true,
		})
	})
	e.GET("/achievement/create-medal-evidences", achievementController.CreateUsersMedalEvidences)
	shareController := imageprocessing.NewShareController(skillService, c.MiniProgramQrcode, c.AliOss.UploadDir)
	e.GET("image-processing/share/cert", shareController.RenderSharedCert)
	e.GET("image-processing/share/medal", shareController.RenderSharedMedal)
	e.GET("image-processing/share/donation", shareController.RenderShareDonation)
	e.GET("image-processing/share/essay", shareController.RenderShareEssay)
	e.GET("image-processing/resource/medal", shareController.RenderSharedMedal)
	e.GET("image-processing/resource/cert", shareController.RenderResourceCert)
	e.GET("image-processing/resource/evidence", shareController.RenderResourceEvidence)

	e.POST("/user/wechat/app/login", userController.WechatAppLogin)
	e.POST("/user/wechat/mini-program/login", userController.WechatMiniProgramLogin)
	e.POST("/user/wechat/mini-program/login-simple", userController.WechatMiniProgramLoginSimple)
	e.POST("/user/generate-uid", userController.GenerateUid)
	e.GET("/user-counts", userController.CountRegisteredUsers)

	e.GET("/aid/infos", aidController.ListHelpInfosPaged)
	e.GET("/aid/infos-hours", aidController.ListOneHoursInfos)
	e.GET("/aid/info", aidController.GetHelpInfo)
	e.POST("/aid/aid-called/:token/:aid", aidController.ActionAidCalled)

	e.GET("/aid/activities-sorted", activityController.ListActivities)
	e.GET("/aid/activities-sorted/latest", activityController.GetLatestActivity)
	e.GET("/aid/activity", activityController.GetOneByID)
	e.GET("/aid/activities", activityController.GetManyByIDs)

	skillController := skill.NewController(skillService)
	//e.GET("/skill/projects", skillController.ListUserProjectUnauth)

	// 为旧数据生成存证
	e.GET("/skill/create-cert-evidences", skillController.CreateCertEvidences)

	projectC := project.Controller{}
	projectR := e.Group("/projects")
	projectR.GET("/:projectId", projectC.GetProjectById)
	projectR.GET("/:projectId/courses", projectC.GetProjectCourses)
	projectR.GET("/courses/:courseId", projectC.GetProjectCoursesById)
	projectR.GET("/courses/articles/:articleId", projectC.GetArticleById)

	clockInC := clock_in.Controller{}
	clockInR := e.Group("/clock-ins")
	clockInR.GET("/stat", clockInC.GetClockInStat)

	configController := controller.ConfigController{}
	e.Group("/system/configs").
		GET("/", configController.GetConfig).
		GET("/all", configController.GetAllConfig)

	donationController := donation.NewController(conf.MiniProgramQrcode.ContentRootPath)
	e.GET("/donations/apply/explain", donationController.ApplyExplain)
	e.GET("/donations", donationController.ListDonations)

	essayController := &essay.Controller{}
	essayGroup := e.Group("/essays")
	essayGroup.GET("", essayController.List)
	essayGroup.GET("/:id", essayController.GetOne)
}
