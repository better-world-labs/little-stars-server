package controller

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/global"
	"errors"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
	"strconv"
	"time"
)

type MarketController struct {
}

func (c MarketController) ListCommoditiesAdmin(ctx *gin.Context) (interface{}, error) {
	commodities, err := interfaces.S.Market.ListCommodities()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"commodities": commodities,
	}, nil
}

func (c MarketController) ListReleasedCommodities(ctx *gin.Context) (interface{}, error) {
	commodities, err := interfaces.S.Market.ListCommoditiesByStatus(entities.CommodityStatusReleased)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"commodities": commodities,
	}, nil
}

func (c MarketController) GetCommodityById(ctx *gin.Context) (interface{}, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	commodity, exists, err := interfaces.S.Market.GetCommodityById(id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("not found")
	}

	return commodity, nil
}

func (c MarketController) Exchange(ctx *gin.Context) (interface{}, error) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	order, err := interfaces.S.Market.Buy(id, userId)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"orderId": order.Id,
	}, nil
}

func (c MarketController) MyOrders(ctx *gin.Context) (interface{}, error) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)
	orders, err := interfaces.S.Market.ListOrders(userId)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"orders": ParseOrderList(orders),
	}, nil
}

func (c MarketController) GetOrderById(ctx *gin.Context) (interface{}, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	order, exists, err := interfaces.S.Market.GetOrderById(id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, err
	}

	return NewOrderDto(order), nil
}

func (c MarketController) AdminVerifyOrder(ctx *gin.Context) (interface{}, error) {
	verifyCode := ctx.Param("verifyCode")
	err := interfaces.S.Market.VerifyOrder(verifyCode)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c MarketController) AdminToBeVerifiedOrders(ctx *gin.Context) (interface{}, error) {
	orders, err := interfaces.S.Market.ListToBeVerifiedOrders()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"orders": orders,
	}, nil
}

func (c MarketController) AdminCreateCommodity(ctx *gin.Context) (interface{}, error) {
	var commodity entities.BaseCommodity
	err := ctx.ShouldBindJSON(&commodity)
	if err != nil {
		return nil, err
	}

	err = interfaces.S.Market.CreateCommodity(entities.Commodity{BaseCommodity: commodity,
		CreatedAt: global.FormattedTime(time.Now()),
		Status:    entities.CommodityStatusNotReleased})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c MarketController) AdminCommodityStandby(ctx *gin.Context) (interface{}, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	err = interfaces.S.Market.CommodityStandBy(id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c MarketController) AdminGetOrderByVerifyCode(ctx *gin.Context) (interface{}, error) {
	verifyCode := ctx.Param("verifyCode")

	order, exists, err := interfaces.S.Market.GetOrderByVerifyCode(verifyCode)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("not found")
	}

	return NewOrderDto(order), nil
}

func NewMarketController() *MarketController {
	return &MarketController{}
}

func (c MarketController) MountAuthRouter(r *route.Router) {
	g := r.Group("/market")
	g.GET("/commodities/:id", c.GetCommodityById)
	g.POST("/exchange/:id", c.Exchange)
	g.GET("/my-orders", c.MyOrders)
	g.GET("/orders/:id", c.GetOrderById)
}

func (c MarketController) MountNoAuthRouter(r *route.Router) {
	r.Group("/market").
		GET("/commodities", c.ListReleasedCommodities)
}

func (c MarketController) MountAdminRouter(r *route.Router) {
	g := r.Group("/market")
	g.POST("/commodities", c.AdminCreateCommodity)
	g.PUT("/commodities-standby/:id", c.AdminCommodityStandby)
	g.GET("/commodities", c.ListCommoditiesAdmin)
	g.GET("/commodities/:id", c.GetCommodityById)
	g.POST("/orders/:verifyCode/verification", c.AdminVerifyOrder)
	g.GET("/orders-by-verifycode/:verifyCode", c.AdminGetOrderByVerifyCode)
}

type (
	OrderDto struct {
		Id          int64                   `json:"id"`
		Sn          string                  `json:"sn"`
		UserId      int64                   `json:"userId"`
		CommodityId int64                   `json:"commodityId"`
		Cost        int                     `json:"cost"`
		CreatedAt   global.FormattedTime    `json:"createdAt"`
		ExpiresAt   global.FormattedTime    `json:"expiresAt"`
		VerifyAt    *global.FormattedTime   `json:"verifyAt"`
		VerifyCode  string                  `json:"verifyCode"`
		Snapshot    *entities.BaseCommodity `json:"snapshot"`
		Status      int                     `json:"status"`
	}
)

func NewOrderDto(order *entities.Order) *OrderDto {
	dto := OrderDto{
		Id:          order.Id,
		Sn:          order.Sn,
		CommodityId: order.CommodityId,
		Cost:        order.Cost,
		CreatedAt:   global.FormattedTime(order.CreatedAt),
		ExpiresAt:   global.FormattedTime(order.ExpiresAt),
		VerifyCode:  order.VerifyCode,
		Snapshot:    order.Snapshot,
		Status:      order.Status(),
	}

	if order.VerifyAt != nil {
		verifyAt := global.FormattedTime(*order.VerifyAt)
		dto.VerifyAt = &verifyAt
	}

	return &dto
}

func ParseOrderList(order []*entities.Order) []*OrderDto {
	var dtos []*OrderDto

	for _, o := range order {
		dtos = append(dtos, NewOrderDto(o))
	}

	return dtos
}
