package imageprocessing

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/asserts"
	"aed-api-server/internal/pkg/base"
	"aed-api-server/internal/service/img"
	"aed-api-server/internal/service/medal"
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"net/http"
	"strconv"
	"time"
)

var shareBg image.Image

func aidService() service.AidService {
	return interfaces.S.Aid
}
func DrawMedalShare(medalId int64, account *entities.User, writer io.Writer, serverHost string) error {
	userMedal, exists, err := medal.GetUserMedal(medalId, account.ID)
	if err != nil {
		return err
	}

	if !exists {
		return base.NewError("imageprocessing", "not found")
	}
	medal, exists, err := medal.GetById(userMedal.MedalID)
	if err != nil {
		return err
	}
	if !exists {
		return base.NewError("imageprocessing", "not found")
	}

	return doDrawMedalShare(medal, userMedal, account, writer, serverHost)
}

func doDrawMedalGeneric(m *entities.Medal, um *entities.UserMedal, u *entities.User, writer io.Writer, serverHost string) error {
	bgBytes, exists := asserts.GetResource(m.ShareBackground)
	if !exists {
		return errors.New("background assert not found")
	}

	bgImage, _, err := image.Decode(bytes.NewReader(bgBytes))
	background := image.NewRGBA(bgImage.Bounds())
	draw.Draw(background, background.Bounds(), bgImage, image.Point{}, draw.Over)

	//drawMedal(background, medal, light)

	blurColor := color.RGBA{235, 90, 50, 255}
	avatarPosition := image.Point{-45, -30}
	err = drawAvatar(background, u.Avatar, avatarPosition)
	if err != nil {
		return err
	}

	err = img.DrawTextBlur(background, u.Nickname, 178, 73, 23, blurColor, 3)
	err = img.DrawText(background, u.Nickname, 175, 75, 23, color.White)

	err = img.DrawTextBlur(background, "人生的路很长很长，让我陪你一起走吧", 178, 113, 18, blurColor, 3)
	err = img.DrawText(background, "人生的路很长很长，让我陪你一起走吧", 175, 115, 18, color.White)
	if err != nil {
		return err
	}

	getMedal := fmt.Sprintf("%s", time.Time(um.Created).Format("2006-01-02 15:04:05"))
	err = img.DrawText(background, getMedal, 167, 1573, 22, color.Black)
	if err != nil {
		return err
	}

	x := background.Bounds().Max.X - 128 - 160
	y := background.Bounds().Max.Y - 128 - 172
	err = img.DrawQrCode(background, 130, fmt.Sprintf("https://%s/share/cert?source=placard&sharer=%d", serverHost, u.ID), -x, -y)
	if err != nil {
		return err
	}

	//err = img.DrawQrCodeRightBottom(background, fmt.Sprintf("%s/share/cert?id=%d", prefix, u.ID))

	return jpeg.Encode(writer, background, &jpeg.Options{Quality: 90})
}
func doDrawMedalSaveLife(m *entities.Medal, um *entities.UserMedal, u *entities.User, writer io.Writer, serverHost string) error {
	bgBytes, exists := asserts.GetResource(m.ShareBackground)
	if !exists {
		return errors.New("background assert not found")
	}

	bgImage, _, err := image.Decode(bytes.NewReader(bgBytes))
	background := image.NewRGBA(bgImage.Bounds())
	draw.Draw(background, background.Bounds(), bgImage, image.Point{}, draw.Over)

	//drawMedal(background, medal, light)

	blurColor := color.RGBA{235, 90, 50, 255}
	avatarPosition := image.Point{-45, -30}
	err = drawAvatar(background, u.Avatar, avatarPosition)
	if err != nil {
		return err
	}

	err = img.DrawTextBlur(background, u.Nickname, 178, 73, 23, blurColor, 3)
	err = img.DrawText(background, u.Nickname, 175, 75, 23, color.White)

	err = img.DrawTextBlur(background, "人生的路很长很长，让我陪你一起走吧", 178, 113, 18, blurColor, 3)
	err = img.DrawText(background, "人生的路很长很长，让我陪你一起走吧", 175, 115, 18, color.White)
	if err != nil {
		return err
	}

	getMedal := fmt.Sprintf("%s", time.Time(um.Created).Format("2006-01-02 15:04:05"))
	err = img.DrawText(background, getMedal, 167, 1573, 22, color.Black)
	if err != nil {
		return err
	}

	if um.BusinessId != "" { // 兼容老数据
		aidId, err := strconv.ParseInt(um.BusinessId, 10, 64)
		if err != nil {
			return err
		}

		helpInfo, exists, err := aidService().GetHelpInfoComposedByID(aidId, nil)
		if exists {
			runes := []rune(helpInfo.Address)
			err = img.DrawText(background, string(runes[:6]), 650, 1573, 22, color.Black)
			if err != nil {
				return err
			}
		}
	}

	x := background.Bounds().Max.X - 128 - 160
	y := background.Bounds().Max.Y - 128 - 172
	err = img.DrawQrCode(background, 130, fmt.Sprintf("https://%s/share/cert?source=placard&sharer=%d", serverHost, u.ID), -x, -y)
	if err != nil {
		return err
	}

	//err = img.DrawQrCodeRightBottom(background, fmt.Sprintf("%s/share/cert?id=%d", prefix, u.ID))

	return jpeg.Encode(writer, background, &jpeg.Options{Quality: 90})
}

func doDrawMedalShare(m *entities.Medal, um *entities.UserMedal, u *entities.User, writer io.Writer, serverHost string) error {
	switch m.ID {
	case entities.MedalIdSaveLife:
		return doDrawMedalSaveLife(m, um, u, writer, serverHost)

	default:
		return doDrawMedalGeneric(m, um, u, writer, serverHost)
	}
}

func drawAvatar(dst draw.Image, avatarUrl string, point image.Point) error {
	res, err := http.Get(avatarUrl)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("get avatar error for cert")
	}

	avatar, _, err := image.Decode(res.Body)
	if err != nil {
		return err
	}

	err = img.DrawAvatar(dst, avatar, point, 80)

	return err
}
