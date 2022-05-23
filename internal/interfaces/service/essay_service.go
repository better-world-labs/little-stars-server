package service

import (
	"aed-api-server/internal/interfaces/entities"
)

type (
	EssayService interface {
		Create(essay *entities.Essay) error

		List() ([]*entities.Essay, error)

		GetById(id int64) (essay *entities.Essay, err error)

		Delete(id int64) error

		Update(essay *entities.Essay) error

		Sort(orderList []int64) error
	}
)
