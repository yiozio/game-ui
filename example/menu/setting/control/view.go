package control

import "github.com/yiozio/game-ui"

var settingWindow = game_ui.NewView([]game_ui.Component{
	gamepadSettingMenuView,
	settingMenuItems[0],
	settingMenuItems[1],
	settingMenuItems[2],
	settingMenuItems[3],
}, game_ui.ViewStyle{
	BackgroundColor:  game_ui.ColorCode1(0x2255aa88),
	BorderWidth:      game_ui.Size1(game_ui.Px(2)),
	BorderColor:      game_ui.ColorCode1(0xffffff88),
	Radius:           game_ui.Radius1(11),
	Padding:          game_ui.Size2(game_ui.Px(10), game_ui.Px(20)),
	PositionVertical: toP(game_ui.Center),
})

var settingMenuItemStyle = game_ui.ViewStyle{
	Margin:      game_ui.Size2(game_ui.Px(5), game_ui.Px(0)),
	Width:       game_ui.Px(50),
	Padding:     game_ui.Size4(game_ui.Px(2), game_ui.Px(5), game_ui.Px(0), game_ui.Px(5)),
	BorderWidth: game_ui.Size4(game_ui.Px(0), game_ui.Px(0), game_ui.Px(1), game_ui.Px(0)),
	BorderColor: game_ui.ColorCode1(0x00000000),
	Direction:   toP(game_ui.Horizontal),
}

var gamepadSettingMenuKeyText = game_ui.NewText("GAMEPAD: ", game_ui.TextStyle{Color: game_ui.Color(0x00aaaaff)})
var gamepadSettingMenuValueText = game_ui.NewText("None", game_ui.TextStyle{Color: game_ui.Color(0x00aaaaff)})
var gamepadUpSettingMenuKeyText = game_ui.NewText("GAMEPAD-UP: ")
var gamepadUpSettingMenuValueText = game_ui.NewText("")
var gamepadDownSettingMenuKeyText = game_ui.NewText("GAMEPAD-DOWN: ")
var gamepadDownSettingMenuValueText = game_ui.NewText("")
var gamepadActionSettingMenuKeyText = game_ui.NewText("GAMEPAD-ACTION: ")
var gamepadActionSettingMenuValueText = game_ui.NewText("")

var gamepadSettingMenuView = game_ui.NewView([]game_ui.Component{
	game_ui.NewView([]game_ui.Component{gamepadSettingMenuKeyText}, game_ui.ViewStyle{Width: game_ui.Px(60)}),
	game_ui.NewView([]game_ui.Component{gamepadSettingMenuValueText}, game_ui.ViewStyle{Width: game_ui.Px(110), PositionHorizontal: toP(game_ui.Last)}),
}, settingMenuItemStyle)

func menuTitleView(components []game_ui.Component) game_ui.View {
	return game_ui.NewView([]game_ui.Component{
		game_ui.NewView([]game_ui.Component{components[0]}, game_ui.ViewStyle{Width: game_ui.Px(140)}),
		game_ui.NewView([]game_ui.Component{components[1]}, game_ui.ViewStyle{Width: game_ui.Px(30), PositionHorizontal: toP(game_ui.Last)}),
	}, settingMenuItemStyle)
}

var gamepadUpSettingMenuView = menuTitleView([]game_ui.Component{gamepadUpSettingMenuKeyText, gamepadUpSettingMenuValueText})
var gamepadDownSettingMenuView = menuTitleView([]game_ui.Component{gamepadDownSettingMenuKeyText, gamepadDownSettingMenuValueText})
var gamepadActionSettingMenuView = menuTitleView([]game_ui.Component{gamepadActionSettingMenuKeyText, gamepadActionSettingMenuValueText})
var closeSettingMenuText = game_ui.NewText("CLOSE")

var settingMenuItems = []game_ui.View{
	gamepadUpSettingMenuView,
	gamepadDownSettingMenuView,
	gamepadActionSettingMenuView,
	game_ui.NewView([]game_ui.Component{closeSettingMenuText}, settingMenuItemStyle),
}

var selectedMenuItemStyle = game_ui.ViewStyle{
	BorderColor: game_ui.ColorCode1(0xffffffff),
}
