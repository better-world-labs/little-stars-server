package essay

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/db"
	"errors"
	"github.com/go-xorm/xorm"
	"time"
)

type Service struct {
}

//go:inject-component
func NewService() *Service {
	return &Service{}
}

func (s *Service) Create(essay *entities.Essay) error {
	essay.CreateAt = time.Now()

	return db.Begin(func(session *xorm.Session) error {
		last, err := s.getLast(session)
		if err != nil {
			return err
		}

		essay.Sort = last.Sort + 1
		_, err = db.Table("essay").Insert(essay)
		return err
	})
}

func (s *Service) List() (slice []*entities.Essay, err error) {
	err = db.Table("essay").Asc("sort").Find(&slice)
	return
}

func (s *Service) ListLimit(limit int) (slice []*entities.Essay, err error) {
	err = db.Table("essay").
		Asc("sort").
		Limit(limit, 0).
		Find(&slice)
	return
}

func (s *Service) GetById(id int64) (*entities.Essay, error) {
	var essay entities.Essay
	_, err := db.Table("essay").Where("id = ?", id).Get(&essay)
	return &essay, err
}

func (s *Service) Delete(id int64) (err error) {
	_, err = db.Table("essay").Exec("delete from essay where id = ?", id)
	return
}

func (s *Service) Update(essay *entities.Essay) error {
	_, err := db.Table("essay").ID(essay.ID).AllCols().Update(essay)
	return err
}

func (s *Service) Sort(orderList []int64) error {
	return db.Begin(func(session *xorm.Session) error {
		count, err := db.Table("essay").Count()
		if err != nil {
			return err
		}

		if count != int64(len(orderList)) {
			return errors.New("number of orderList not match for essays")
		}

		for sort, i := range orderList {
			_, err := db.Exec("update essay set sort = ? where id = ?", sort, i)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Service) getLast(session *xorm.Session) (*entities.Essay, error) {
	var e entities.Essay
	_, err := session.Table("essay").Desc("sort").ForUpdate().Get(&e)
	return &e, err
}
