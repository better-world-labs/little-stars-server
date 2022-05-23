package entities

import "time"

type (
	EvidenceCreateResponse struct {
		ID              string `json:"id"`
		TransactionHash string `json:"transaction_hash"`
		Type            int    `json:"type"`
	}

	EvidenceDetailInfo struct {
		EvidenceCreateResponse

		Content string `json:"content"`
		Created int64  `json:"timestamp"`
	}

	Evidence struct {
		ID              int64            `xorm:"id pk autoincr"`
		AccountID       int64            `xorm:"account_id"`
		BusinessKey     string           `xorm:"business_key"`
		Category        EvidenceCategory `xorm:"category"`
		EvidenceID      string           `xorm:"evidence_id"`
		ContentHash     string           `xorm:"content_hash"`
		TransactionHash string           `xorm:"transaction_hash"`
		CredentialID    string           `xorm:"credential_id"`
		FileBytes       int64            `xorm:"file_bytes"`
		Time            time.Time        `xorm:"time"`
	}
)

const (
	EvidenceCategoryCert     EvidenceCategory = 1
	EvidenceCategoryMedal                     = 2
	EvidenceCategoryDonation                  = 3
)

type EvidenceCategory int
