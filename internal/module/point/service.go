package point

import (
	"aed-api-server/internal/pkg/db"
)

type service struct {
}

func NewService() Service {
	return service{}
}

func (s service) Detail(accountID int64) ([]*Point, error) {
	sess := db.GetSession()
	defer sess.Close()

	list := make([]*Point, 0)
	err := sess.Where("account_id = ?", accountID).OrderBy("id desc").Find(&list)

	return list, err
}

func (s service) TotalPoints(accountID int64) (float64, error) {
	sess := db.GetSession()
	defer sess.Close()

	table := new(Point)
	total, err := sess.Where("account_id = ?", accountID).Sum(table, "points")

	return float64(total), err
}

func (s service) AddPoint(point Point) error {
	sess := db.GetSession()
	defer sess.Close()

	_, err := sess.Insert(point)
	return err
}
