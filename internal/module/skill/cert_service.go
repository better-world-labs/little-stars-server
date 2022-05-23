package skill

import (
	service2 "aed-api-server/internal/interfaces/service"
	"strconv"
)

type certService struct {
	service Service
}

func NewCertService(service Service) service2.CertService {
	return &certService{service: service}
}

func (c certService) CreateCert(projectId, userId int64) (string, string, error) {
	certImg, certNum, err := c.service.GenCert(userId, strconv.FormatInt(projectId, 10))
	return certImg, certNum, err
}
