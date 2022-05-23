package credential

import (
	"aed-api-server/internal/interfaces/entities"
)

type Service interface {
	CreateCredential(claim entities.IClaim) (*Info, error)
}
