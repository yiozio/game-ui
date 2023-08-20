package common

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
)

var (
	EmptyImage = ebiten.NewImage(1, 1)

	_EmptyImage   = ebiten.NewImage(3, 3)
	EmptySubImage = _EmptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	EmptyImage.Fill(color.White)
	_EmptyImage.Fill(color.White)
}

func ColorCodeToRGBA(color uint32) (float32, float32, float32, float32) {
	var r = float32((0xff000000&color)>>0o30) / 0xff
	var g = float32((0x00ff0000&color)>>0o20) / 0xff
	var b = float32((0x0000ff00&color)>>0o10) / 0xff
	var a = float32(0x000000ff&color) / 0xff

	return r, g, b, a
}

func ColorCodeToVertex(vertex *ebiten.Vertex, color uint32) {
	vertex.ColorR, vertex.ColorG, vertex.ColorB, vertex.ColorA = ColorCodeToRGBA(color)
}
