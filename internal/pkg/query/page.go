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

func (q *Query) GetLimit() (int, int) {
	return (q.Page - 1) * q.Size, q.Size
}

type Result[T any] struct {
	Total int `json:"total"`
	List  []T `json:"list"`
}

func NewResult[T any](list []T, total int) Result[T] {
	if list == nil {
		list = make([]T, 0)
	}

	return Result[T]{List: list, Total: total}
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
