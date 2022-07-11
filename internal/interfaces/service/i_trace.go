package service

import "aed-api-server/internal/interfaces/entities"

type TraceService interface {
	Create(code string, trace entities.Trace) (*entities.Trace, error)
	GetEarliestSharerTrace(openid string) (*entities.Trace, bool, error)
}
