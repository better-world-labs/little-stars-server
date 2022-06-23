package controller

import (
	"aed-api-server/internal/module/achievement"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type AchievementController struct {
}

func NewAchievementController() *AchievementController {
	return &AchievementController{}
}

func (con AchievementController) MountAuthRouter(r *route.Router) {
	// achievement
	r.GET("/achievement/medals", con.ListAllMedalMeta)
	r.GET("/achievement/user-medals", con.ListUsersMedal)
	r.GET("/achievement/medal/toast", con.ListUsersMedalToast)
}

func (con AchievementController) ListAllMedalMeta(c *gin.Context) (interface{}, error) {
	list, err := achievement.ListMedals()
	utils.MustNil(err, err)
	return map[string]interface{}{"medals": list}, nil
}

func (con AchievementController) ListUsersMedal(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	list, err := achievement.ListUsersMedal(accountID)
	utils.MustNil(err, err)

	return map[string]interface{}{"medals": list}, nil
}

func (con AchievementController) ListUsersMedalToast(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	list, err := achievement.ListUsersMedalToast(accountID)
	utils.MustNil(err, err)

	return map[string]interface{}{"medals": list}, nil
}
