package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"strconv"
	"yioz.io/game-ui"
)

type settingMenu struct {
	game_ui.Component
	selectedIndex   int
	waitButtonIndex int
	onExit          func()
	actionEffect    func(x, y int)
	isDisabled      func() bool
}

func NewSettingMenu(onExit func(), actionEffect func(x, y int), isDisabled func() bool, mode ControlMode) *settingMenu {
	var selectedIndex = -1
	if mode == Gamepad {
		selectedIndex = 0
	}
	return &settingMenu{game_ui.NewWindow([]game_ui.Component{game_ui.NewView([]game_ui.Component{
		game_ui.NewView([]game_ui.Component{
			game_ui.NewView([]game_ui.Component{gamepadSettingMenuText[0]}, game_ui.ViewStyle{Width: "60"}),
			game_ui.NewView([]game_ui.Component{gamepadSettingMenuText[1]}, game_ui.ViewStyle{Width: "110", PositionHorizontal: game_ui.Last}),
		}, settingMenuItemStyle),
		settingMenuItems[0],
		settingMenuItems[1],
		settingMenuItems[2],
		settingMenuItems[3],
	}, game_ui.ViewStyle{
		BackgroundColor:  "#2255aa88",
		BorderWidth:      "2",
		BorderColor:      "#ffffff88",
		Radius:           "11",
		Padding:          "10 20",
		PositionVertical: game_ui.Center,
	})}), selectedIndex, -1, onExit, actionEffect, isDisabled}
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
		var min, max = view.GetActionArea()
		if min.X <= mouseX && mouseX <= max.X &&
			min.Y <= mouseY && mouseY <= max.Y {
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
		var min, max = view.GetActionArea()
		if min.X <= touchX && touchX <= max.X &&
			min.Y <= touchY && touchY <= max.Y {
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
	gamepadUpSettingMenuText[1].ChangeText(strconv.Itoa(int(buttonSetting.Up)))
	gamepadDownSettingMenuText[1].ChangeText(strconv.Itoa(int(buttonSetting.Down)))
	gamepadActionSettingMenuText[1].ChangeText(strconv.Itoa(int(buttonSetting.Action)))

	if gamepadId != nil && m.waitButtonIndex >= 0 && !actioned {
		var buttons = inpututil.AppendJustPressedGamepadButtons(*gamepadId, nil)
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
}

//// draw

func (m *settingMenu) Draw(screen *ebiten.Image, now int64, screenSizeX, screenSizeY int) {
	// draw window
	for i := range settingMenuItems {
		if i == m.waitButtonIndex {
			settingMenuItems[i].ReplaceStyle(0, game_ui.ViewStyle{
				BorderColor: "#ffffffff",
				BorderWidth: "0 0 1 10",
			})
		} else if i == m.selectedIndex && !m.isDisabled() {
			settingMenuItems[i].ReplaceStyle(0, game_ui.ViewStyle{
				BorderColor: "#ffffffff",
			})
		} else {
			settingMenuItems[i].PopStyle()
		}
	}
	var size = m.Component.GetSize()
	m.Component.Draw(screen, (screenSizeX-size.X-2)/2, (screenSizeY-size.Y-2)/2)
}

var settingMenuItemStyle = game_ui.ViewStyle{
	Margin:      "5 0",
	Width:       "50",
	Padding:     "2 5 0 5",
	BorderWidth: "0 0 1 0",
	BorderColor: "#00000000",
	Direction:   game_ui.Horizontal,
}

var gamepadSettingMenuText = [2]game_ui.Text{game_ui.NewText("GAMEPAD: ", game_ui.TextStyle{Color: "#0aa"}), game_ui.NewText("None", game_ui.TextStyle{Color: "#0aa"})}
var gamepadUpSettingMenuText = [2]game_ui.Text{game_ui.NewText("GAMEPAD-UP: "), game_ui.NewText("")}
var gamepadDownSettingMenuText = [2]game_ui.Text{game_ui.NewText("GAMEPAD-DOWN: "), game_ui.NewText("")}
var gamepadActionSettingMenuText = [2]game_ui.Text{game_ui.NewText("GAMEPAD-ACTION: "), game_ui.NewText("")}
var closeSettingMenuText = game_ui.NewText("CLOSE")

var settingMenuItems []game_ui.View

func init() {
	for _, text := range [][2]game_ui.Text{gamepadUpSettingMenuText, gamepadDownSettingMenuText, gamepadActionSettingMenuText} {
		settingMenuItems = append(settingMenuItems, game_ui.NewView([]game_ui.Component{
			game_ui.NewView([]game_ui.Component{text[0]}, game_ui.ViewStyle{Width: "140"}),
			game_ui.NewView([]game_ui.Component{text[1]}, game_ui.ViewStyle{Width: "30", PositionHorizontal: game_ui.Last}),
		}, settingMenuItemStyle))
	}
	settingMenuItems = append(settingMenuItems, game_ui.NewView([]game_ui.Component{closeSettingMenuText}, settingMenuItemStyle))
}
