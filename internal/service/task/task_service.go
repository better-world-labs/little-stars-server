package task

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/cursor"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/location"
	page "aed-api-server/internal/pkg/query"
	"aed-api-server/internal/pkg/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

const (
	DEFAULT_JOB_ID   int64 = 999999999999
	DEFAULT_PAGESIZE       = 20
	MAX_PAGESIZE           = 100
)

type userJobPageCursor struct {
	BeginId        int64 `json:"i"`
	Status         []int `json:"s"`
	PageSize       int   `json:"p"`
	IncludeExpired bool  `json:"in"`
}

type Service struct {
	User service.UserServiceOld `inject:"-"`
}

//go:inject-component
func NewTaskService() *Service {
	return &Service{}
}

// GetUserTasks 获取任务列表
func (s Service) GetUserTasks(
	userId int64,
	pageSize int,
	statusStr string,
	cursorStr string,
	includeExpired bool,
	coordinate location.Coordinate,
) (*entities.UserTypePage, error) {
	//参数解析
	pageSize, status, beginId, includeExpired := dealParams(pageSize, statusStr, includeExpired, cursorStr)

	//查询jobs
	jobs, err := getUserJobs(userId, pageSize, status, beginId, includeExpired)
	if err != nil {
		return nil, err
	}

	//填充job的设备信息
	err = s.patchJobsDeviceInfo(jobs, coordinate)
	if err != nil {
		return nil, err
	}

	//生成游标
	var nextCursor = ""
	if len(jobs) >= pageSize {
		nextCursor = newCursor(pageSize, status, includeExpired, jobs)
	}

	return &entities.UserTypePage{
		Records:    jobs,
		NextCursor: nextCursor,
	}, nil
}

// ReadUserTask 读任务
func (s Service) ReadUserTask(userId int64, jobId string) error {
	id, err := strconv.ParseInt(jobId, 10, 64)
	if err != nil {
		return err
	}
	return readUserTask(id, userId)
}

func (s Service) GetUserTaskStat(userId int64) (*entities.UserTaskStat, error) {
	return findUserTaskStat(userId)
}

// GenJobsByUserLocation 根据用户的定位信息生成Job
func (s Service) GenJobsByUserLocation(userId int64, coordinate location.Coordinate) ([]*entities.Job, error) {
	if isReportedRecently(userId, coordinate) {
		return nil, nil
	}

	//1.读取所有设备纠察任务
	tasks, err := findPicketTasks()
	if err != nil || len(tasks) == 0 {
		return nil, err
	}

	//2.搜索附近的AED设备
	list, err := interfaces.S.Device.ListDevicesWithoutDistance(coordinate, Nearby, page.Query{Page: 0, Size: AedCount})
	if err != nil {
		return nil, err
	}

	jobs := make([]*entities.Job, 0)

	//4. 遍历设备生成任务
	for _, d := range list {
		//4.1 查询用户是否已经拥有该设备的纠察任务
		if !userHasUnCompletePicketJob(userId, d.Id) {
			//4.2 匹配任务，生成Job
			for _, task := range tasks {
				if job := task.genJob(userId, d); job != nil {
					jobs = append(jobs, job)
				}
			}
		}
	}
	return jobs, nil
}

const recentTime = 5 * time.Minute //"最近"定义为5分钟
const nearby = 100                 //"附近" 定义为5m

//isReportedRecently "最近" 在 "附近" 有上报
func isReportedRecently(userId int64, coordinate location.Coordinate) bool {
	lastRecord, err := findLastedByUserId(userId, time.Now().Add(-recentTime))
	if err != nil {
		log.Error("isReportedRecently call findLastedByUserId error:", err)
	}

	isReported := lastRecord != nil

	if isReported {
		latitude, _ := strconv.ParseFloat(lastRecord.Latitude, 64)
		longitude, _ := strconv.ParseFloat(lastRecord.Longitude, 64)

		lastLocation := location.Coordinate{
			Latitude:  latitude,
			Longitude: longitude,
		}
		dis := coordinate.DistanceOf(lastLocation)

		//100米范围内
		if dis > nearby {
			isReported = false
		}
	}

	if !isReported {
		record, err2 := createRecord(&GenRecord{
			UserId:    userId,
			Longitude: fmt.Sprintf("%f", coordinate.Longitude),
			Latitude:  fmt.Sprintf("%f", coordinate.Latitude),
		})

		if err2 != nil {
			log.Error("GenJobsByUserLocation createRecord error:", err2, record)
		}
	}
	return isReported
}

func (s Service) FindUserTaskByUserIdAndDeviceId(userId int64, deviceId string) (*entities.UserTask, error) {
	if deviceId == "" {
		return nil, nil
	}
	job, err := findUserTaskByUserIdAndDeviceId(userId, deviceId)
	if err != nil {
		return nil, err
	}
	position := s.getUserPosition(userId)
	err = patchOneJobDeviceInfo(position, job)
	if err != nil {
		return nil, err
	}
	return job, err
}

func (s Service) CompleteJob(userId int64, jobId int64) error {
	_, err := completeJob(userId, jobId)
	return err
}

func (s Service) FindJobByPageLink(userId int64, link string) (*entities.Job, error) {
	job, _, err := findJobByUserIdAndLink(userId, link)
	return job, err
}

func (s Service) getUserPosition(userId int64) location.Coordinate {
	position, err := s.User.GetPositionByUserID(userId)
	if err != nil {
		log.Warn("patchJobsDeviceInfo error:", err)
		position.Coordinate = &location.Coordinate{}
	}

	return *position.Coordinate
}

// patchJobsDeviceInfo 补齐 job的设备信息
func (s Service) patchJobsDeviceInfo(jobs []*entities.UserTask, position location.Coordinate) error {
	defer utils.TimeStat("patchJobsDeviceInfo")()

	deviceIds := make([]string, 0, len(jobs))
	for i := range jobs {
		deviceIds = append(deviceIds, jobs[i].DeviceId)
	}

	_map, err := interfaces.S.ClockIn.GetBatchDeviceLastClockIn(position, deviceIds)
	if err != nil {
		return err
	}
	for i := range jobs {
		job := jobs[i]
		if _map[job.DeviceId] != nil {
			job.DeviceInfo = _map[job.DeviceId]
		}
	}
	return nil
}

func (Service) IsTodayHasBubble(userId int64) (bool, error) {
	//1. 有待办
	//2. 今日没有完成

	type Count struct {
		Total int `xorm:"total"`
		Today int `xorm:"today"`
	}
	var count Count
	exist, err := db.SQL(`
		select
			count(if(a.status = 0, 1, null)) total,
			count(if(a.status = 10 and a.completed_at > CURRENT_DATE(), 1, null)) as today
		from task_job as a
		inner join task as b
			on b.id = a.task_id
			and b.not_show = 0
		where
			a.user_id = ?
			and (
				a.is_time_limit = 0
				or (a.begin_limit <= now() and now() <= a.end_limit)
			)
	`, userId).Get(&count)
	if err != nil {
		return false, err
	}

	if !exist {
		return false, nil
	}
	return count.Total > 0 && count.Today == 0, nil
}

func patchOneJobDeviceInfo(from location.Coordinate, job *entities.UserTask) error {
	if job == nil || job.DeviceId == "" {
		return nil
	}
	clockIn, err := interfaces.S.ClockIn.GetDeviceLastClockIn(from, job.DeviceId)
	if err != nil {
		return err
	}
	job.DeviceInfo = clockIn
	return nil
}

func newCursor(size int, status []int, includeExpired bool, jobs []*entities.UserTask) string {
	var beginId = DEFAULT_JOB_ID
	for _, job := range jobs {
		if job.Id < beginId {
			beginId = job.Id
		}
	}
	cursorStr, err := cursor.ToString(&userJobPageCursor{
		BeginId:        beginId,
		Status:         status,
		PageSize:       size,
		IncludeExpired: includeExpired,
	})
	if err != nil {
		return ""
	}
	return cursorStr
}

func dealParams(pageSize int, statusStr string, includeExpired bool, cursorStr string) (int, []int, int64, bool) {
	var (
		status  []int
		beginId int64
	)
	if cursorStr != "" {
		var cur = userJobPageCursor{}
		err := cursor.FromString(cursorStr, &cur)
		if err == nil {
			pageSize = cur.PageSize
			status = cur.Status
			beginId = cur.BeginId
			includeExpired = cur.IncludeExpired
		}
	} else {
		status = splitToInt(statusStr)
	}

	if len(status) == 0 {
		status = []int{JOB_STATUS_INIT, JOB_STATUS_COMPLETE}
	}

	if pageSize == 0 {
		pageSize = DEFAULT_PAGESIZE
	}
	if pageSize > MAX_PAGESIZE {
		pageSize = MAX_PAGESIZE
	}
	if beginId == 0 {
		beginId = DEFAULT_JOB_ID
	}
	return pageSize, status, beginId, includeExpired
}

func splitToInt(str string) []int {
	list := strings.Split(str, ",")
	var intList []int

	for _, it := range list {
		item, err := strconv.Atoi(it)
		if err == nil {
			intList = append(intList, item)
		}
	}
	return intList
}
