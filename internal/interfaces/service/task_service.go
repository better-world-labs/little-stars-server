package service

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/location"
)

type TaskService interface {
	GetUserTasks(
		userId int64,
		pageSize int,
		statusStr string,
		cursorStr string,
		includeExpired bool,
		coordinate location.Coordinate,
	) (*entities.UserTypePage, error)

	ReadUserTask(userId int64, jobId string) error

	GetUserTaskStat(userId int64) (*entities.UserTaskStat, error)

	GenJobsByUserLocation(userId int64, coordinate location.Coordinate) ([]*entities.Job, error)

	// FindUserTaskByUserIdAndDeviceId 通过设备ID查询
	FindUserTaskByUserIdAndDeviceId(userId int64, deviceId string) (*entities.UserTask, error)

	//CompleteJob 完成任务
	CompleteJob(userId int64, jobId int64) error

	IsTodayHasBubble(userId int64) (bool, error)

	FindJobByPageLink(userId int64, link string) (*entities.Job, error)
}
