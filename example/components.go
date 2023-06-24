package main

import (
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var tt, _ = opentype.Parse(fonts.MPlus1pRegular_ttf)
var face, _ = opentype.NewFace(tt, &opentype.FaceOptions{
	Size:    32,
	DPI:     72,
	Hinting: font.HintingVertical,
})
