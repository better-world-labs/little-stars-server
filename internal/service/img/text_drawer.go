package img

import (
	"errors"
	"fmt"
	"github.com/golang/freetype"
	"image"
	"image/color"
	"image/draw"
	"sync"
)

type TextDrawer struct {
	con       *freetype.Context
	mu        sync.Mutex
	fontBytes []byte
}

func NewTextDrawer(fontBytes []byte) (*TextDrawer, error) {
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}
	f := freetype.NewContext()
	f.SetFont(font)

	return &TextDrawer{
		con:       f,
		mu:        sync.Mutex{},
		fontBytes: fontBytes,
	}, nil
}

func (t *TextDrawer) DrawText(dst draw.Image, text string, x int, y int, fontSize float64, color color.Color) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.con.SetDPI(100)
	t.con.SetFontSize(fontSize)
	t.con.SetClip(dst.Bounds())
	t.con.SetDst(dst)
	t.con.SetSrc(image.NewUniform(color))
	pt := freetype.Pt(x, y)
	_, err := t.con.DrawString(text, pt)
	return err
}

func (t *TextDrawer) CreateTextImage(bg draw.Image, text string, fontSize float64, color color.Color) (image.Image, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.con.SetDPI(100)
	t.con.SetFontSize(fontSize)
	t.con.SetClip(bg.Bounds())
	t.con.SetDst(bg)
	t.con.SetSrc(image.NewUniform(color))
	y := int(t.con.PointToFixed(fontSize) >> 6)
	pt := freetype.Pt(0, y)
	p, err := t.con.DrawString(text, pt)
	if err != nil {
		return nil, err
	}

	rect := image.Rect(0, y, int(p.X), int(p.Y))
	fmt.Printf("%d,%d\n", p.X, p.Y)
	rgba := image.NewRGBA(rect)
	inter := bg.Bounds().Intersect(rect)
	if inter.Empty() {
		return nil, errors.New("empty image")
	}
	draw.Draw(rgba, rgba.Bounds(), bg, image.Point{}, draw.Over)
	//return bg,nil
	return rgba, nil

}
