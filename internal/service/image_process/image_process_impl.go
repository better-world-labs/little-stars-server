package image_process

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/star"
	"aed-api-server/internal/pkg/utils"
	"errors"
	"fmt"
	"time"
)

type ImageProcess struct {
	*ImageBot `inject:"-"`

	Donation service.DonationService `inject:"-"`
	User     service.UserServiceOld  `inject:"-"`
}

const (
	DonationShareTpl = "complete-donation-card"
)

//go:inject-component
func NewImageProcess() service.IImageProcess {
	return &ImageProcess{}
}

func (i *ImageProcess) DrawDonationShareImage(recordId, userId int64) (string, error) {
	basePromise, err := utils.PromiseAll(func() (interface{}, error) {
		record, exists, err := i.Donation.GetRecordById(recordId)
		if err != nil {
			return nil, err
		}

		if !exists {
			return nil, errors.New("record not found")
		}

		return record, nil
	}, func() (interface{}, error) {
		return i.User.GetUserByID(userId)
	})

	if err != nil {
		return "", err
	}

	record := basePromise[0].(*entities.DonationRecord)
	user := basePromise[1].(*entities.User)

	all, err := utils.PromiseAll(func() (interface{}, error) {
		return i.Donation.GetDonationHonor(user)
	}, func() (interface{}, error) {
		donation, exists, err := i.Donation.GetDonationDetail(record.DonationId)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("donation not found")
		}
		return donation, nil
	}, func() (interface{}, error) {
		return i.Donation.StatDonationByUserId(userId)
	})
	if err != nil {
		return "", err
	}

	honor := all[0].(*entities.DonationHonor)
	return i.Call(DonationShareTpl, map[string]interface{}{
		"userAvatar":        user.Avatar,
		"username":          user.Nickname,
		"donationPoints":    record.Points,
		"projectName":       all[1].(*entities.Donation).Title,
		"todayStr":          time.Now().Format("2006年06月02日"),
		"medals":            honor.Medals,
		"cashValue":         honor.EquivalentRMB,
		"addStarDays":       honor.RegisteredDays,
		"moreThanRatio":     honor.ExceedPercents,
		"donationAllPoints": honor.TotalDonatedPoints,
		"qrContent":         star.GetPlaceCardSharedQrCodeContent(user.ID),
	}, fmt.Sprintf("%d-%d.png", userId, record.DonationId))
}
