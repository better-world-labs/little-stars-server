package img

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"golang.org/x/image/colornames"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"testing"
)

func TestCreateTextImage(t *testing.T) {
	fontBytes, err := os.ReadFile("/usr/share/fonts/jet-brans-mono/JetBrains Mono Bold Italic Nerd Font Complete Mono.ttf")
	require.Nil(t, err)

	drawer, err := NewTextDrawer(fontBytes)
	require.Nil(t, err)

	certBackgroundBytes, err := ioutil.ReadFile("../../../assert/cert_origin.png")
	require.Nil(t, err)
	background, _, err := image.Decode(bytes.NewReader(certBackgroundBytes))
	dst := image.NewRGBA(background.Bounds())

	img, err := drawer.CreateTextImage(dst, "Souththth_Tae🍑", 26, colornames.Yellow)
	require.Nil(t, err)
	f, err := os.Create("img.png")
	defer f.Close()
	require.Nil(t, err)
	err = jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
	require.Nil(t, err)
}
