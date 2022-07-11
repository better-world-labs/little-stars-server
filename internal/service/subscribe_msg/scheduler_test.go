package subscribe_msg

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg/config"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/service/point"
	"aed-api-server/internal/service/user"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"testing"
)

func Test_sendPointsExpiringMsg(t *testing.T) {
	c, err := config.LoadConfig("../../../")
	if err != nil {
		panic("get config error")
	}
	interfaces.InitConfig(c)
	interfaces.S.User = user.NewService()
	interfaces.S.Points = point.NewService()
	interfaces.S.UserConfig = user.NewUserConfigService()
	db.InitEngine(c.Database)
	log.Init(c.Log)

	sendPointsExpiringMsg()
}
