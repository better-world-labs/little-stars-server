package img

import (
	"aed-api-server/internal/pkg/asserts"
	"bytes"
	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
	"image"
	"image/color"
	"image/draw"
	"math"
)

var textDrawer *TextDrawer

func Init() error {
	fontBytes, exists := asserts.GetResource("SourceHanSansSC-Medium.ttf")
	if !exists {
		panic("font file not found")
	}

	drawer, err := NewTextDrawer(fontBytes)
	if err != nil {
		return err
	}

	textDrawer = drawer

	return nil
}

func DrawAvatar(dst draw.Image, avatar image.Image, p image.Point, size uint) error {
	avatar = resize.Resize(size, size, avatar, resize.Lanczos2)
	avatarCircle := image.NewRGBA(avatar.Bounds())
	rx := avatar.Bounds().Max.X
	ry := avatar.Bounds().Max.Y
	r := int(math.Min(float64(rx), float64(ry))) / 2
	draw.DrawMask(avatarCircle, avatar.Bounds(), avatar, image.Point{}, &circle{image.Point{r, r}, r}, image.Point{}, draw.Over)
	draw.Draw(dst, dst.Bounds(), avatarCircle, p, draw.Over)
	return nil
}

func DrawCircleBlur(dst draw.Image, c color.Color, p image.Point, size uint, blur float64) error {
	rgba := image.NewRGBA(image.Rect(0, 0, int(size), int(size)))
	for i := 0; i < int(blur); i++ {
		for j := 0; j < int(blur); j++ {
			rgba.Set(i, j, c)
		}
	}
	rx := rgba.Bounds().Max.X
	ry := rgba.Bounds().Max.Y
	r := int(math.Min(float64(rx), float64(ry))) / 2
	draw.DrawMask(rgba, rgba.Bounds(), rgba, image.Point{}, &circle{image.Point{r, r}, r}, image.Point{}, draw.Over)
	draw.Draw(dst, dst.Bounds(), imaging.Blur(rgba, blur), p, draw.Over)
	return nil
}

func DrawTextBlur(dst draw.Image, text string, x int, y int, fontSize float64, color color.Color, blur float64) error {
	back := image.NewRGBA(dst.Bounds())
	err := DrawText(back, text, x, y, fontSize, color)
	if err != nil {
		return err
	}

	draw.Draw(dst, dst.Bounds(), imaging.Blur(back, blur), image.Point{}, draw.Over)
	return nil
}

//TODO 无法确定Text的占用Size,不能写通用的定位逻辑
func DrawText(dst draw.Image, text string, x int, y int, fontSize float64, color color.Color) error {
	return textDrawer.DrawText(dst, text, x, y, fontSize, color)
}

func DrawTextAutoBreakRune(dst draw.Image, text []rune, x, y, wordsCount, lineHeight int, fontSize float64, color color.Color) error {
	for len(text) > 0 {
		if len(text) > wordsCount {
			err := DrawText(dst, string(text[:wordsCount]), x, y, fontSize, color)
			if err != nil {
				return err
			}

			text = text[wordsCount:]
			y += lineHeight + int(fontSize)
			continue
		}

		err := DrawText(dst, string(text), x, y, fontSize, color)
		if err != nil {
			return err
		}

		text = nil
	}

	return nil
}

func DrawTextAutoBreakASCII(dst draw.Image, text string, x, y, wordsCount, lineHeight int, fontSize float64, color color.Color) error {
	for len(text) > 0 {
		if len(text) > wordsCount {
			err := DrawText(dst, text[:wordsCount], x, y, fontSize, color)
			if err != nil {
				return err
			}

			text = text[wordsCount:]
			y += lineHeight + int(fontSize)
			continue
		}

		err := DrawText(dst, text, x, y, fontSize, color)
		if err != nil {
			return err
		}

		text = ""
	}

	return nil
}

func DrawQrCode(dst draw.Image, qrSize int, content string, x int, y int) error {
	b, err := qrcode.Encode(content, qrcode.Medium, qrSize)
	if err != nil {
		return err
	}

	qrCodeImg, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return err
	}

	draw.Draw(dst, dst.Bounds(), qrCodeImg, image.Point{x, y}, draw.Over)
	return err
}

func DrawQrCodeRightBottom(dst draw.Image, qrSize int, content string, marginRight, marginBottom int) error {
	x := dst.Bounds().Max.X - qrSize - marginRight
	y := dst.Bounds().Max.Y - qrSize - marginBottom
	err := DrawQrCode(dst, qrSize, content, -x, -y)
	return err
}
