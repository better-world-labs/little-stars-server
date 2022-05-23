package imageprocessing

import (
	"aed-api-server/internal/module/aid"
	"aed-api-server/internal/module/user"
	"github.com/disintegration/imaging"
	"image"
	"image/draw"
)

func Init() {
	aidService = aid.NewService(user.NewService(nil))
}

func drawMedal(dst draw.Image, medal image.Image, light image.Image) {
	medal = imaging.Resize(medal, 250, 250, imaging.Lanczos)
	light = imaging.Resize(light, 400, 400, imaging.Lanczos)

	centerX := (dst.Bounds().Max.X - dst.Bounds().Min.X) / 2
	centerY := 225

	draw.Draw(dst, dst.Bounds(), light, image.Point{-(centerX - light.Bounds().Dx()/2), -(centerY - light.Bounds().Dx()/2)}, draw.Over)
	draw.Draw(dst, dst.Bounds(), medal, image.Point{-(centerX - medal.Bounds().Dx()/2), -(centerY - medal.Bounds().Dx()/2)}, draw.Over)
}
