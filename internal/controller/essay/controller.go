package essay

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	service2 "aed-api-server/internal/interfaces/service"
	"errors"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type EssayList struct {
	entities.Essay

	CreateBy *entities.SimpleUser `json:"createBy"`
}

type Dto struct {
	Title       string   `json:"title"      binding:"required"`
	Description string   `json:"description"`
	Type        int      `json:"type"       binding:"required,min=1,max=4"`
	FrontCover  []string `json:"frontCover" binding:"required"`
	Content     string   `json:"content"    binding:"required"`
	Extra       string   `json:"extra"`
}

type Controller struct {
}

func (c *Controller) MountNoAuthRouter(r *route.Router) {
	essayGroup := r.Group("/essays")
	essayGroup.GET("", c.List)
	essayGroup.GET("/:id", c.GetOne)
}

func (c *Controller) MountAdminRouter(r *route.Router) {
	essayGroup := r.Group("/essays")
	essayGroup.POST("", c.Create)
	essayGroup.DELETE("/:id", c.Delete)
	essayGroup.GET("/:id", c.GetOne)
	essayGroup.PUT("/:id", c.Update)
	r.POST("/essays-sorts", c.Sort)
	essayGroup.GET("", c.List)
}

func NewEssayController() *Controller {
	return &Controller{}
}

func service() service2.EssayService {
	return interfaces.S.Essay
}

func (c *Controller) Create(ctx *gin.Context) (interface{}, error) {
	var param Dto
	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, err
	}

	err = service().Create(&entities.Essay{
		Title:       param.Title,
		Description: param.Description,
		Type:        param.Type,
		Content:     param.Content,
		FrontCover:  param.FrontCover,
		Extra:       param.Extra,
		CreateAt:    time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *Controller) List(ctx *gin.Context) (interface{}, error) {
	list, err := service().List()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"essays": parseEssayList(list)}, nil
}

func (c *Controller) Delete(ctx *gin.Context) (interface{}, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	err = service().Delete(id)
	if err != nil {
		return nil, err
	}

	return nil, nil
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

func (c *Controller) GetOne(ctx *gin.Context) (interface{}, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	essay, err := service().GetById(id)
	if err != nil {
		return nil, err
	}

	list := parseEssayList([]*entities.Essay{essay})
	if len(list) == 0 {
		return nil, errors.New("not found")
	}

	return list[0], nil
}

func (c *Controller) Update(ctx *gin.Context) (interface{}, error) {
	var param Dto
	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, err
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	err = service().Update(&entities.Essay{
		ID:          id,
		Title:       param.Title,
		Description: param.Description,
		Type:        param.Type,
		Content:     param.Content,
		FrontCover:  param.FrontCover,
		Extra:       param.Extra,
		CreateAt:    time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *Controller) Sort(ctx *gin.Context) (interface{}, error) {
	var param struct {
		OrderList []int64 `json:"orderList"`
	}

	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, err
	}

	err = service().Sort(param.OrderList)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
