package img

import "image"

type Position interface {
	GetPoint(dst, src image.Image) (x, y int)
}

type RightTop struct {
	Right, Top int
}

func (a RightTop) GetPoint(dst, src image.Image) (x, y int) {
	x = dst.Bounds().Max.X - src.Bounds().Dx()
	y = dst.Bounds().Max.Y + src.Bounds().Dy()
	return x - a.Right, y + a.Top
}

type Absolute struct {
	X, Y int
}

func (a Absolute) GetPoint(dst, src image.Image) (x, y int) {
	return a.X, a.Y
}
