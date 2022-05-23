package task

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

func mkTaskData() {
	db.Exec(`INSERT INTO 
			task(id,name,image,description,url,user_range_type,aed_range_type,time_limit,point,theme_id,created_at,is_picket,picket_condition,level) 
			VALUES (?,'纠察任务','','未经过纠察（限时7天，2倍积分）',NULL,0,0,604800,20.0,1,'2022-02-14 17:09:12',1,1,1);`, testTaskId)
}

func cleanTaskData() {
	db.Exec(`delete from task where id = ?`, testTaskId)
}

func mkUserRangeData() {
	db.Exec(`INSERT INTO task_range_user ( task_id, user_id) VALUES ( ?, ?);`, testTaskId, testUserId)
}

func cleanUserRangeData() {
	db.Exec(`delete from task_range_user where task_id = ?`, testTaskId)
}

func mkAedRangeData() {
	db.Exec(`INSERT INTO task_range_aed ( task_id, device_id) VALUES ( ?, ?);`, testTaskId, testDeviceId)
}

func cleanAedRangeData() {
	db.Exec(`delete from task_range_aed where task_id = ?`, testTaskId)
}

func Test_userInRange(t *testing.T) {
	t.Cleanup(InitDbAndConfig())
	defer cleanUserRangeData()
	defer cleanTaskData()

	inRange := userInRange(testTaskId, testUserId)
	assert.False(t, inRange)

	mkTaskData()
	mkUserRangeData()
	inRange = userInRange(testTaskId, testUserId)
	assert.True(t, inRange)
}

func Test_deviceInRange(t *testing.T) {
	t.Cleanup(InitDbAndConfig())
	defer cleanAedRangeData()
	defer cleanTaskData()

	inRange := deviceInRange(testTaskId, testDeviceId)
	assert.False(t, inRange)

	mkTaskData()
	mkAedRangeData()
	inRange = deviceInRange(testTaskId, testDeviceId)
	assert.True(t, inRange)
}

func Test_findPicketTasks(t *testing.T) {
	t.Cleanup(InitDbAndConfig())
	defer cleanTaskData()

	mkTaskData()
	tasks, err := findPicketTasks()
	assert.Nil(t, err)
	matched := false
	for _, task := range tasks {
		if task.Id == testTaskId {
			matched = true
		}
	}
	assert.True(t, matched)
}

func Test_genJob(t *testing.T) {
	mockPicketCondition()
	t.Cleanup(InitDbAndConfig())
	defer cleanTaskData()
	mkTaskData()

	task := findTaskById(testTaskId)
	assert.NotNil(t, task)

	job := task.genJob(testUserId, &entities.Device{
		Id: testDeviceId,
	})

	assert.NotNil(t, job)
}
