package evidence

import (
	"aed-api-server/internal/interfaces/entities"
	"github.com/go-xorm/xorm"
)

type IPersistence interface {
	CreateOrUpdateWithSession(session *xorm.Session, evi entities.Evidence) error
	GetByUUID(session *xorm.Session, uuid string) (*entities.Evidence, bool, error)
	GetOneByBusinessKeyWithSession(session *xorm.Session, category entities.EvidenceCategory, businessKey string) (*entities.Evidence, bool, error)
}
