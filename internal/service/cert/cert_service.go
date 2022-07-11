package cert

import (
	"aed-api-server/internal/interfaces/service"
	"strconv"
)

type certService struct {
	Service service.SkillService `inject:"-"`
}

//go:inject-component
func NewCertService() service.CertService {
	return &certService{}
}

func (c certService) CreateCert(projectId, userId int64) (string, string, error) {
	certImg, certNum, err := c.Service.GenCert(userId, strconv.FormatInt(projectId, 10))
	return certImg, certNum, err
}
