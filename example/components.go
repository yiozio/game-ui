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

var menuItemStyle = game_ui.ViewStyle{
	Margin:          "5 45",
	Width:           "200",
	Padding:         "2 40 1 10",
	BorderWidth:     "1 0 1 1",
	BorderColor:     "#00000000",
	BackgroundColor: "#00000000",
	Radius:          "20 0 0 20",
}

var startText = game_ui.NewText("START")
var startView = game_ui.NewView([]game_ui.Component{startText}, menuItemStyle)

var settingText = game_ui.NewText("SETTING")
var settingView = game_ui.NewView([]game_ui.Component{settingText}, menuItemStyle)

var exitText = game_ui.NewText("EXIT")
var exitView = game_ui.NewView([]game_ui.Component{exitText}, menuItemStyle)

var startMenuItems = []game_ui.View{
	startView,
	settingView,
	exitView,
}

var selectedStartMenuItemIndex = -1
