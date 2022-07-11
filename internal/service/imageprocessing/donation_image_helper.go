package imageprocessing

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/asserts"
	"aed-api-server/internal/pkg/star"
	"aed-api-server/internal/pkg/utils"
	"aed-api-server/internal/service/img"
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"strconv"
)

const (
	MaxLineRunes = 15
)

func DrawDonationShare(donationId int64, u *entities.User, w *io.PipeWriter) error {
	donation, exists, err := interfaces.S.Donation.GetDonationDetail(donationId)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("donation not found")
	}

	donationsPoinsts, err := interfaces.S.Donation.ListUserPointsForDonations(u.ID)
	if err != nil {
		return err
	}

	bgBytes, exists := asserts.GetResource("donation_share_bg.png")
	if !exists {
		return errors.New("background assert not found")
	}

	bgImage, _, err := image.Decode(bytes.NewReader(bgBytes))
	background := image.NewRGBA(bgImage.Bounds())
	draw.Draw(background, background.Bounds(), bgImage, image.Point{}, draw.Over)

	// 头像昵称
	avatarPosition := image.Point{-78, -596}
	err = drawAvatar(background, u.Avatar, avatarPosition)
	if err != nil {
		return err
	}

	err = img.DrawText(background, u.Nickname, 166, 649, 22, color.Black)
	if err != nil {
		return err
	}

	err = img.DrawText(background, "捐赠积分", 450, 649, 22, color.RGBA{222, 222, 222, 255})
	if err != nil {
		return err
	}

	err = img.DrawText(background, strconv.Itoa(donationsPoinsts[donation.Id]), 590, 649, 22, color.RGBA{255, 116, 97, 255})
	if err != nil {
		return err
	}

	err = img.DrawText(background, "已收到积分:", 100, 1061, 18, color.Black)
	if err != nil {
		return err
	}

	err = img.DrawText(background, utils.PointsString(donation.ActualPoints), 255, 1061, 18, color.RGBA{255, 116, 97, 255})
	if err != nil {
		return err
	}

	err = img.DrawText(background, "尚缺:", 480, 1061, 18, color.Black)
	if err != nil {
		return err
	}

	err = img.DrawText(background, utils.PointsString(donation.TargetPoints-donation.ActualPoints), 557, 1061, 18, color.RGBA{255, 116, 97, 255})
	if err != nil {
		return err
	}

	err = img.DrawQrCodeRightBottom(background, 128, star.GetPlaceCardSharedQrCodeContent(u.ID), 50, 50)
	if err != nil {
		return err
	}

	title := utils.StringLimitHidden(donation.Title, 35, "...")
	titleRunes := []rune(fmt.Sprintf("\"%s\"", title))
	y := getTitlePositionY(titleRunes)

	err = img.DrawText(background, "支持项目", 120, y, 22, color.Black)
	if err != nil {
		return err
	}

	err = img.DrawTextAutoBreakRune(background, titleRunes, 250, y, MaxLineRunes, 30, 22, color.RGBA{255, 116, 97, 255})
	if err != nil {
		return err
	}

	return jpeg.Encode(w, background, &jpeg.Options{Quality: 90})
}

// getTitlePositionY 过去 title 绘制的 Y 坐标,因为如果是两行，需要整体上移居中
func getTitlePositionY(titleRune []rune) int {
	if len(titleRune) > MaxLineRunes {
		return 825
	}

	return 851
}
