package controller

import (
	"aed-api-server/internal/module/achievement"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (con Controller) ListAllMedalMeta(c *gin.Context) {
	list, err := achievement.ListMedals()
	utils.MustNil(err, err)
	response.ReplyOK(c, map[string]interface{}{"medals": list})
}

func (con Controller) ListUsersMedal(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	list, err := achievement.ListUsersMedal(accountID)
	utils.MustNil(err, err)

	response.ReplyOK(c, map[string]interface{}{"medals": list})
}

func (con Controller) CreateUsersMedalEvidences(c *gin.Context) {
	//um, err := ListAll()
	//utils.MustNil(err, err)
	//
	//for _, u := range um {
	//	_, exist, err := con.evidenceService.GetEvidenceByBusinessKey(strconv.FormatInt(u.ID, 10), entities.EvidenceCategoryMedal)
	//	if err != nil {
	//		response.ReplyError(c, err)
	//		return
	//	}
	//
	//	if exist {
	//		continue
	//	}
	//
	//	account, err := con.userService.GetUserByID(u.UserID)
	//	if err != nil {
	//		response.ReplyError(c, err)
	//		return
	//	}
	//
	//	medal, exists, err := GetById(u.MedalID)
	//	if err != nil {
	//		response.ReplyError(c, err)
	//		return
	//	}
	//
	//	if !exists {
	//		response.ReplyError(c, errors.New("medal not found"))
	//		return
	//	}
	//
	//	errChan := evidenceService.CreateEvidenceAsync(&claim.Medal{
	//		Mobile: account.Mobile,
	//		Medal:  "见义勇为",
	//	}, "见义勇为", account.ID, entities.EvidenceCategoryMedal, strconv.FormatInt(u.ID, 10))
	//
	//	if err := <-errChan; err != nil {
	//		if err != nil {
	//			response.ReplyError(c, err)
	//			return
	//		}
	//	}
	//}

	response.ReplyOK(c, nil)
}

func (con Controller) ListUsersMedalToast(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	list, err := achievement.ListUsersMedalToast(accountID)
	utils.MustNil(err, err)

	response.ReplyOK(c, map[string]interface{}{"medals": list})
}
