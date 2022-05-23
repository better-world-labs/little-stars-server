package server

import (
	_ "aed-api-server/internal/module/exam"
	_ "aed-api-server/internal/service/activity"
	_ "aed-api-server/internal/service/clock_in"
	_ "aed-api-server/internal/service/donation"
	_ "aed-api-server/internal/service/essay"
	"aed-api-server/internal/service/feedback"
	"aed-api-server/internal/service/medal"
	"aed-api-server/internal/service/merit_tree"
	"aed-api-server/internal/service/merit_tree/task_bubble"
	_ "aed-api-server/internal/service/point"
	_ "aed-api-server/internal/service/project"
	"aed-api-server/internal/service/stat"
	"aed-api-server/internal/service/system_config"
	_ "aed-api-server/internal/service/task"
	_ "aed-api-server/internal/service/user"
	_ "aed-api-server/internal/service/user_config"
)

func init() {
	merit_tree.InitMeritTree()
	system_config.Init()
	feedback.Init()
	medal.Init()
	task_bubble.Init()
	stat.Init()
}
