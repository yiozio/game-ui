package game_ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
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

type sizeSeg struct {
	t    sizeType
	v    float32
	calc []sizeSeg
}

type sizeType int

const (
	px sizeType = 0
	vw sizeType = 1
	vh sizeType = 2
)

func Px(value int, calc ...sizeSeg) *sizeSeg {
	return &sizeSeg{px, float32(value), calc}
}
func Vw(value float32, calc ...sizeSeg) *sizeSeg {
	return &sizeSeg{vw, value, calc}
}
func Vh(value float32, calc ...sizeSeg) *sizeSeg {
	return &sizeSeg{vh, value, calc}
}

func calcSize(screenSize image.Point, size sizeSeg) int {
	var calc = 0
	for _, size := range size.calc {
		calc += calcSize(screenSize, size)
	}
	switch size.t {
	case px:
		return int(size.v) + calc
	case vw:
		return int(float32(screenSize.X)*size.v) + calc
	case vh:
		return int(float32(screenSize.Y)*size.v) + calc
	}
	return 0
}

func Size1(size *sizeSeg) *[4]sizeSeg {
	return &[4]sizeSeg{*size, *size, *size, *size}
}
func Size2(vertical, horizontal *sizeSeg) *[4]sizeSeg {
	return &[4]sizeSeg{*vertical, *horizontal, *vertical, *horizontal}
}
func Size3(top, horizontal, bottom *sizeSeg) *[4]sizeSeg {
	return &[4]sizeSeg{*top, *horizontal, *bottom, *horizontal}
}
func Size4(top, right, bottom, left *sizeSeg) *[4]sizeSeg {
	return &[4]sizeSeg{*top, *right, *bottom, *left}
}

func Radius1(px int) *[4]int {
	return &[4]int{px, px, px, px}
}
func Radius4(topLeft, topRight, bottomRight, bottomLeft int) *[4]int {
	return &[4]int{topLeft, topRight, bottomRight, bottomLeft}
}

func ColorCode1(code uint32) *[4]color.Color {
	var c = Color(code)
	return &[4]color.Color{*c, *c, *c, *c}
}
func ColorCodeHorizontal(left, right uint32) *[4]color.Color {
	var c1 = Color(left)
	var c2 = Color(right)
	return &[4]color.Color{*c1, *c2, *c2, *c1}
}
func ColorCodeVertical(top, bottom uint32) *[4]color.Color {
	var c1 = Color(top)
	var c2 = Color(bottom)
	return &[4]color.Color{*c1, *c1, *c2, *c2}
}
func ColorCodeGradation(topLeft, topRight, bottomRight, bottomLeft uint32) *[4]color.Color {
	var c1 = Color(topLeft)
	var c2 = Color(topRight)
	var c3 = Color(bottomRight)
	var c4 = Color(bottomLeft)
	return &[4]color.Color{*c1, *c2, *c3, *c4}
}

func Color(code uint32) *color.Color {
	// #RRGGBBAA
	var r = (code & 0xff000000) >> (6 * 4)
	var g = (code & 0x00ff0000) >> (4 * 4)
	var b = (code & 0x0000ff00) >> (2 * 4)
	var a = (code & 0x000000ff) >> (0 * 4)

	var c color.Color = color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}

	return &c
}
