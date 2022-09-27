package entities

import (
	"aed-api-server/internal/pkg/global"
	"aed-api-server/internal/pkg/response"
	"encoding/json"
	"time"
)

const (
	OrderStatusToBeVerified = 0
	OrderStatusVerified     = 1
	OrderStatusExpired      = 2
)

const (
	CommodityStatusNotReleased CommodityStatus = 0
	CommodityStatusReleased    CommodityStatus = 1
)

type (
	CommodityStatus int

	BaseCommodity struct {
		Id         int64    `json:"id"`
		Name       string   `json:"name" binding:"required"`
		Sn         string   `json:"sn"`
		Price      int      `json:"price" binding:"required"`
		FrontCover string   `json:"frontCover" binding:"required"`
		Images     []string `json:"images" binding:"required"`
		Stock      int      `json:"stock" binding:"required"`
	}

	Commodity struct {
		BaseCommodity `xorm:"extends"`

		CreatedAt global.FormattedTime `json:"createdAt"`
		Status    CommodityStatus      `json:"status"`
	}

	Order struct {
		Id          int64
		Sn          string
		UserId      int64
		CommodityId int64
		Cost        int
		CreatedAt   time.Time
		ExpiresAt   time.Time
		VerifyAt    *time.Time
		VerifyCode  string
		Snapshot    *BaseCommodity
	}
)

func (c *Commodity) SubStock() error {
	if c.Status != CommodityStatusReleased {
		return response.ErrorCommodityNotReleased
	}

	if c.Stock > 0 {
		c.Stock--
		return nil
	}

	return response.ErrorCommodityStockIsNotEnough
}

func (e *BaseCommodity) FromDB(b []byte) error {
	return json.Unmarshal(b, e)
}

func (e *BaseCommodity) ToDB() ([]byte, error) {
	return json.Marshal(e)
}

func (o *Order) Expired() bool {
	return o.Status() == OrderStatusExpired
}

func (o *Order) Status() int {
	if o.VerifyAt != nil {
		return OrderStatusVerified
	}

	if o.ExpiresAt.Before(time.Now()) {
		return OrderStatusExpired
	}

	return OrderStatusToBeVerified
}
