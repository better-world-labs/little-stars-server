package service

type MedalService interface {
	AwardMedalSaveLife(userId, helpInfoId int64) error
	AwardMedalFirstDonation(userId, donationRecordId int64) error
	AwardMedalInspector(userId int64) error
}

type UserMedalService interface {
	GetUserMedalUrl(userId int64) ([]string, error)
}
