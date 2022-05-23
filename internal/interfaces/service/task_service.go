package service

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/location"
	"time"
)

type Job struct {
	Id          int64
	TaskId      int64
	UserId      int64
	IsRead      int8
	Status      int8
	CreatedAt   time.Time
	CompletedAt time.Time
	IsTimeLimit int8
	BeginLimit  time.Time
	EndLimit    time.Time
	DeviceId    string
}

type UserTypePage struct {
	Records    []*entities.UserTask `json:"records"`
	NextCursor string               `json:"nextCursor"`
}

type UserTaskStat struct {
	Unread int `json:"unread"`
	Todo   int `json:"todo"`
}

type TaskService interface {
	GetUserTasks(
		userId int64,
		pageSize int,
		statusStr string,
		cursorStr string,
		includeExpired bool,
		coordinate location.Coordinate,
	) (*UserTypePage, error)

	ReadUserTask(userId int64, jobId string) error

	GetUserTaskStat(userId int64) (*UserTaskStat, error)

	GenJobsByUserLocation(userId int64, coordinate location.Coordinate) ([]*Job, error)

	// FindUserTaskByUserIdAndDeviceId 通过设备ID查询
	FindUserTaskByUserIdAndDeviceId(userId int64, deviceId string) (*entities.UserTask, error)

	//CompleteJob 完成任务
	CompleteJob(userId int64, jobId int64) error

	IsTodayHasBubble(userId int64) (bool, error)
}
