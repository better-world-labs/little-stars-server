package credential

type Info struct {
	Claim           map[string]interface{}
	VcId            string `json:"vc_id"`
	IssuanceDate    int64  `json:"issuance_date"`
	SignatureValue  string `json:"signature_value"`
	EvidenceFileURL string `json:"evidence_file_url"`
}
