package img

import (
	"errors"
	"github.com/disintegration/imaging"
	"github.com/stretchr/testify/require"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	_ "image/png"
	"net/http"
	"os"
	"testing"
)

func TestDrawText(t *testing.T) {

	bg := image.NewRGBA(image.Rect(0, 0, 128, 128))
	err := DrawQrCode(bg, 128, "https://dev.star.openviewtech.com/subcontract/article/detail?url=https%3A%2F%2Fdev.cms.openviewtech.com%2Farticle%2F3%3FuserId%3D49%26id%3D8", 0, 0)
	//err := DrawQrCode(bg, 128, "https://dev.star.openviewtech.com/share/cert?source=placard&sharer=50", 0, 0)

	require.Nil(t, err)

	create, err := os.Create("qrcode_test.png")
	defer create.Close()
	require.Nil(t, err)
	err = jpeg.Encode(create, bg, &jpeg.Options{Quality: 100})
	require.Nil(t, err)
}

func TestDrawMedalDetail(t *testing.T) {
	err := Init()
	require.Nil(t, err)

	bgImg, err := loadBackground()
	require.Nil(t, err)

	medal, err := loadMedal()
	require.Nil(t, err)

	light, err := loadMedalLight()
	require.Nil(t, err)

	background := image.NewRGBA(bgImg.Bounds())
	draw.Draw(background, background.Bounds(), bgImg, image.Point{}, draw.Over)

	drawMedal(background, medal, light)

	// 右上角
	err = DrawText(background, "2022-01-05 获得勋章", background.Bounds().Max.X-270, background.Bounds().Min.Y+45, 18, color.RGBA{215, 212, 215, 255})
	require.Nil(t, err)

	//err = DrawTextBlur(background, "及时发布救援", 222, 515, 18, color.RGBA{255, 118, 99,255},3)
	//err = DrawText(background, "及时发布救援", 222, 510, 18, color.White)

	err = DrawTextBlur(background, "从死神手里抢回一条生命，", 138, 545, 20, color.RGBA{255, 118, 99, 255}, 4)
	err = DrawText(background, "从死神手里抢回一条生命", 138, 540, 20, color.White)

	err = DrawText(background, "感谢坚守在小星星 100 天的你", 100, 700, 20, color.RGBA{215, 212, 215, 255})
	err = DrawText(background, "因为这份坚守，才让Ta 看见了光", 100, 750, 20, color.RGBA{215, 212, 215, 255})
	create, err := os.Create("medal_out.png")
	defer create.Close()
	require.Nil(t, err)
	err = jpeg.Encode(create, background, &jpeg.Options{Quality: 100})
}

func drawMedal(dst draw.Image, medal image.Image, light image.Image) {
	medal = imaging.Resize(medal, 250, 250, imaging.Lanczos)
	light = imaging.Resize(light, 400, 400, imaging.Lanczos)

	centerX := (dst.Bounds().Max.X - dst.Bounds().Min.X) / 2
	centerY := 225

	draw.Draw(dst, dst.Bounds(), light, image.Point{-(centerX - light.Bounds().Dx()/2), -(centerY - light.Bounds().Dx()/2)}, draw.Over)
	draw.Draw(dst, dst.Bounds(), medal, image.Point{-(centerX - medal.Bounds().Dx()/2), -(centerY - medal.Bounds().Dx()/2)}, draw.Over)
}

func loadMedal() (image.Image, error) {
	res, err := http.Get("https://openview-oss.oss-cn-chengdu.aliyuncs.com/star-static/medal_active.png")
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("get resource with error code")
	}

	i, _, err := image.Decode(res.Body)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func loadMedalLight() (image.Image, error) {
	f, err := os.Open("/home/shenweijie/medal_light.png")
	if err != nil {
		return nil, err
	}

	i, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func loadBackground() (image.Image, error) {
	bgFile, err := os.Open("/home/shenweijie/medal_background.png")
	if err != nil {
		return nil, err
	}

	background, _, err := image.Decode(bgFile)

	if err != nil {
		return nil, err
	}
	return background, nil
}

func TestDrawQrCode(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 1024, 768))
	err := DrawQrCode(img, 128, "https://openscan.openviewtech.com/#/transaction/transactionDetail?pageSize=10&pageNumber=1&v_page=transaction&pkHash=0x5650d4124a73408b2d40c31115147567340b3a89aef173a3cb9d2161d95b6097", 0, 0)
	require.Nil(t, err)
	f, err := os.Create("qrcode.png")
	require.Nil(t, err)
	defer f.Close()
	err = png.Encode(f, img)
	require.Nil(t, err)
}
