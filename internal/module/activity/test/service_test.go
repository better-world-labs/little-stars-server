package test

import (
	"aed-api-server/internal/module/activity"
	"aed-api-server/internal/pkg/async"
	"aed-api-server/internal/pkg/global"
	"aed-api-server/internal/pkg/record"
	activity2 "aed-api-server/internal/service/activity"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func init() {
	//db.InitEngine("mysql", "root:1qaz.2wsx@tcp(116.62.220.222:3306)/aed?charset=utf8mb4")
}

func TestCreateActivity(t *testing.T) {
	service := activity2.GetService()
	i := int64(6)
	err := service.Create(&activity.Activity{
		HelpInfoID: 6,
		Class:      record.ClassVolunteerNotified,
		UserID:     &i,
		Created:    global.FormattedTime(time.Now()),
		Record: map[string]interface{}{
			"userCount":     5,
			"firstUserName": "张三",
		},
	})

	assert.Nil(t, err)
}

func TestCreateOrUpdateActivity(t *testing.T) {
	service := activity2.GetService()
	u := uuid.NewString()
	i := int64(6)
	a := activity.Activity{
		HelpInfoID: 6,
		Uuid:       u,
		Class:      record.ClassVolunteerNotified,
		UserID:     &i,
		Created:    global.FormattedTime(time.Now()),
		Record: map[string]interface{}{
			"userCount":     5,
			"firstUserName": "张三",
		},
	}
	err := service.CreateOrUpdateByUUID(&a)

	i2 := int64(2)
	a.UserID = &i2
	err = service.CreateOrUpdateByUUID(&a)
	assert.Nil(t, err)

	r, err := service.GetOneByID(a.ID)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), *r.UserID)
}
func TestUpdateActivity(t *testing.T) {
	//service := activity.GetService()
	//o, err := service.GetOneByUUID("19bd7099-4d75-472e-bfc4-69acff714232")
	//assert.Nil(t, err)
	//if o != nil {
	//	p := float64(4)
	//	o.Points.Points = &p
	//	err := service.Update(o)
	//	assert.Nil(t, err)
	//
	//	o2, err := service.GetOneByUUID("19bd7099-4d75-472e-bfc4-69acff714232")
	//	assert.Equal(t, *o2.Points.Points, p)
	//}
}

func TestGetActivityById(t *testing.T) {
	service := activity2.GetService()
	o, err := service.GetOneByID(6)
	assert.Nil(t, err)
	fmt.Printf("%v", o)
}

func TestListByAID(t *testing.T) {
	service := activity2.GetService()
	o, err := service.ListByAID(105, 0)
	assert.Nil(t, err)
	fmt.Printf("%v", o)
}

func TestListByAIDs(t *testing.T) {
	service := activity2.GetService()
	o, err := service.ListByAIDs([]int64{105, 106, 107})
	for _, v := range o {
		for _, e := range v {
			fmt.Printf("%s", time.Time(e.Activity.Created))
		}
	}
	assert.Nil(t, err)
	fmt.Printf("%v+", o)
}

func TestListLatest(t *testing.T) {
	service := activity2.GetService()
	o, err := service.ListLatestCategorySorted(6, 3)
	assert.Nil(t, err)
	str, err := json.Marshal(o)
	assert.Nil(t, err)
	fmt.Printf("%s", string(str))
}

func TestListMultiLatest(t *testing.T) {
	service := activity2.GetService()
	o, err := service.ListMultiLatestCategorySorted([]int64{108, 106}, 3)
	assert.Nil(t, err)
	str, err := json.Marshal(o)
	assert.Nil(t, err)
	fmt.Printf("%s", string(str))
}

func TestListLatestAsync(t *testing.T) {
	async.TaskPool = async.NewPool(async.Config{
		GoRoutinePoolSize: 1000000,
		MaxTaskQueueSize:  10000000,
	})

	async.TaskPool.Start()

	service := activity2.GetService()
	_ = service.ListLatestCategorySortedAsync(105, 3)
}

func TestCreatePoint(t *testing.T) {
	service := activity2.GetService()
	point := float64(101)
	err := service.CreatePoint(&activity.Point{
		ActivityUuid: "19707825-1ca8-483a-be16-185b92bf89a9",
		Point:        &point,
	})
	assert.Nil(t, err)
}
