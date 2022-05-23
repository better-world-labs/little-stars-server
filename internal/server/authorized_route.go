package server

import (
	"aed-api-server/internal/controller"
	"aed-api-server/internal/controller/clock_in"
	"aed-api-server/internal/controller/donation"
	"aed-api-server/internal/controller/exam"
	"aed-api-server/internal/controller/friends"
	"aed-api-server/internal/controller/merit_tree"
	"aed-api-server/internal/controller/point"
	"aed-api-server/internal/controller/project"
	"aed-api-server/internal/controller/task"
	"aed-api-server/internal/controller/user_config"
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/middleware"
	"aed-api-server/internal/module/activity"
	"aed-api-server/internal/module/aid"
	"aed-api-server/internal/module/device"
	"aed-api-server/internal/module/imageprocessing"
	"aed-api-server/internal/module/oss"
	"aed-api-server/internal/module/skill"
	"aed-api-server/internal/module/user"
	"aed-api-server/internal/pkg/config"
	"aed-api-server/internal/pkg/sms"

	"github.com/gin-gonic/gin"
)

func InitAuthorizedRoutes(routerGroup *gin.RouterGroup, c *config.AppConfig, client user.WechatClient) {
	e := routerGroup.Group("", middleware.Authorize)
	imageprocessing.Init()
	sms.InitSmsClient(c.SmsClient)

	us := user.NewService(client)
	userController := user.NewController(c.Backend, us)
	aidController := aid.NewController(aid.NewService(us))
	activityController := activity.NewController()

	skillService := skill.NewService(&c.AliOss)
	achievementController := controller.NewController()
	deviceController := device.NewController(interfaces.S.Device)
	skillController := skill.NewController(skillService)

	e.PUT("/user/position", userController.UpdatePosition)
	e.PUT("/user/mobile", userController.UpdateMobile)
	e.PUT("/user/info", userController.UpdateUserInfo)
	e.GET("/user/info", userController.GetUserInfo)
	e.GET("/user/infos", userController.GetUserInfos)
	e.POST("/user/check-token", userController.CheckUserToken)
	e.POST("/user/events", userController.DealReportedEvents)
	e.GET("/user/charity-card", userController.GetUserCharityCard)

	// aid
	e.POST("/aid/publish", aidController.PublishHelpInfo)
	e.POST("/aid/arrived", aidController.ActionArrived)
	e.POST("/aid/going-to-scene", aidController.ActionGoingToScene)
	e.POST("/aid/called", aidController.ActionCalled)
	e.GET("/aid/me/published", aidController.ListMyHelpInfosPaged)
	e.GET("/aid/me/participated", aidController.ListHelpInfosParticipatedPaged)
	e.GET("/aid/me/all", aidController.ListHelpInfosMyAll)
	e.GET("/aid/track", aidController.GeAidTrackForUser)
	e.POST("/aid/going-to-learning", activityController.GoingToLearning)

	// activity
	e.POST("/aid/scene-report", activityController.CreateScene)
	e.POST("/aid/going-to-device", activityController.GoingToDevice)

	// oss upload
	ossUploadService := oss.NewService(&c.AliOss)
	ossController := oss.NewController(ossUploadService)
	e.POST("/common/photo", ossController.Upload)
	e.GET("/common/upload_token", ossController.GetUploadToken)

	e.GET("/devices/:deviceId/gallery", deviceController.DeviceGallery)

	e.GET("/aed/devices", deviceController.ListDevices)
	e.POST("/devices", deviceController.MarkDevice)
	e.POST("/aed/add/device", deviceController.AddDevice)
	e.GET("/aed/device", deviceController.InfoDevice)
	e.POST("/aed/borrow", activityController.GetDevice)
	e.POST("/aed/device/add_guide", deviceController.AddGuide)
	e.GET("/aed/device/guide", deviceController.ListGuide)
	e.GET("/aed/device/guide_info", deviceController.GetDeviceGuideInfoById)

	// achievement
	e.GET("/achievement/medals", achievementController.ListAllMedalMeta)
	e.GET("/achievement/user-medals", achievementController.ListUsersMedal)
	e.GET("/achievement/medal/toast", achievementController.ListUsersMedalToast)

	e.GET("/skill/cert", skillController.MyCertificate)
	e.GET("/skill/certs", skillController.ListCerts)

	taskC := task.Controller{}
	taskR := e.Group("/task-jobs")
	e.GET("/devices/:deviceId/picket-task", taskC.FindPicketTaskByDeviceId)
	taskR.GET("", taskC.GetUserTasks)
	taskR.GET("/count", taskC.GetUserTaskStat)
	taskR.PUT("/:jobId/read", taskC.ReadTask)

	projectC := project.Controller{}
	projectR := e.Group("/projects")
	projectR.GET("/:projectId/check-video-completed", projectC.CheckVideoCompleted)
	projectR.GET("/:projectId/courses/learnt", projectC.GetLearntCourses)
	projectR.PUT("/:projectId/video/completed", projectC.CompletedProjectVideo)
	projectR.PUT("/courses/:courseId/learnt", projectC.LearntCourseById)
	projectR.GET("/:projectId/level", projectC.GetProjectUserLevel)

	service := skill.NewCertService(skillService)
	examC := exam.NewController(service)
	projectR.GET("/:projectId/exams/submitted", examC.ListSubmitted)
	projectR.GET("/:projectId/exams/unsubmitted/latest", examC.GetUnSubmittedLatest)
	projectR.POST("/:projectId/exams", examC.Start)
	projectR.POST("/exams/:examId/save", examC.Save)
	projectR.POST("/exams/:examId/submit", examC.Submit)
	projectR.GET("/exams/:examId", examC.GetByID)

	userConfigC := user_config.Controller{}
	configR := e.Group("/configs")
	configR.GET("", userConfigC.GetConfig)
	configR.PUT("", userConfigC.PutConfig)

	clockInC := clock_in.Controller{}
	clockInR := e.Group("/clock-ins")
	clockInR.POST("", clockInC.PostClockIn)
	clockInR.GET("", clockInC.GetDeviceClockInList)
	clockInR.GET("/last", clockInC.GetDeviceLastClockIn)

	// point
	pointTreeC := point.Controller{}
	e.Group("/points").
		GET("/details", pointTreeC.GetPointDetail).
		GET("/total", pointTreeC.GetUserPointsCount).
		GET("/strategies", pointTreeC.GetPointsStrategies)

	meritTreeC := merit_tree.Controller{}
	e.Group("/merit-tree").
		GET("", meritTreeC.GetUserMeritTreeInfo).
		GET("/bubbles/count", meritTreeC.GetUserMeritTreeBubblesCount).
		PUT("/bubbles/count", meritTreeC.ReadUserMeritTreeBubblesCount).
		PUT("/bubbles", meritTreeC.AcceptBubblePoints).
		POST("/walk-points", meritTreeC.GetWalkConvertInfo).
		PUT("/walk-points", meritTreeC.ConvertWalkToPoints).
		PUT("/sign-early", meritTreeC.SignEarly)


	friendsController := friends.NewController()
	friendsGroup := e.Group("/friends")
	friendsGroup.GET("/add-points", friendsController.ListFriendsPoints)

	feedbackController := controller.FeedbackController{}
	e.Group("/user-feedbacks").
		POST("/", feedbackController.SubmitFeedback)

	donationController := donation.NewController(conf.MiniProgramQrcode.ContentRootPath)
	e.GET("/donations-donated", donationController.ListDonationsByDonator)
	e.GET("/donations/:id", donationController.GetDonation)
	e.POST("/donations/:id/records", donationController.Donate)
	e.GET("/donations/:id/records", donationController.ListLatestRecords)
	e.GET("/donations/:id/records/top", donationController.TopRecords)
	e.GET("/donations/:id/evidence", donationController.GetDonationEvidence)
	e.POST("/donations/apply", donationController.ApplyDonation)

}
