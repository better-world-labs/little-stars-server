package controller

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
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

type EssayDto struct {
	Title       string   `json:"title"      binding:"required"`
	Description string   `json:"description"`
	Type        int      `json:"type"       binding:"required,min=1,max=4"`
	FrontCover  []string `json:"frontCover" binding:"required"`
	Content     string   `json:"content"    binding:"required"`
	Extra       string   `json:"extra"`
}

type EssayController struct {
	Essay service.EssayService `inject:"-"`
}

func (c *EssayController) MountNoAuthRouter(r *route.Router) {
	essayGroup := r.Group("/essays")
	essayGroup.GET("", c.List)
	essayGroup.GET("/:id", c.GetOne)
}

func (c *EssayController) MountAdminRouter(r *route.Router) {
	essayGroup := r.Group("/essays")
	essayGroup.POST("", c.Create)
	essayGroup.DELETE("/:id", c.Delete)
	essayGroup.GET("/:id", c.GetOne)
	essayGroup.PUT("/:id", c.Update)
	r.POST("/essays-sorts", c.Sort)
	essayGroup.GET("", c.List)
}

//go:inject-component
func NewEssayController() *EssayController {
	return &EssayController{}
}

func (c *EssayController) Create(ctx *gin.Context) (interface{}, error) {
	var param EssayDto
	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, err
	}

	err = c.Essay.Create(&entities.Essay{
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

func (c *EssayController) List(ctx *gin.Context) (interface{}, error) {
	essays := make([]*entities.Essay, 0)
	var err error

	if sizeQuery, exists := ctx.GetQuery("size"); exists {
		size, err := strconv.Atoi(sizeQuery)
		if err != nil {
			return nil, err
		}

		essays, err = c.Essay.ListLimit(size)
		if err != nil {
			return nil, err
		}
	} else {
		essays, err = c.Essay.List()
		if err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"essays": parseEssayList(essays)}, nil
}

func (c *EssayController) Delete(ctx *gin.Context) (interface{}, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	err = c.Essay.Delete(id)
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

func (c *EssayController) GetOne(ctx *gin.Context) (interface{}, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	essay, err := c.Essay.GetById(id)
	if err != nil {
		return nil, err
	}

	list := parseEssayList([]*entities.Essay{essay})
	if len(list) == 0 {
		return nil, errors.New("not found")
	}

	return list[0], nil
}

func (c *EssayController) Update(ctx *gin.Context) (interface{}, error) {
	var param EssayDto
	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, err
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	err = c.Essay.Update(&entities.Essay{
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

func (c *EssayController) Sort(ctx *gin.Context) (interface{}, error) {
	var param struct {
		OrderList []int64 `json:"orderList"`
	}

	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, err
	}

	err = c.Essay.Sort(param.OrderList)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
