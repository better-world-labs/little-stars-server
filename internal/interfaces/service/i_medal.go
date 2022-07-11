package service

import "aed-api-server/internal/interfaces/entities"

type IMedal interface {
	AwardMedalSaveLife(userId, helpInfoId int64) error
	AwardMedalFirstDonation(userId, donationRecordId int64) error
	AwardMedalInspector(userId int64) error
	GetById(id int64) (*entities.Medal, bool, error)
}
type IUserMedal interface {
	GetUserMedalUrl(userId int64) ([]string, error)
}
