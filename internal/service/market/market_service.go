package market

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/global"
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"time"
)

const (
	TableNameCommodity = "market_commodity"
	TableNameOrder     = "market_order"
)

var (
	OrderExpiresIn = 60 * 24 * time.Hour
)

type marketService struct {
}

//go:inject-component
func NewMarketService() *marketService {
	return &marketService{}
}

func (m *marketService) CreateCommodity(commodity entities.Commodity) error {
	commodity.Sn = m.generateSn()
	commodity.CreatedAt = global.FormattedTime(time.Now())
	commodity.Status = entities.CommodityStatusNotReleased
	_, err := db.Table(TableNameCommodity).Insert(&commodity)
	return err
}

func (m *marketService) ListCommoditiesByStatus(status int) ([]*entities.Commodity, error) {
	var commodities []*entities.Commodity
	err := db.Table(TableNameCommodity).Where("status = ?", status).Asc("sort").Find(&commodities)
	return commodities, err
}

func (m *marketService) ListCommodities() ([]*entities.Commodity, error) {
	var commodities []*entities.Commodity
	err := db.Table(TableNameCommodity).Asc("sort").Find(&commodities)
	return commodities, err
}

func (m *marketService) GetCommodityById(id int64) (*entities.Commodity, bool, error) {
	var commodity entities.Commodity
	exists, err := db.Table(TableNameCommodity).Where("id = ?", id).Get(&commodity)
	return &commodity, exists, err
}

func (m *marketService) GetCommodityByIdForUpdate(session *xorm.Session, id int64) (*entities.Commodity, bool, error) {
	var commodity entities.Commodity
	exists, err := session.Table(TableNameCommodity).Where("id = ?", id).ForUpdate().Get(&commodity)
	return &commodity, exists, err
}

func (m *marketService) Buy(commodityId, userId int64) (*entities.Order, error) {
	var order *entities.Order

	return order, db.Transaction(func(session *xorm.Session) error {
		commodity, exists, err := m.GetCommodityByIdForUpdate(session, commodityId)
		if err != nil {
			return err
		}

		if !exists {
			return errors.New("not found ")
		}

		err = commodity.SubStock()
		if err != nil {
			return err
		}

		err = m.setStock(session, commodityId, commodity.Stock)
		if err != nil {
			return err
		}

		err = interfaces.S.Points.AddPoint(userId, -commodity.Price, fmt.Sprintf("兑换“%s”", commodity.Name), entities.PointsEventTypeTransaction)
		if err != nil {
			return err
		}

		order, err = m.CreateOrder(session, commodity, userId)
		return err
	})
}

func (m *marketService) CreateOrder(session *xorm.Session, commodity *entities.Commodity, userId int64) (*entities.Order, error) {
	now := time.Now()
	order := &entities.Order{
		Sn:          m.generateSn(),
		UserId:      userId,
		CommodityId: commodity.Id,
		Cost:        commodity.Price,
		CreatedAt:   now,
		ExpiresAt:   now.Add(OrderExpiresIn),
		Snapshot:    &commodity.BaseCommodity,
	}

	_, err := session.Table(TableNameOrder).Insert(order)
	if err != nil {
		return nil, err
	}

	_, err = session.Exec(fmt.Sprintf("update %s set verify_code = ? where id = ?", TableNameOrder), m.generateVerifyCode(order.Id), order.Id)
	return order, err
}

func (m *marketService) ListOrders(userId int64) ([]*entities.Order, error) {
	var orders []*entities.Order
	err := db.Table(TableNameOrder).Where("user_id = ?", userId).Desc("id").Find(&orders)
	return orders, err
}

func (m *marketService) ListToBeVerifiedOrders() ([]*entities.Order, error) {
	var orders []*entities.Order
	err := db.Table(TableNameOrder).Where("expires_at > ? and verify_at is null", time.Now()).Find(&orders)
	return orders, err
}

func (m *marketService) GetOrderById(id int64) (*entities.Order, bool, error) {
	var order entities.Order
	exists, err := db.Table(TableNameOrder).Where("id = ?", id).Get(&order)
	return &order, exists, err
}

func (m *marketService) GetOrderByVerifyCode(verifyCode string) (*entities.Order, bool, error) {
	var order entities.Order
	exists, err := db.Table(TableNameOrder).Where("verify_code = ?", verifyCode).Get(&order)
	return &order, exists, err
}

func (m *marketService) setStock(session *xorm.Session, id int64, value int) error {
	_, err := session.Exec(fmt.Sprintf("update %s set stock = ? where id = ?", TableNameCommodity), value, id)
	return err
}

func (m *marketService) compareAndSetStock(session *xorm.Session, id int64, excepted, value int) (bool, error) {
	res, err := session.Exec(fmt.Sprintf("update %s set stock = ? where id = ? and stock = ?", TableNameCommodity), value, id, excepted)
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func (m *marketService) VerifyOrder(verifyCode string) error {
	order, exists, err := m.GetOrderByVerifyCode(verifyCode)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("not found")
	}

	if order.Expired() {
		return errors.New("the verify code is expired")
	}

	now := time.Now()
	order.VerifyAt = &now

	_, err = db.Table(TableNameOrder).ID(order.Id).Update(order)
	return err
}

func (m *marketService) generateVerifyCode(orderId int64) string {
	u := uuid.New()
	orderIdStr := strconv.FormatInt(orderId, 10)
	formatted := strconv.FormatUint(uint64(u.ID()), 10)
	formatted = orderIdStr + formatted
	l := len(formatted)
	if l < 10 {
		for i := 0; i < 10-i; i++ {
			formatted += strconv.Itoa(rand.Int())
		}
	}
	return fmt.Sprintf("%s", formatted[:10])
}

func (m *marketService) generateSn() string {
	u := uuid.New()
	formatted := strconv.FormatUint(uint64(u.ID()), 10)
	l := len(formatted)
	if l < 7 {
		for i := 0; i < 10-i; i++ {
			formatted += strconv.Itoa(rand.Int())
		}
	}
	now := time.Now().UnixMilli()
	return fmt.Sprintf("%d%s", now, formatted[:7])
}

func (m *marketService) CommodityStandBy(id int64) error {
	_, err := db.Table(TableNameCommodity).Exec(fmt.Sprintf("update %s set status = ? where id = ?", TableNameCommodity), entities.CommodityStatusReleased, id)
	return err
}
