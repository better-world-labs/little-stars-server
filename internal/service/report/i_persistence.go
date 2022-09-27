package report

import "aed-api-server/internal/interfaces/entities"

type IPersistence interface {
	Create(report *entities.FeedReport) error

	List() ([]*entities.FeedReport, error)

	GetById(id int64) (*entities.FeedReport, bool, error)

	Update(report *entities.FeedReport) error

	Delete(id int64) error
}
