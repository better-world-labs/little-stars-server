package imageprocessing

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg/config"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"os"
	"testing"
)

func init() {
	c, err := config.LoadConfig("../../..")
	if err != nil {
		panic(err)
	}
	interfaces.InitConfig(c)
}

func Test_upload(t *testing.T) {

	//使用方法
	upload("test/upload/gen.png", func(writer *io.PipeWriter) {
		defer writer.Close()
		genImg(writer)
	})
}

func genImg(writer io.Writer) error {
	println("read file")
	fd, err := os.Open("testdata/image.jpeg")
	if err != nil {
		println(err)
		panic(err)
	}

	bgFile, _ := os.Open("testdata/medal_background.png")
	bgImg, _, _ := image.Decode(bgFile)

	lightFile, _ := os.Open("testdata/medal_light.png")
	light, _, _ := image.Decode(lightFile)

	println("创建背景")
	background := image.NewRGBA(bgImg.Bounds())

	println("载入图片")
	img, _, err := image.Decode(fd)

	println("画点东西")
	draw.Draw(background, background.Bounds(), bgImg, image.Point{}, draw.Over)

	println("something")
	drawMedal(background, img, light)
	err = png.Encode(writer, background)
	if err != nil {
		fmt.Printf("err is: %v", err)
	}
	return err
}

func TestGenFile(t *testing.T) {
	file, _ := os.OpenFile("testdata/gen.png", os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()
	genImg(file)
}
