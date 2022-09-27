package report

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/db"
)

type PersistenceImpl struct {
	Table string
}

func (p *PersistenceImpl) Create(report *entities.FeedReport) error {
	_, err := db.Table(p.Table).Insert(report)
	return err
}

func (p *PersistenceImpl) List() ([]*entities.FeedReport, error) {
	var res []*entities.FeedReport
	err := db.Table(p.Table).Find(&res)
	return res, err
}

func (p *PersistenceImpl) GetById(id int64) (*entities.FeedReport, bool, error) {
	var report *entities.FeedReport
	exists, err := db.Table(p.Table).Where("id = ?", id).Get(&report)
	return report, exists, err
}

func (p *PersistenceImpl) Update(report *entities.FeedReport) error {
	_, err := db.Table(p.Table).ID(report.Id).Update(&report)
	return err
}

func (p *PersistenceImpl) Delete(id int64) error {
	_, err := db.Table(p.Table).Delete(entities.FeedReport{Id: id})
	return err
}

func NewPersistenceImpl(table string) IPersistence {
	return &PersistenceImpl{Table: table}
}
