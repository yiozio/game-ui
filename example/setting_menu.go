package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yiozio/game-ui"
)

type settingMenu struct {
	game_ui.Window
	selectedIndex       int
	waitButtonIndex     int
	prevWaitButtonIndex int
	onExit              func()
	actionEffect        func(x, y int)
	isDisabled          func() bool
}

func toP[T string](str T) *T {
	return &str
}

func NewSettingMenu(onExit func(), actionEffect func(x, y int), isDisabled func() bool, mode ControlMode) *settingMenu {
	var selectedIndex = -1
	if mode == Gamepad {
		selectedIndex = 0
	}
	return &settingMenu{game_ui.NewWindow([]game_ui.Component{game_ui.NewView([]game_ui.Component{
		game_ui.NewView([]game_ui.Component{
			game_ui.NewView([]game_ui.Component{gamepadSettingMenuText[0]}, game_ui.ViewStyle{Width: game_ui.Px(60)}),
			game_ui.NewView([]game_ui.Component{gamepadSettingMenuText[1]}, game_ui.ViewStyle{Width: game_ui.Px(110), PositionHorizontal: toP(game_ui.Last)}),
		}, settingMenuItemStyle),
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
	})}), selectedIndex, -1, -1, onExit, actionEffect, isDisabled}
}

//// controls

func (m *settingMenu) ChangeControlMode(mode ControlMode) {
	switch mode {
	case Mouse:
		m.selectedIndex = -1
	case Gamepad:
		m.selectedIndex = 0
	case Touch:
		m.selectedIndex = -1
	}
}

func (m *settingMenu) OnMouseMove(mouseX, mouseY int, justClick bool) {
	if m.isDisabled() || m.waitButtonIndex >= 0 {
		return
	}
	var hovered = false
	for i, view := range settingMenuItems {
		var area = view.Area()
		if area.Min.X <= mouseX && mouseX <= area.Max.X &&
			area.Min.Y <= mouseY && mouseY <= area.Max.Y {
			m.selectedIndex = i
			hovered = true
			break
		}
	}

	if !hovered {
		m.selectedIndex = -1
	} else if justClick {
		m.onAction(mouseX, mouseY)
	}
}

func (m *settingMenu) OnTouch(touchX, touchY int) {
	if m.isDisabled() || m.waitButtonIndex >= 0 {
		return
	}
	var action = false
	for i, view := range settingMenuItems {
		var area = view.Area()
		if area.Min.X <= touchX && touchX <= area.Max.X &&
			area.Min.Y <= touchY && touchY <= area.Max.Y {
			m.selectedIndex = i
			action = true
			break
		}
	}

	if action {
		m.onAction(touchX, touchY)
	}
}

func (m *settingMenu) OnGamepadUp() {
	if m.isDisabled() || m.waitButtonIndex >= 0 {
		return
	}
	m.selectedIndex -= 1
	if m.selectedIndex < 0 {
		m.selectedIndex = len(settingMenuItems) - 1
	}
}

func (m *settingMenu) OnGamepadDown() {
	if m.isDisabled() || m.waitButtonIndex >= 0 {
		return
	}
	m.selectedIndex += 1
	if m.selectedIndex >= len(settingMenuItems) {
		m.selectedIndex = 0
	}
}

func (m *settingMenu) OnGamepadAction() {
	if m.isDisabled() || m.waitButtonIndex >= 0 {
		return
	}
	m.onAction(-0xffff, -0xffff)
}

func (m *settingMenu) onAction(x, y int) {
	switch m.selectedIndex {
	case 0:
		m.waitButtonIndex = 0
	case 1:
		m.waitButtonIndex = 1
	case 2:
		m.waitButtonIndex = 2
	case 3:
		m.onExit()
	}
	m.actionEffect(x, y)
}

func (m *settingMenu) Update(now int64) {
	if gamepadId != nil {
		var name = ebiten.GamepadName(*gamepadId)
		if len(name) > 17 {
			name = name[0:17] + "…"
		}
		gamepadSettingMenuText[1].ChangeText(name)
	} else {
		gamepadSettingMenuText[1].ChangeText("None")
	}
	gamepadUpSettingMenuText[1].ChangeText(buttonToString(buttonSetting.Up))
	gamepadDownSettingMenuText[1].ChangeText(buttonToString(buttonSetting.Down))
	gamepadActionSettingMenuText[1].ChangeText(buttonToString(buttonSetting.Action))

	if gamepadId != nil && m.prevWaitButtonIndex >= 0 {
		var buttons = inpututil.AppendJustPressedStandardGamepadButtons(*gamepadId, nil)
		debugMessage = "Waiting"
		if len(buttons) > 0 {
			switch m.waitButtonIndex {
			case 0:
				buttonSetting.Up = buttons[0]
			case 1:
				buttonSetting.Down = buttons[0]
			case 2:
				buttonSetting.Action = buttons[0]
			}
			m.waitButtonIndex = -1
		}
	}
	m.prevWaitButtonIndex = m.waitButtonIndex
}

//// draw

func (m *settingMenu) Draw(screen *ebiten.Image, now int64, screenSizeX, screenSizeY int) {
	// draw window
	for i := range settingMenuItems {
		if i == m.waitButtonIndex {
			settingMenuItems[i].ReplaceStyle(0, game_ui.ViewStyle{
				BorderColor: game_ui.ColorCode1(0xffffffff),
				BorderWidth: game_ui.Size4(game_ui.Px(0), game_ui.Px(0), game_ui.Px(1), game_ui.Px(10)),
			})
		} else if i == m.selectedIndex && !m.isDisabled() {
			settingMenuItems[i].ReplaceStyle(0, game_ui.ViewStyle{
				BorderColor: game_ui.ColorCode1(0xffffffff),
			})
		} else {
			settingMenuItems[i].PopStyle()
		}
	}
	var size = m.Window.GetSize()
	m.Window.Draw(screen, (screenSizeX-size.X-2)/2, (screenSizeY-size.Y-2)/2)
}

var settingMenuItemStyle = game_ui.ViewStyle{
	Margin:      game_ui.Size2(game_ui.Px(5), game_ui.Px(0)),
	Width:       game_ui.Px(50),
	Padding:     game_ui.Size4(game_ui.Px(2), game_ui.Px(5), game_ui.Px(0), game_ui.Px(5)),
	BorderWidth: game_ui.Size4(game_ui.Px(0), game_ui.Px(0), game_ui.Px(1), game_ui.Px(0)),
	BorderColor: game_ui.ColorCode1(0x00000000),
	Direction:   toP(game_ui.Horizontal),
}

var gamepadSettingMenuText = [2]game_ui.Text{game_ui.NewText("GAMEPAD: ", game_ui.TextStyle{Color: game_ui.Color(0x00aaaaff)}), game_ui.NewText("None", game_ui.TextStyle{Color: game_ui.Color(0x00aaaaff)})}
var gamepadUpSettingMenuText = [2]game_ui.Text{game_ui.NewText("GAMEPAD-UP: "), game_ui.NewText("")}
var gamepadDownSettingMenuText = [2]game_ui.Text{game_ui.NewText("GAMEPAD-DOWN: "), game_ui.NewText("")}
var gamepadActionSettingMenuText = [2]game_ui.Text{game_ui.NewText("GAMEPAD-ACTION: "), game_ui.NewText("")}
var closeSettingMenuText = game_ui.NewText("CLOSE")

var settingMenuItems []game_ui.View

func init() {
	for _, text := range [][2]game_ui.Text{gamepadUpSettingMenuText, gamepadDownSettingMenuText, gamepadActionSettingMenuText} {
		settingMenuItems = append(settingMenuItems, game_ui.NewView([]game_ui.Component{
			game_ui.NewView([]game_ui.Component{text[0]}, game_ui.ViewStyle{Width: game_ui.Px(140)}),
			game_ui.NewView([]game_ui.Component{text[1]}, game_ui.ViewStyle{Width: game_ui.Px(30), PositionHorizontal: toP(game_ui.Last)}),
		}, settingMenuItemStyle))
	}
	settingMenuItems = append(settingMenuItems, game_ui.NewView([]game_ui.Component{closeSettingMenuText}, settingMenuItemStyle))
}
