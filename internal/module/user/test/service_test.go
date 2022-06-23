package test

import (
	//"aed-api-server/internal/module/user"
	//"aed-api-server/internal/pkg/location"
	//"aed-api-server/internal/service"
	//"github.com/stretchr/testify/assert"
	"testing"
)

//var svc = service.NewService(&mockWechatClient{})

func init() {
	//db.InitEngine("mysql", "root:1qaz.2wsx@tcp(116.62.220.222:3306)/aed?charset=utf8mb4")
}

func TestCreateOrUpdatePosition(t *testing.T) {
	//err := svc.UpdatePosition(&user.Position{AccountID: 4, Coordinate: &location.Coordinate{Longitude: 23.9845684994126, Latitude: 89.5534423}})
	//assert.Nil(t, err)
}

func TestListUserPositionByUserIDs(t *testing.T) {
	//res, err := svc.ListPositionByUserIDs([]int64{1, 2, 3, 4})
	//assert.Nil(t, err)
	//assert.NotNil(t, res)
}
func TestGetUserPositionByUserID(t *testing.T) {
	//res, err := svc.GetPositionByUserID(2)
	//assert.Nil(t, err)
	//assert.NotNil(t, res)
}
func TestUpdateMobile(t *testing.T) {
	//err := svc.UpdateMobile(1, "184xxxxxxxx")
	//assert.Nil(t, err)
}

func TestWechatMiniProgramLogin(t *testing.T) {
	//account, token, sessionKey, err := svc.WechatMiniProgramLogin("ccc", "", "nDeb5/FIiRynRIyTI2cx/A==")
	//
	//assert.Nil(t, err)
	//assert.NotNil(t, account)
	//assert.NotEqual(t, "", token)
	//assert.NotEqual(t, "", sessionKey)
	//assert.NotEqual(t, "", account.Mobile)
	//assert.NotEqual(t, "", account.ID)
	//assert.NotEqual(t, "", account.Openid)
}

func TestUpdateUserInfo(t *testing.T) {
	//err := svc.UpdateUserInfo(&user.User{ID: 6, Avatar: "http://xxx"})
	//assert.Nil(t, err)
}
