package command

import (
	"aed-api-server/internal/controller"
	"aed-api-server/internal/controller/clock_in"
	"aed-api-server/internal/controller/donation"
	"aed-api-server/internal/controller/essay"
	"aed-api-server/internal/controller/exam"
	"aed-api-server/internal/controller/friends"
	"aed-api-server/internal/controller/merit_tree"
	"aed-api-server/internal/controller/point"
	"aed-api-server/internal/controller/project"
	"aed-api-server/internal/controller/task"
	"aed-api-server/internal/controller/user_config"
	"aed-api-server/internal/module/aid"
	"aed-api-server/internal/module/device"
	"aed-api-server/internal/module/imageprocessing"
	"aed-api-server/internal/module/oss"
	"aed-api-server/internal/module/user"
	"aed-api-server/internal/pkg/config"
	"gitlab.openviewtech.com/openview-pub/gopkg/inject"
)

func LoadController(c *config.AppConfig, component *inject.Component) {
	component.Load(user.NewController(c.Backend))
	component.Load(aid.NewController())
	component.Load(controller.NewActivityController())
	component.Load(controller.NewTraceController())
	component.Load(controller.NewAchievementController())
	component.Load(controller.NewSkillController())
	component.Load(project.NewController())
	component.Load(point.NewController())
	component.Load(controller.NewConfigController())
	component.Load(donation.NewController(c.MiniProgramQrcode.ContentRootPath))
	component.Load(imageprocessing.NewShareController(c.MiniProgramQrcode, c.AliOss.UploadDir))
	component.Load(essay.NewEssayController())
	component.Load(controller.NewVoteController())
	component.Load(controller.NewMarketController())
	component.Load(controller.NewFeedbackController())
	component.Load(controller.NewStatController())
	component.Load(oss.NewController(oss.NewService(&c.AliOss)))
	component.Load(device.NewController())
	component.Load(task.NewController())
	component.Load(exam.NewController())
	component.Load(user_config.NewController())
	component.Load(clock_in.NewController())
	component.Load(merit_tree.NewController())
	component.Load(friends.NewController())
}
