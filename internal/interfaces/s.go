package interfaces

import (
	"aed-api-server/internal/interfaces/service"
)

type ServiceKeeper struct {
	Aid             service.AidService              `inject:"-"`
	Task            service.TaskService             `inject:"-"`
	Device          service.DeviceService           `inject:"-"`
	PicketCondition service.PicketConditionService  `inject:"-"`
	Project         service.ProjectService          `inject:"-"`
	Course          service.CourseService           `inject:"-"`
	Exam            service.ExamService             `inject:"-"`
	UserConfig      service.UserConfigService       `inject:"-"`
	User            service.UserService             `inject:"-"`
	UserOld         service.UserServiceOld          `inject:"-"`
	ClockIn         service.ClockInService          `inject:"-"`
	Trace           service.TraceService            `inject:"-"`
	Points          service.PointsService           `inject:"-"`
	PointsScheduler service.PointsScheduler         `inject:"-"`
	MeritTree       service.MeritTreeService        `inject:"-"`
	Friends         service.FriendsService          `inject:"-"`
	Essay           service.EssayService            `inject:"-"`
	Walk            service.WalkService             `inject:"-"`
	Activity        service.ActivityService         `inject:"-"`
	Early           service.EarlyService            `inject:"-"`
	Donation        service.DonationService         `inject:"-"`
	Evidence        service.IEvidenceService        `inject:"-"`
	Medal           service.IMedal                  `inject:"-"`
	UserMedal       service.IUserMedal              `inject:"-"`
	Stat            service.StatService             `inject:"-"`
	Config          service.ConfigService           `inject:"-"`
	Feedback        service.FeedbackService         `inject:"-"`
	TaskBubble      service.MeritTreeTaskTaskBubble `inject:"-"`
	Vote            service.VoteService             `inject:"-"`
	SubscribeMsg    service.SubscribeMsg            `inject:"-"`
	Wx              service.IWechat                 `inject:"-"`
	Market          service.MarketService           `inject:"-"`
}

var S = &ServiceKeeper{}
