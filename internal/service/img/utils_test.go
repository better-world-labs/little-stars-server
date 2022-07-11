package img

import (
	"github.com/stretchr/testify/require"
	"image"
	"image/jpeg"
	"image/png"
	_ "image/png"
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
