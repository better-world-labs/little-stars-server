package test

import (
	"aed-api-server/internal/module/aid"
	"aed-api-server/internal/module/user"
	"aed-api-server/internal/pkg/db"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var svc = aid.NewService(user.NewService(nil))

func init() {
	db.InitEngine("mysql", "root:1qaz.2wsx@tcp(116.62.220.222:3306)/aed?charset=utf8mb4")
}

func TestPublishHelpInfo(t *testing.T) {
	_, err := svc.PublishHelpInfo(2, &aid.PublishDTO{
		Address:       "成都市",
		DetailAddress: "local detail",
		Longitude:     20.12364,
		Latitude:      4453.23434,
		Images: []*aid.ImageDTO{
			{
				Thumbnail: "http://xxx/xxx/thum",
				Origin:    "http://xx/xx/origin",
			},
		},
	})
	assert.Nil(t, err)
}

func TestGetHelpImageByIDs(t *testing.T) {
	r, err := svc.GetHelpImagesByHelpInfoIDs([]int64{4, 5, 3, 2})
	assert.Nil(t, err)
	fmt.Println(r)
}
