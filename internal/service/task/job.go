package task

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/pkg/url_args"
	"aed-api-server/internal/pkg/utils"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

const (
	JobTableName        = "task_job"
	JOB_STATUS_INIT     = 0
	JOB_STATUS_COMPLETE = 10

	ScanVideoToCompletedProcess = 90
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
			and not_show = 0
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

func readUserTask(jobId int64, userId int64) error {
	conn := db.GetSession()
	_, err := conn.Exec(`update task_job set is_read = 1 where id = ? and user_id = ?`, jobId, userId)
	if err != nil {
		return err
	}
	return nil
}

func findUserTaskStat(userId int64) (*entities.UserTaskStat, error) {
	stat := entities.UserTaskStat{}
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

func saveJob(job *entities.Job) error {
	_, err := db.GetSession().Table(JobTableName).Insert(job)
	if err != nil {
		return err
	}
	return nil
}

func completeJob(userId int64, jobId int64) (bool, error) {
	rst, err := db.Exec(`
		update task_job 
		set 
			status = 10,
			completed_at = now() 
		where
			id = ? 
			and user_id = ? 
			and status = 0`, jobId, userId)

	affected, err := rst.RowsAffected()
	return affected == 1, err
}

func findJobByUserIdAndLink(userId int64, link string) (job *entities.Job, keyHash string, err error) {
	keyHash = buildKeyHash(userId, link)

	jobs := make([]*entities.Job, 0)
	err = db.Table("task_job").Where("key_hash = ? and status = 0", keyHash).Find(&jobs)
	if err != nil {
		return nil, keyHash, err
	}

	for _, job = range jobs {
		var evt events.UserOpenTreasureChest
		err = json.Unmarshal([]byte(job.Param), &evt)
		if err != nil {
			log.Errorf("json.Unmarshal([]byte(j.Param), &evt) err=%v, j=%v", err, job)
			continue
		}

		if url_args.Compare(link, evt.Link, evt.LinkArgs) {
			return job, keyHash, nil
		}
	}
	return nil, keyHash, nil
}

const ScanPageTaskId = 10

func scanPage(userId int64, pageUrl string) {
	job, _, err := findJobByUserIdAndLink(userId, pageUrl)
	if err != nil {
		log.Error("findJobByUserIdAndLink(userId, pageUrl)", err)
		return
	}

	if job == nil {
		return
	}
	if job.TaskId != ScanPageTaskId {
		return
	}

	suc, err := completeJob(userId, job.Id)
	if err != nil {
		log.Error("completeJob(userId, job.Id)", err)
		return
	}
	if suc {
		var evt = events.UserOpenTreasureChest{}
		err = json.Unmarshal([]byte(job.Param), &evt)
		if err != nil {
			log.Error("json.Unmarshal([]byte(job.Param), &evt)", err)
			return
		}

		pointsEvt := interfaces.S.PointsScheduler.BuildPointsEventTypeReward(userId, job.Id, evt.Points, evt.TreasureChestName)
		err := emitter.Emit(pointsEvt)
		if err != nil {
			log.Error("emitter.Emit(pointsEvt)", err)
			return
		}
	}
}

const ScanVideoTaskId = 11

func scanVideo(userId int64, pageUrl string, process int) {
	job, _, err := findJobByUserIdAndLink(userId, pageUrl)
	if err != nil {
		log.Error("findJobByUserIdAndLink(userId, pageUrl)", err)
		return
	}
	if job == nil {
		return
	}
	if job.TaskId != ScanVideoTaskId {
		return
	}

	if process < ScanVideoToCompletedProcess {
		return
	}

	suc, err := completeJob(userId, job.Id)
	if err != nil {
		log.Error("completeJob(userId, job.Id)", err)
		return
	}
	if suc {
		var evt = events.UserOpenTreasureChest{}
		err = json.Unmarshal([]byte(job.Param), &evt)
		if err != nil {
			log.Error("json.Unmarshal([]byte(job.Param), &evt)", err)
			return
		}

		pointsEvt := interfaces.S.PointsScheduler.BuildPointsEventTypeReward(userId, job.Id, evt.Points, evt.TreasureChestName)
		err := emitter.Emit(pointsEvt)
		if err != nil {
			log.Error("emitter.Emit(pointsEvt)", err)
			return
		}
	}
}
