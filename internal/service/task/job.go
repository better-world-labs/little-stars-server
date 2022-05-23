package task

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/utils"
	"strings"

	"gitlab.openviewtech.com/openview-pub/gopkg/log"
)

const (
	JobTableName        = "task_job"
	JOB_STATUS_INIT     = 0
	JOB_STATUS_COMPLETE = 10
)

func getUserJobs(userId int64, pageSize int, status []int, beginId int64, includeExpired bool) ([]*entities.UserTask, error) {
	defer utils.TimeStat("getUserJobs")()

	conn := db.GetSession()
	jobs := make([]*entities.UserTask, 0)

	condition := ""
	if !includeExpired {
		condition = "and (job.end_limit > now() or job.is_time_limit = 0)"
	}

	sql := `
		select
			job.id as id,
			task.name as name,
			task.point as point,
			job.is_time_limit as is_time_limit,
			job.end_limit - now()  as time_limit,
			job.is_read as is_read,
			task.image as image,
			task.description as description,
			job.device_id as device_id
		from task_job as job
		inner join task
			on task.id = job.task_id
		where
			job.id < ?
			and job.user_id = ?
			` + condition + `
			and status in (` + db.ArrayPlaceholder(len(status)) + `)
		order by id desc
		limit ?`

	params := make([]interface{}, 0)
	params = append(params, beginId, userId)
	for _, it := range status {
		params = append(params, it)
	}
	params = append(params, pageSize)

	err := conn.SQL(sql, params...).Find(&jobs)
	for _, job := range jobs {
		job.IsExpired = job.TimeLimit <= 0
	}

	if err != nil {
		return nil, err
	}
	return jobs, err
}

func findUserTaskByUserIdAndDeviceId(userId int64, deviceId string) (*entities.UserTask, error) {
	job := entities.UserTask{}
	isSuc, err := db.SQL(`
		select
			job.id as id,
			task.name as name,
			task.point as point,
			job.is_time_limit as is_time_limit,
			job.end_limit - now() as time_limit,
			job.is_read as is_read,
			job.status as status,
			task.image as image,
			task.description as description,
			job.device_id as device_id
		from task_job as job
		inner join task
			on task.id = job.task_id
		where
			job.user_id = ?
			and job.status = 0
			and job.device_id = ?
			and task.is_picket = 1
			and (job.end_limit > now() or job.is_time_limit = 0)
		limit 1
	`, userId, deviceId).Get(&job)

	if err != nil {
		return nil, err
	}
	if !isSuc {
		return nil, err
	}
	job.IsExpired = job.TimeLimit <= 0
	return &job, err
}

func arrayPlaceholder(arr []int) string {
	n := len(arr)
	p := make([]string, n)
	for i := 0; i < n; i++ {
		p[i] = "?"
	}
	return strings.Join(p, ",")
}

func readUserTask(jobId int64, userId int64) error {
	conn := db.GetSession()
	_, err := conn.Exec(`update task_job set is_read = 1 where id = ? and user_id = ?`, jobId, userId)
	if err != nil {
		return err
	}
	return nil
}

func findUserTaskStat(userId int64) (*service.UserTaskStat, error) {
	stat := service.UserTaskStat{}
	_, err := db.SQL(`
		select
			count(if(is_read = 0, 1, null)) as unread,
			count(if(status = 0, 1, null)) as todo
		from task_job
		where
			user_id = ?
			and end_limit > now()
	`, userId).Get(&stat)
	return &stat, err
}

func userHasUnCompletePicketJob(userId int64, deviceId string) bool {
	has, err := db.Exist(`
		select
			1 as existed 
		from task_job 
		where 
			user_id = ? 
			and	device_id = ?
			and end_limit > now()
			and status = 0
	`, userId, deviceId)
	if err != nil {
		log.Error("userHasUnCompletePicketJob error:userId=", userId, err)
		return false
	}
	return has
}

func saveJob(job *service.Job) error {
	_, err := db.GetSession().Table(JobTableName).Insert(job)
	if err != nil {
		return err
	}
	return nil
}

func completeJob(userId int64, jobId int64) error {
	_, err := db.Exec(`
		update task_job 
		set 
			status = 10,
			completed_at = now() 
		where
			id = ? 
			and user_id = ? 
			and status = 0`, jobId, userId)
	return err
}
