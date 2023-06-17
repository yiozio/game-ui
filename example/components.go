package main

import (
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"yioz.io/game-ui"
)

var tt, _ = opentype.Parse(fonts.MPlus1pRegular_ttf)
var face, _ = opentype.NewFace(tt, &opentype.FaceOptions{
	Size:    32,
	DPI:     72,
	Hinting: font.HintingVertical,
})
var titleFont = game_ui.NewTextFont(
	face,
	0,
	0,
)

var titleText = game_ui.NewText("SAMPLE", game_ui.TextStyle{
	Font: &titleFont,
})
var titleView = game_ui.NewView([]game_ui.Component{titleText}, game_ui.ViewStyle{Margin: "10 50 20"})

var startText = game_ui.NewText("START")
var startView = game_ui.NewView([]game_ui.Component{startText}, game_ui.ViewStyle{Margin: "10 50"})

var settingText = game_ui.NewText("SETTING")
var settingView = game_ui.NewView([]game_ui.Component{settingText}, game_ui.ViewStyle{Margin: "10 50"})

var exitText = game_ui.NewText("EXIT")
var exitView = game_ui.NewView([]game_ui.Component{exitText}, game_ui.ViewStyle{Margin: "10 50"})
