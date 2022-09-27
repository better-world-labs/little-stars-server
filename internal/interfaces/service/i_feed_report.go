package service

import "aed-api-server/internal/interfaces/entities"

type IFeedReport interface {
	Create(report *entities.FeedReport) error
	List() ([]*entities.FeedReport, error)
	GetById(id int64) (*entities.FeedReport, bool, error)
	UpdateStatus(id int64, status entities.FeedReportStatus) error
}
