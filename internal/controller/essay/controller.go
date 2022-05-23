package essay

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	service2 "aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/response"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type EssayList struct {
	entities.Essay

	CreateBy *entities.SimpleUser `json:"createBy"`
}

type Dto struct {
	Title      string   `json:"title"      binding:"required"`
	Type       int      `json:"type"       binding:"required,min=1,max=3"`
	FrontCover []string `json:"frontCover" binding:"required"`
	Content    string   `json:"content"    binding:"required"`
	Extra      string   `json:"extra"`
}

type Controller struct {
}

func service() service2.EssayService {
	return interfaces.S.Essay
}

func (c *Controller) Create(ctx *gin.Context) {
	var param Dto
	err := ctx.BindJSON(&param)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	err = service().Create(&entities.Essay{
		Title:      param.Title,
		Type:       param.Type,
		Content:    param.Content,
		FrontCover: param.FrontCover,
		Extra:      param.Extra,
		CreateAt:   time.Now(),
	})
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	response.ReplyOK(ctx, nil)
}

func (c *Controller) List(ctx *gin.Context) {
	list, err := service().List()
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	response.ReplyOK(ctx, map[string]interface{}{
		"essays": parseEssayList(list)})
}

func (c *Controller) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	err = service().Delete(id)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	response.ReplyOK(ctx, nil)
}

func parseEssayList(essays []*entities.Essay) []*EssayList {
	list := make([]*EssayList, len(essays))

	for i, e := range essays {
		list[i] = &EssayList{
			Essay: *e,
			CreateBy: &entities.SimpleUser{
				Nickname: "小星星志愿者",
				Avatar:   "https://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLKxhW6g5dxI5vXULkpHEGRothHq9VzeuATmE0U2hLtwjMGaLK0e684HbyO9jxY0fwVyuc67Ed2ng/132",
			},
		}
	}

	return list
}

func (c *Controller) GetOne(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	essay, err := service().GetById(id)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	list := parseEssayList([]*entities.Essay{essay})
	if len(list) == 0 {
		response.ReplyError(ctx, errors.New("not found"))
		return
	}

	response.ReplyOK(ctx, list[0])
}

func (c *Controller) Update(ctx *gin.Context) {
	var param Dto
	err := ctx.BindJSON(&param)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	err = service().Update(&entities.Essay{
		ID:         id,
		Title:      param.Title,
		Type:       param.Type,
		Content:    param.Content,
		FrontCover: param.FrontCover,
		Extra:      param.Extra,
		CreateAt:   time.Now(),
	})
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	response.ReplyOK(ctx, nil)
}

func (c *Controller) Sort(ctx *gin.Context) {
	var param struct {
		OrderList []int64 `json:"orderList"`
	}

	err := ctx.BindJSON(&param)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	err = service().Sort(param.OrderList)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	response.ReplyOK(ctx, nil)
}
