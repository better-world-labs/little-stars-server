package evidence

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/base"
	"github.com/go-xorm/xorm"
	"time"
)

type persistence struct {
}

func NewPersistence() IPersistence {
	return &persistence{}
}

func (s *persistence) CreateOrUpdateWithSession(session *xorm.Session, evi entities.Evidence) error {
	now := time.Now()
	_, err := session.Exec("insert into `evidence` (`account_id`,`category`,`business_key`,`evidence_id`, `content_hash`,`transaction_hash`,`credential_id`, `file_bytes`, `time`)"+
		" values(?, ?, ?, ?, ?, ?, ?, ?, ?) on duplicate key update `evidence_id`=?, `content_hash`=?, transaction_hash=?, credential_id=?, file_bytes=?, time=?",
		evi.AccountID, evi.Category, evi.BusinessKey, evi.EvidenceID, evi.ContentHash, evi.TransactionHash, evi.CredentialID, evi.FileBytes, now,
		evi.EvidenceID, evi.ContentHash, evi.TransactionHash, evi.CredentialID, evi.FileBytes, now)

	return err
}

func (s *persistence) GetOneByBusinessKeyWithSession(session *xorm.Session, category entities.EvidenceCategory, businessKey string) (*entities.Evidence, bool, error) {
	var evi entities.Evidence
	exists, err := session.Table("evidence").Where("business_key = ? and category = ?", businessKey, category).Get(&evi)
	if err != nil {
		return nil, false, base.WrapError("evidence.persistence", "get evidence error", err)
	}

	if !exists {
		return nil, false, nil
	}

	return &evi, true, nil
}

func (s *persistence) GetByUUID(session *xorm.Session, uuid string) (*entities.Evidence, bool, error) {
	var evi entities.Evidence
	exists, err := session.Table("evidence").Where("uuid = ?", uuid).Get(&evi)
	if err != nil {
		return nil, false, base.WrapError("evidence.storage", "get evidence error", err)
	}

	if !exists {
		return nil, false, base.NewError("evidence.storage", "get evidence error")
	}

	return &evi, true, nil
}
