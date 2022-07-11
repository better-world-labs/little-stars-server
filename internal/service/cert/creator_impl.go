package cert

import (
	"aed-api-server/internal/pkg/asserts"
	"aed-api-server/internal/service/img"
	"bytes"
	"errors"
	"fmt"
	"go.uber.org/zap/buffer"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"time"
)

const (
	FileCertBackground = "cert_origin.png"
	FileCertLogo       = "cert_logo.png"
)

type imageCreator struct {
	backgroundBytes []byte
	certLogoBytes   []byte
}

func NewImageCreatorDefaultAssert() (ImageCreator, error) {
	return NewImageCreator()
}

func NewImageCreator() (ImageCreator, error) {
	certBackgroundBytes, exists := asserts.GetResource(FileCertBackground)
	if !exists {
		return nil, errors.New("file not found")
	}

	certLogoBytes, exists := asserts.GetResource(FileCertLogo)
	if !exists {
		return nil, errors.New("file not found")
	}

	return &imageCreator{
		backgroundBytes: certBackgroundBytes,
		certLogoBytes:   certLogoBytes,
	}, nil
}

func (i *imageCreator) Create(avatarUrl string, nickname string, description string, t time.Time, writer io.Writer) error {
	background, _, err := image.Decode(bytes.NewReader(i.backgroundBytes))
	if err != nil {
		return err
	}

	dst := image.NewRGBA(background.Bounds())
	draw.Draw(dst, background.Bounds(), background, image.Point{}, draw.Over)

	err = i.drawLogo(dst)
	if err != nil {
		return err
	}

	err = i.drawAvatar(dst, avatarUrl)
	if err != nil {
		return err
	}

	err = i.drawText(dst, nickname, description, t)
	if err != nil {
		return err
	}

	err = jpeg.Encode(writer, dst, &jpeg.Options{Quality: 90})
	if err != nil {
		return err
	}

	return nil
}

func (i *imageCreator) drawLogo(dst draw.Image) error {
	certLogo, _, err := image.Decode(bytes.NewReader(i.certLogoBytes))

	if err != nil {
		return err
	}

	draw.Draw(dst, dst.Bounds(), certLogo, image.Point{-255, -30}, draw.Over)
	return err
}

func (imageCreator) drawAvatar(dst draw.Image, avatarUrl string) error {
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

	err = img.DrawAvatar(dst, avatar, image.Point{-37, -1186}, 60)

	return err
}

func (i *imageCreator) drawText(dst draw.Image, nickname string, description string, t time.Time) error {
	err := img.DrawText(dst, nickname, 107, 1210, 18, color.White)
	if err != nil {
		return nil
	}

	err = img.DrawText(dst, description, 107, 1243, 14, color.White)
	if err != nil {
		return nil
	}

	date := t.Format("2006/01/02")
	b := buffer.Buffer{}
	for i := 0; i < len(date); i++ {
		b.AppendByte(date[i])
		b.AppendString("  ")
	}

	err = img.DrawText(dst, fmt.Sprintf("小星星项目组    %s", b.String()), 36, 1288, 12, color.White)

	return nil
}
