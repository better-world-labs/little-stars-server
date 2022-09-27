package service

import (
	"aed-api-server/internal/interfaces/entities"
	page "aed-api-server/internal/pkg/query"
)

type DonationService interface {

	// Donate 积分捐献
	Donate(record *entities.DonationRecord) (*entities.DonationRecord, error)

	// CreateDonation 创建一个捐献项目
	CreateDonation(donation *entities.Donation) error

	// UpdateDonation 根据 ID 更新捐献项目
	UpdateDonation(donation *entities.Donation) error

	// DeleteDonation 删除一个捐献项目
	DeleteDonation(id int64) error

	// GetDonationById 根据 ID 读取一个捐献项目
	GetDonationById(id int64) (*entities.Donation, bool, error)

	// GetDonationDetail 根据 ID 读取一个捐献项目
	GetDonationDetail(id int64) (*entities.Donation, bool, error)

	// GetRecordById 根据记录ID查询记录
	GetRecordById(recordId int64) (*entities.DonationRecord, bool, error)

	// ListDonation 捐献项目列表
	ListDonation() ([]*entities.Donation, error)

	// ListDonationSorted 捐献项目列表
	ListDonationSorted(query page.Query, userId int64) ([]*entities.Donation, error)

	// ListDonorsDonation 根据捐献者查询捐献项目列表
	ListDonorsDonation(userId int64) ([]*entities.DonationWithUserDonated, error)

	// ListUserPointsForDonations 列出用户对项目的捐献情况
	ListUserPointsForDonations(userId int64) (map[int64]int, error)

	// ListRecords 读取项目的最近 n 条捐献记录， latest <= 0 则不限制
	ListRecords(donationId int64, latest int) ([]*entities.DonationRecord, error)

	// ListUsersRecordsTop 读取用户捐献前 n 名
	ListUsersRecordsTop(donationId int64, top int) ([]*entities.DonationRecord, error)

	Apply(apply entities.DonationApply, userId int64) error

	UpdateCrowdfunding(id int64, actualCrowdfunding float32) error

	GetDonationHonor(user *entities.User) (*entities.DonationHonor, error)

	//CountUserRecord 获取用户捐献次数
	CountUserRecord(userId int64) (int, error)

	StatDonationByUserId(userId int64) (stat entities.DonationStat, err error)

	StatUsersDonations(userIds []int64) (list []*entities.UserDonationPoints, err error)
}
