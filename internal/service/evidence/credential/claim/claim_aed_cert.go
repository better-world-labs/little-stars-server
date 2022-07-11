package claim

import "aed-api-server/internal/interfaces"

type AedCert struct {
	User   string `json:"user"`
	Detail string `json:"detail"`
}

func (AedCert) CptID() int {
	config := interfaces.GetConfig()
	return config.CptAedCert
}
