package claim

import (
	"aed-api-server/internal/interfaces"
)

type Medal struct {
	Mobile string `json:"mobile"`
	Medal  string `json:"medal"`
}

func (Medal) CptID() int {
	config := interfaces.GetConfig()
	return config.CptMedal
}
