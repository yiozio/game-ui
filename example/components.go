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

var titleText = game_ui.NewTextWithStyle("SAMPLE", game_ui.TextStyle{
	Font: &titleFont,
})
var titleView = game_ui.NewViewWithStyle([]game_ui.Component{titleText}, game_ui.ViewStyle{Margin: "10 50 20"})

var startText = game_ui.NewTextWithStyle("START", game_ui.TextStyle{})
var startView = game_ui.NewViewWithStyle([]game_ui.Component{startText}, game_ui.ViewStyle{
	Margin:          "5 45",
	Padding:         "2 40 1 10",
	BorderWidth:     "1 0 1 1",
	BorderColor:     "#ffffff99 #ffffff00 #ffffff00 #ffffff99",
	BackgroundColor: "#5599cc55 #5599cc00 #5599cc00 #5599cc55",
	Radius:          "20 0 0 20",
})

var settingText = game_ui.NewTextWithStyle("SETTING", game_ui.TextStyle{})
var settingView = game_ui.NewViewWithStyle([]game_ui.Component{settingText}, game_ui.ViewStyle{Margin: "10 50"})

var exitText = game_ui.NewTextWithStyle("EXIT", game_ui.TextStyle{})
var exitView = game_ui.NewViewWithStyle([]game_ui.Component{exitText}, game_ui.ViewStyle{Margin: "10 50"})

var menuView = game_ui.NewViewWithStyle([]game_ui.Component{
	titleView,
	startView,
	settingView,
	exitView,
}, game_ui.ViewStyle{
	Width:            "640",
	Height:           "480",
	PositionVertical: game_ui.Center,
})
