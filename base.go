package game_ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
	"strconv"
)

var (
	emptyImage = ebiten.NewImage(3, 3)

	// emptySubImage is an internal sub image of emptyImage.
	// Use emptySubImage at DrawTriangles instead of emptyImage in order to avoid bleeding edges.
	emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

type Component interface {
	GetSize() image.Point
	Draw(screen *ebiten.Image, x, y int)
}

func init() {
	emptyImage.Fill(color.White)
}

func colorCodeToColor(colorCode string) color.Color {
	var r int64
	var g int64
	var b int64
	var a int64
	var num, _ = strconv.ParseInt(colorCode[1:], 16, 64)

	if len(colorCode) == 4 {
		// #RGB
		r = ((num & 0xf00) >> (2 * 4)) * 17
		g = ((num & 0x0f0) >> (1 * 4)) * 17
		b = ((num & 0x00f) >> (0 * 4)) * 17
		a = 0xff
	} else if len(colorCode) == 7 {
		// #RRGGBB
		r = (num & 0xff0000) >> (4 * 4)
		g = (num & 0x00ff00) >> (2 * 4)
		b = (num & 0x0000ff) >> (0 * 4)
		a = 0xff
	} else if len(colorCode) == 9 {
		// #RRGGBBAA
		r = (num & 0xff000000) >> (6 * 4)
		g = (num & 0x00ff0000) >> (4 * 4)
		b = (num & 0x0000ff00) >> (2 * 4)
		a = (num & 0x000000ff) >> (0 * 4)
	} else {
		panic("illegal color code: " + colorCode)
	}

	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
}
