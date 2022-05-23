package service

import (
	"aed-api-server/internal/interfaces/entities"
)

type EvidenceService interface {

	//CreateCertEvidenceAsync 异步创建证书存证
	CreateCertEvidenceAsync(accountID int64, uid string, desc string) chan error

	// CreateEvidence 创建存证
	CreateEvidence(credentialClaim entities.IClaim, name string, accountID int64, category entities.EvidenceCategory, businessKey string) error

	// CreateTextEvidence 创建文本存证
	CreateTextEvidence(name string, accountID int64, category entities.EvidenceCategory, businessKey string, payload map[string]interface{}) error

	// CreateEvidenceAsync 异步创建存证
	CreateEvidenceAsync(credentialClaim entities.IClaim, name string, accountID int64, category entities.EvidenceCategory, businessKey string) chan error

	// GetEvidenceByBusinessKey 读取存证信息
	GetEvidenceByBusinessKey(businessKey string, category entities.EvidenceCategory) (*entities.Evidence, bool, error)

	// GetTransactionViewLink 获取交易验证链接
	GetTransactionViewLink(transactionId string) string

	// GetTransactionViewLinkByBusinessKey 根据 BusinessKey 和 category 获取交易验证链接
	GetTransactionViewLinkByBusinessKey(businessKey string, category entities.EvidenceCategory) (string, error)
}
