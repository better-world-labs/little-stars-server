package evidence

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/module/evidence/credential"
	"aed-api-server/internal/module/evidence/credential/claim"
	"aed-api-server/internal/module/user"
	"aed-api-server/internal/pkg/base"
	"aed-api-server/internal/pkg/config"
	"aed-api-server/internal/pkg/db"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gitlab.openviewtech.com/openview-pub/gopkg/uuid2"
	"net/http"
)

type service struct {
	credentialService credential.Service
	storage           Storage
	api               Api
	accountService    user.Service
}

func Init(conf *config.AppConfig) {
	interfaces.S.Evidence = NewService(*conf)
}
func NewService(config config.AppConfig) *service {
	uuid2.InitSnowFlake(int64(uuid.New().ID()) % 1023)
	return &service{
		credentialService: credential.NewService(config.CredentialConfig),
		api:               NewApi(config.EvidenceConfig),
		accountService:    user.NewService(nil),
		storage:           NewStorage(),
	}
}

func (s service) CreateCertEvidenceAsync(accountID int64, uid string, desc string) chan error {
	errChan := make(chan error, 1)
	go func() {
		defer close(errChan)
		account, err := s.accountService.GetUserByID(accountID)
		if err != nil {
			errChan <- err
			return
		}

		if account == nil {
			errChan <- base.NewError("skill.service", "account not found")
			return
		}

		err = s.CreateEvidence(&claim.AedCert{
			User:   account.Mobile,
			Detail: desc,
		}, "User AED Cert Evidence", accountID, entities.EvidenceCategoryCert, uid)

		errChan <- err
	}()

	return errChan
}

func (s *service) CreateEvidence(credentialClaim entities.IClaim, name string, accountID int64, category entities.EvidenceCategory, businessKey string) error {
	credentialInfo, err := s.credentialService.CreateCredential(credentialClaim)
	if err != nil {
		return base.WrapError("evidence.service", "CreateCredential error", err)
	}

	fileURL := credentialInfo.EvidenceFileURL
	if fileURL == "" {
		return base.NewError("evidence.CreateEvidence", "create credential without file_url")
	}

	resp, err := http.Get(fileURL)
	if err != nil {
		return base.WrapError("evidence.CreateEvidence", "download credential failed", err)
	}

	length := resp.ContentLength
	evi, err := s.api.CreateFileEvidence(name, resp.Body)
	if err != nil {
		return base.WrapError("evidence.onProjectCertificated", "create file evidence error", err)
	}

	info, err := s.api.GetEvidenceInfo(evi.ID)
	if err != nil {
		return base.WrapError("evidence.onProjectCertificated", "get evidence info error", err)
	}

	session := db.GetSession()
	defer session.Close()
	err = s.storage.CreateOrUpdateWithSession(session, entities.Evidence{
		AccountID:       accountID,
		Category:        category,
		BusinessKey:     businessKey,
		EvidenceID:      info.ID,
		ContentHash:     info.Content,
		TransactionHash: info.TransactionHash,
		CredentialID:    credentialInfo.VcId,
		FileBytes:       length,
	})

	if err != nil {
		return base.WrapError("evidence.onProjectCertificated", "store evidence error", err)
	}

	return nil
}

func (s *service) CreateEvidenceAsync(credentialClaim entities.IClaim, name string, accountID int64, category entities.EvidenceCategory, businessKey string) chan error {
	resChan := make(chan error, 1)
	go func() {
		defer close(resChan)
		err := s.CreateEvidence(credentialClaim, name, accountID, category, businessKey)
		resChan <- err
	}()

	return resChan
}

func (s *service) GetEvidenceByBusinessKey(businessKey string, category entities.EvidenceCategory) (*entities.Evidence, bool, error) {
	session := db.GetSession()
	defer session.Close()
	return s.storage.GetOneByBusinessKeyWithSession(session, category, businessKey)
}

func (s *service) CreateTextEvidence(name string, accountID int64, category entities.EvidenceCategory, businessKey string, payload map[string]interface{}) error {
	evidence, err := s.api.CreateTextEvidence(name, payload)
	if err != nil {
		return err
	}

	info, err := s.api.GetEvidenceInfo(evidence.ID)
	if err != nil {
		return base.WrapError("evidence.onProjectCertificated", "get evidence info error", err)
	}

	session := db.GetSession()
	defer session.Close()

	err = s.storage.CreateOrUpdateWithSession(session, entities.Evidence{
		AccountID:       accountID,
		Category:        category,
		BusinessKey:     businessKey,
		EvidenceID:      info.ID,
		ContentHash:     info.Content,
		TransactionHash: info.TransactionHash,
	})

	return err
}

func (s *service) GetTransactionViewLink(transactionId string) string {
	return fmt.Sprintf("https://openscan.openviewtech.com/#/transaction/transactionDetail?pageSize=10&pageNumber=1&v_page=transaction&pkHash=%s", transactionId)
}

func (s *service) GetTransactionViewLinkByBusinessKey(businessKey string, category entities.EvidenceCategory) (string, error) {
	evidence, exists, err := s.GetEvidenceByBusinessKey(businessKey, category)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", errors.New("evidence not found")
	}

	return fmt.Sprintf("https://openscan.openviewtech.com/#/transaction/transactionDetail?pageSize=10&pageNumber=1&v_page=transaction&pkHash=%s", evidence.TransactionHash), nil
}
