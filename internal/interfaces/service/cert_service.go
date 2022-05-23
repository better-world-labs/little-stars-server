package service

type CertService interface {
	CreateCert(projectId, userId int64) (string, string, error)
}
