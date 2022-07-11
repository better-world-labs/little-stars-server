package credential

import (
	"aed-api-server/internal/interfaces/entities"
)

type IService interface {
	CreateCredential(claim entities.IClaim) (*Info, error)
}
