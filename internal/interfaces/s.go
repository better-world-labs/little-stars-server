package interfaces

import (
	"aed-api-server/internal/interfaces/service"
)

type ServiceKeeper struct {
	Task            service.TaskService
	Device          service.DeviceService
	PicketCondition service.PicketConditionService
	Project         service.ProjectService
	Course          service.CourseService
	Exam            service.ExamService
	UserConfig      service.UserConfigService
	User            service.UserService
	ClockIn         service.ClockInService
	Trace           service.TraceService
	Points          service.PointsService
	PointsScheduler service.PointsScheduler
	MeritTree       service.MeritTreeService
	Friends         service.FriendsService
	Essay           service.EssayService
	Walk            service.WalkService
	Activity        service.ActivityService
	Early           service.EarlyService
	Donation        service.DonationService
	Evidence        service.EvidenceService
	Medal           service.MedalService
	UserMedal       service.UserMedalService
	Stat            service.StatService

	Config service.ConfigService

	Feedback   service.FeedbackService
	TaskBubble service.MeritTreeTaskTaskBubble
}

var S = ServiceKeeper{}

// --- task/service.go ---
//func init() {
//	interfaces.S.Task = &Service{} //初始化服务 👈🏻
//}
//

// --- task/controller.go ---
//type Controller struct {
//	task interfaces.TaskService
//}
//
//func NewController() *Controller {
//	return &Controller{
//		task: interfaces.S.Task, //使用服务 👈🏻
//	}
//}
