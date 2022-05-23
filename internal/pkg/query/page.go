package page

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Query struct {
	Page int `form:"page" json:"page,omitempty"`
	Size int `form:"size" json:"size,omitempty"`
}

const DefaultPageSize = "20"
const DefaultPage = "0"

func (q *Query) getLimit() (int, int) {
	return (q.Page - 1) * q.Size, q.Size
}

type Result struct {
	Total int         `json:"total"`
	List  interface{} `json:"list"`
}

func BindPageQuery(c *gin.Context) (*Query, error) {
	page := c.DefaultQuery("page", DefaultPage)
	size := c.DefaultQuery("size", DefaultPageSize)
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}

	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		return nil, err
	}

	return &Query{
		Page: pageInt,
		Size: sizeInt,
	}, nil
}
