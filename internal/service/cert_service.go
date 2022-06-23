package service

import (
	service2 "aed-api-server/internal/interfaces/service"
	"strconv"
)

type certService struct {
	Service service2.SkillService `inject:"-"`
}

func NewCertService() service2.CertService {
	return &certService{}
}

func (c certService) CreateCert(projectId, userId int64) (string, string, error) {
	certImg, certNum, err := c.Service.GenCert(userId, strconv.FormatInt(projectId, 10))
	return certImg, certNum, err
}
