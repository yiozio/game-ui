package start

import (
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/yiozio/game-ui"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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
var titleView = game_ui.NewView([]game_ui.Component{titleText}, game_ui.ViewStyle{Margin: game_ui.Size3(game_ui.Px(10), game_ui.Px(50), game_ui.Px(20))})

var startMenuItemStyle = game_ui.ViewStyle{
	Margin:          game_ui.Size2(game_ui.Px(5), game_ui.Px(45)),
	Width:           game_ui.Px(200),
	Padding:         game_ui.Size4(game_ui.Px(2), game_ui.Px(40), game_ui.Px(1), game_ui.Px(10)),
	BorderWidth:     game_ui.Size4(game_ui.Px(1), game_ui.Px(0), game_ui.Px(1), game_ui.Px(1)),
	BorderColor:     game_ui.ColorCode1(0x00000000),
	BackgroundColor: game_ui.ColorCode1(0x00000000),
	Radius:          game_ui.Radius4(20, 0, 0, 20),
}

var startText = game_ui.NewText("START")
var startView = game_ui.NewView([]game_ui.Component{startText}, startMenuItemStyle)

var settingText = game_ui.NewText("SETTING")
var settingView = game_ui.NewView([]game_ui.Component{settingText}, startMenuItemStyle)

var exitText = game_ui.NewText("EXIT")
var exitView = game_ui.NewView([]game_ui.Component{exitText}, startMenuItemStyle)

var startMenuItems = []game_ui.View{
	startView,
	settingView,
	exitView,
}
