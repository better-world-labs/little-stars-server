package service

import (
	"aed-api-server/internal/domains"
	"aed-api-server/internal/interfaces/entities"
)

type ExamService interface {
	SetCertService(service CertService)
	Start(projectId int64, userId int64, examType int) (*domains.Exam, error)
	Save(examID, userId int64, paper map[int64][]int) error
	Submit(examID, userId int64, paper map[int64][]int) (string, string, []*entities.DealPointsEventRst, error)
	ListLatestSubmitted(projectId int64, userId int64, examType int, latest int) ([]*domains.Exam, error)
	GetLatestUnSubmitted(projectId int64, userId int64, examType int) (*domains.Exam, bool, error)
	GetByID(examID int64) (*domains.Exam, bool, error)
}
