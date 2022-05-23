package point

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/global"
)

type UserPointsRecordDto struct {
	Id          string               `json:"id"`
	Description string               `json:"description"`
	Points      int                  `json:"points"`
	Time        global.FormattedTime `json:"time"` // *global.FormattedTime 会在查询引发 panic
}

func parseDto(record *entities.UserPointsRecord) *UserPointsRecordDto {
	return &UserPointsRecordDto{
		Id:          record.Id,
		Description: record.Description,
		Points:      record.Points,
		Time:        global.FormattedTime(*record.Time),
	}
}

func parseDtos(records []*entities.UserPointsRecord) []*UserPointsRecordDto {
	var dtos []*UserPointsRecordDto

	for _, r := range records {
		dtos = append(dtos, parseDto(r))
	}

	return dtos
}
