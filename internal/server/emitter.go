package server

import (
	"aed-api-server/internal/module/achievement"
	"aed-api-server/internal/module/activity"
	"aed-api-server/internal/module/device"
	"aed-api-server/internal/module/speech"
	"aed-api-server/internal/service/merit_tree/task_bubble"
	"aed-api-server/internal/service/point"
)

func initEventHandler() {
	device.InitEventHandler()
	point.InitEventHandler()
	speech.InitEventHandler()
	activity.InitEventHandler()
	task_bubble.InitEventHandler()
	achievement.InitEventHandler()
}
