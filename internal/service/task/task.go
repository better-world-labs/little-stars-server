package task

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/pkg/crypto"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type RangeType int8

const (
	RANGE_All        RangeType = 0
	RANGE_WHITE_LIST RangeType = 1
	RANGE_BLACK_LIST RangeType = 2

	TaskTableName      = "task"
	RangeAedTableName  = "task_range_aed"
	RangeUserTableName = "task_range_user"

	True  = 1 //未读
	False = 0 //已读

	PicketNone              = 1 //未纠察
	PicketLastTwiceConflict = 2 //最近两次纠察结果不一致
	PicketLastTwiceFalse    = 3 //最近两次纠察结果都为 "设备不存在"
)

type RangeAed struct {
	Id       int64
	TaskId   int64
	DeviceId string
}

type RangeUser struct {
	Id     int64
	TaskId int64
	UserId int64
}

type PicketConditionType int8

type Task struct {
	Id              int64
	Name            string `json:"name" xorm:"name"`
	Image           string `json:"image" xorm:"image"`
	Description     string
	Url             string
	UserRangeType   RangeType
	AedRangeType    RangeType
	TimeLimit       time.Duration
	Point           float64 `json:"point,omitempty"`
	ThemeId         int
	CreatedAt       time.Time `json:"createdAt" xorm:"created_at"`
	IsPicket        int8
	PicketCondition PicketConditionType `json:"picketCondition"`
	Level           int8                `json:"level"`
}

func findPicketTasks() ([]*Task, error) {
	conn := db.GetSession()

	var tasks = make([]*Task, 0)
	err := conn.Table(TaskTableName).Where("is_picket=1").Asc("level").Find(&tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

const (
	Nearby   = 500
	AedCount = 100
)

func buildKeyHash(userId int64, link string) string {
	index := strings.Index(link, "?")
	if index != -1 {
		link = link[0:index]
	}
	return crypto.Md5(fmt.Sprintf("%v|%s", userId, link))
}

func (t *Task) genJobByTreasureChest(chest *events.UserOpenTreasureChest) {
	if !t.matchUserRange(chest.UserId) {
		return
	}

	job, keyHash, err := findJobByUserIdAndLink(chest.UserId, chest.Link)

	if err != nil {
		log.Error("find job by key err", err)
		return
	}
	if job != nil {
		log.Info("genJobByTreasureChest: job exited")
		return
	}

	job = &entities.Job{
		TaskId:      t.Id,
		UserId:      chest.UserId,
		IsRead:      False,
		IsTimeLimit: False,
		Status:      JOB_STATUS_INIT,
		CreatedAt:   time.Now(),
		Points:      chest.Points,
		KeyHash:     keyHash,
		Param:       utils.Json(chest),
	}

	err = saveJob(job)
	if err != nil {
		log.Error("saveJob error", err)
	}
}

func (t Task) genJob(userId int64, d *entities.Device) *entities.Job {
	//用户范围匹配
	if !t.matchUserRange(userId) {
		return nil
	}

	//设备范围匹配
	if !t.matchDeviceRange(d.Id) {
		return nil
	}

	//设备状态匹配
	if !t.matchDevicePicketCondition(d) {
		return nil
	}

	job := entities.Job{
		TaskId:      t.Id,
		UserId:      userId,
		IsRead:      False,
		IsTimeLimit: False,
		Status:      JOB_STATUS_INIT,
		DeviceId:    d.Id,
		CreatedAt:   time.Now(),
	}

	if t.TimeLimit > 0 {
		job.IsTimeLimit = True
		job.BeginLimit = time.Now()
		job.EndLimit = job.BeginLimit.Add(t.TimeLimit * time.Second)
	}
	err := saveJob(&job)
	if err != nil {
		return nil
	}
	return &job
}

func findTaskById(taskId int64) *Task {
	var task Task
	_, _ = db.GetSession().Table(TaskTableName).Where("id=?", taskId).Get(&task)
	return &task
}

func (t Task) matchDevicePicketCondition(d *entities.Device) bool {
	condition := interfaces.S.PicketCondition
	switch t.PicketCondition {
	case PicketNone:
		return condition.IsPicketNone(d.Id)
	case PicketLastTwiceConflict:
		return condition.IsLastTwiceConflict(d.Id)
	case PicketLastTwiceFalse:
		return condition.IsLastTwiceFalse(d.Id)
	default:
		return false
	}
}

func (t Task) matchUserRange(userId int64) bool {
	if t.UserRangeType == RANGE_All {
		return true
	}
	if t.UserRangeType == RANGE_WHITE_LIST {
		return userInRange(t.Id, userId)
	}

	if t.UserRangeType == RANGE_BLACK_LIST {
		return !userInRange(t.Id, userId)
	}
	return false
}

func (t Task) matchDeviceRange(deviceId string) bool {
	if t.AedRangeType == RANGE_All {
		return true
	}

	if t.AedRangeType == RANGE_WHITE_LIST {
		return deviceInRange(t.Id, deviceId)
	}

	if t.AedRangeType == RANGE_BLACK_LIST {
		return !deviceInRange(t.Id, deviceId)
	}
	return false
}

func userInRange(taskId int64, userId int64) bool {
	has, err := db.Exist(`select 1 from task_range_user where task_id = ? and user_id = ?`, taskId, userId)
	if err != nil {
		log.Error("userInRange error", err)
	}
	return has
}

func deviceInRange(taskId int64, deviceId string) bool {
	has, err := db.Exist(`select 1 from task_range_aed where task_id = ? and device_id = ?`, taskId, deviceId)
	if err != nil {
		log.Error("userInRange error", err)
	}
	return has
}
