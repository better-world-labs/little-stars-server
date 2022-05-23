package task

import (
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/db"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func mkData() {
	mkTaskData()
	//有时间限制的 且 未结束 且 未完成 且 已读
	saveJob(&service.Job{
		UserId:      testUserId,
		TaskId:      testTaskId,
		IsRead:      True,
		Status:      JOB_STATUS_INIT,
		CreatedAt:   time.Now(),
		IsTimeLimit: True,
		BeginLimit:  time.Now(),
		EndLimit:    time.Now().Add(_1d),
		DeviceId:    testDeviceId,
	})

	//有时间限制的 且 结束 且 未完成 且 已读
	saveJob(&service.Job{
		UserId:      testUserId,
		TaskId:      testTaskId,
		IsRead:      True,
		Status:      JOB_STATUS_INIT,
		CreatedAt:   time.Now(),
		IsTimeLimit: True,
		BeginLimit:  time.Now(),
		EndLimit:    time.Now().Add(-_1d),
		DeviceId:    testDeviceId,
	})

	//有时间限制的 且 结束 且 完成 且 已读
	saveJob(&service.Job{
		UserId:      testUserId,
		TaskId:      testTaskId,
		IsRead:      True,
		Status:      JOB_STATUS_COMPLETE,
		CreatedAt:   time.Now(),
		IsTimeLimit: True,
		BeginLimit:  time.Now(),
		EndLimit:    time.Now().Add(-_1d),
	})

	//有时间限制的 且 未结束 且 未完成 且 未读
	saveJob(&service.Job{
		UserId:      testUserId,
		TaskId:      testTaskId,
		IsRead:      False,
		Status:      JOB_STATUS_INIT,
		CreatedAt:   time.Now(),
		IsTimeLimit: True,
		BeginLimit:  time.Now(),
		EndLimit:    time.Now().Add(_1d),
	})

	//有时间限制的 且 未结束 且 完成 且 未读
	saveJob(&service.Job{
		UserId:      testUserId,
		TaskId:      testTaskId,
		IsRead:      False,
		Status:      JOB_STATUS_COMPLETE,
		CreatedAt:   time.Now(),
		IsTimeLimit: True,
		BeginLimit:  time.Now(),
		EndLimit:    time.Now().Add(_1d),
	})
}

func cleanData() {
	db.Exec(`delete from task_job where user_id = ?`, testUserId)
	cleanTaskData()
}

func reset() func() {
	mkData()
	return cleanData
}

func Test_Job(t *testing.T) {
	t.Cleanup(InitDbAndConfig())
	t.Run("saveJob", func(t *testing.T) {
		t.Cleanup(reset())
		err := saveJob(&service.Job{
			UserId:      testUserId,
			TaskId:      testTaskId,
			IsRead:      True,
			Status:      JOB_STATUS_INIT,
			CreatedAt:   time.Now(),
			IsTimeLimit: True,
			BeginLimit:  time.Now(),
			EndLimit:    time.Now().Add(_1d),
			DeviceId:    testDeviceId,
		})
		assert.Nil(t, err)

	})

	t.Run("getUserJobs", func(t *testing.T) {
		//reset()
		t.Cleanup(reset())
		jobs, err := getUserJobs(testUserId, pageSize, []int{JOB_STATUS_INIT}, DEFAULT_JOB_ID, false)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(jobs))

		jobs, err = getUserJobs(testUserId, pageSize, []int{JOB_STATUS_INIT, JOB_STATUS_COMPLETE}, DEFAULT_JOB_ID, false)
		assert.Nil(t, err)
		assert.Equal(t, 3, len(jobs))

		jobs, err = getUserJobs(testUserId, pageSize, []int{JOB_STATUS_INIT, JOB_STATUS_COMPLETE}, DEFAULT_JOB_ID, true)
		assert.Nil(t, err)
		assert.Equal(t, 5, len(jobs))
	})

	t.Run("findUserTaskStat & readUserTask", func(t *testing.T) {
		t.Cleanup(reset())
		stat, _ := findUserTaskStat(testUserId)
		assert.Equal(t, 2, stat.Unread)
		jobs, err := getUserJobs(testUserId, pageSize, []int{JOB_STATUS_INIT}, DEFAULT_JOB_ID, false)
		assert.Nil(t, err)

		readUserTask(jobs[0].Id, testUserId)
		stat, _ = findUserTaskStat(testUserId)
		assert.Equal(t, 1, stat.Unread)
	})

	t.Run("userHasUnCompletePicketJob", func(t *testing.T) {
		t.Cleanup(reset())
		existed := userHasUnCompletePicketJob(testUserId, testDeviceId)
		assert.True(t, existed)

		job, err := findUserTaskByUserIdAndDeviceId(testUserId, testDeviceId)
		assert.Nil(t, err)

		err = completeJob(testUserId, job.Id)
		assert.Nil(t, err)

		existed = userHasUnCompletePicketJob(testUserId, testDeviceId)
		assert.False(t, existed)
	})

	t.Run("findUserTaskByUserIdAndDeviceId", func(t *testing.T) {
		t.Cleanup(reset())
		job, err := findUserTaskByUserIdAndDeviceId(testUserId, testDeviceId)
		assert.Nil(t, err)
		assert.Equal(t, testDeviceId, job.DeviceId)
		assert.True(t, job.TimeLimit > 0)
	})
}
