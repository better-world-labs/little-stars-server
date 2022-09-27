package dto

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/global"
	"time"
)

type (
	Dto struct {
		entities.Donation                                 //TODO 废弃
		ProcessPercent             float32                `json:"processPercent"`             //TODO 废弃
		CrowdFundingProcessPercent float32                `json:"crowdfundingProcessPercent"` //TODO 废弃
		StartAt                    global.FormattedTime   `json:"startAt"`                    //TODO 废弃
		CompleteAt                 *global.FormattedTime  `json:"completeAt"`                 //TODO 废弃
		ExpiredAt                  global.FormattedTime   `json:"expiredAt"`                  //TODO 废弃
		CreatedAt                  global.FormattedTime   `json:"createdAt"`                  //TODO 废弃
		Record                     map[string]interface{} `json:"record"`
	}

	WithDonatedDto struct {
		Dto

		DonatedPoints int `json:"donatedPoints"`
	}

	RecordDto struct {
		DonationId int64               `json:"donationId"`
		Points     int                 `json:"points"`
		Donator    entities.SimpleUser `json:"donator"`
	}

	EvidenceDto struct {
		ViewLink         string `json:"viewLink"`
		EvidenceImageUrl string `json:"evidenceImageUrl"`
	}
)

func DtoFromEntity(donation *entities.Donation) *Dto {
	if donation == nil {
		return nil
	}

	dto := Dto{
		Donation:                   *donation,
		StartAt:                    global.FormattedTime(donation.StartAt),
		ExpiredAt:                  global.FormattedTime(donation.ExpiredAt),
		CreatedAt:                  global.FormattedTime(donation.CreatedAt),
		ProcessPercent:             donation.GetProcessPercents(),
		CrowdFundingProcessPercent: donation.GetCrowdfundingProcessPercents(),
	}
	if donation.CompleteAt != nil {
		completeAt := global.FormattedTime(*donation.CompleteAt)
		dto.CompleteAt = &completeAt
	}

	return &dto
}

func DtosFromEntities(donations []*entities.Donation) []*Dto {
	dtos := make([]*Dto, len(donations))

	for i := range donations {
		dtos[i] = DtoFromEntity(donations[i])
	}

	return dtos
}

func EntityFromDto(dto *Dto) *entities.Donation {
	return entities.NewDonation(
		dto.Id,
		dto.Title,
		dto.Images,
		dto.Description,
		dto.TargetPoints,
		0,
		time.Time(dto.StartAt),
		nil,
		time.Time(dto.ExpiredAt),
		dto.ArticleId,
		dto.Executor,
		dto.ExecutorNumber,
		dto.Feedback,
		0,
		dto.Plan,
		dto.PlanImage,
		dto.Budget)
}

func WithDonatedDtoFromEntity(donated *entities.DonationWithUserDonated) *WithDonatedDto {
	return &WithDonatedDto{
		Dto:           *DtoFromEntity(&donated.Donation),
		DonatedPoints: donated.DonatedPoints,
	}
}

func WithDonatedDtosFromEntities(donations []*entities.DonationWithUserDonated) []*WithDonatedDto {
	dtos := make([]*WithDonatedDto, len(donations))

	for i := range donations {
		dtos[i] = WithDonatedDtoFromEntity(donations[i])
	}

	return dtos
}

func RecordDtoFromEntity(record *entities.DonationRecord, user *entities.SimpleUser) *RecordDto {
	return &RecordDto{
		DonationId: record.DonationId,
		Donator:    *user,
		Points:     record.Points,
	}
}

func RecordDtosFromEntities(records []*entities.DonationRecord, users []*entities.SimpleUser) []*RecordDto {
	dtos := make([]*RecordDto, len(records))
	usersMap := make(map[int64]*entities.SimpleUser, 0)
	for _, u := range users {
		usersMap[u.ID] = u
	}

	for i, r := range records {
		user := usersMap[r.UserId]
		dtos[i] = RecordDtoFromEntity(r, user)
	}

	return dtos
}
