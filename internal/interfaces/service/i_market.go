package service

import "aed-api-server/internal/interfaces/entities"

type MarketService interface {
	CreateCommodity(commodity entities.Commodity) error
	CommodityStandBy(id int64) error
	ListCommodities() ([]*entities.Commodity, error)
	ListCommoditiesByStatus(size int, status entities.CommodityStatus) ([]*entities.Commodity, error)
	GetCommodityById(id int64) (*entities.Commodity, bool, error)
	Buy(commodityId, userId int64) (*entities.Order, error)

	ListOrders(userId int64) ([]*entities.Order, error)
	ListToBeVerifiedOrders() ([]*entities.Order, error)
	GetOrderById(id int64) (*entities.Order, bool, error)
	GetOrderByVerifyCode(verifyCode string) (*entities.Order, bool, error)
	VerifyOrder(verifyCode string) error
}
