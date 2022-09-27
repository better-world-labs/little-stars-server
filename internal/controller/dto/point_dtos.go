package dto

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

func ParseDto(record *entities.UserPointsRecord) *UserPointsRecordDto {
	return &UserPointsRecordDto{
		Id:          record.Id,
		Description: record.Description,
		Points:      record.Points,
		Time:        global.FormattedTime(*record.Time),
	}
}

func ParseDtos(records []*entities.UserPointsRecord) []*UserPointsRecordDto {
	var dtos []*UserPointsRecordDto

	for _, r := range records {
		dtos = append(dtos, ParseDto(r))
	}

	return dtos
}

type AwardFlowDto struct {
	entities.AwardPointFlow

	User *entities.SimpleUser `json:"user"`
}

func ParseAwardFlowDto(flow *entities.AwardPointFlow, user *entities.SimpleUser) *AwardFlowDto {
	return &AwardFlowDto{AwardPointFlow: *flow, User: user}
}

func ParseAwardFlowDtos(flows []*entities.AwardPointFlow, users map[int64]*entities.SimpleUser) []*AwardFlowDto {
	var res []*AwardFlowDto

	for _, f := range flows {
		user := users[f.UserId]
		res = append(res, ParseAwardFlowDto(f, user))
	}

	return res
}
