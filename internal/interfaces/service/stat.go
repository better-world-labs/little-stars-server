package service

import "aed-api-server/internal/interfaces/entities"

type StatService interface {
	DoKipStat() (*entities.KpiStatItem, error)
}
