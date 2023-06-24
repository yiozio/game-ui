package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"strconv"
	"yioz.io/game-ui"
)

type settingMenu struct {
	game_ui.Component
	selectedIndex int
	onExit        func()
	actionEffect  func(x, y int)
	isDisabled    func() bool
}

const (
	settingMenuWidth  = 300
	settingMenuHeight = 200
)

func NewSettingMenu(onExit func(), actionEffect func(x, y int), isDisabled func() bool, mode ControlMode) *settingMenu {
	var selectedIndex = -1
	if mode == Gamepad {
		selectedIndex = 0
	}
	return &settingMenu{game_ui.NewWindow([]game_ui.Component{game_ui.NewView([]game_ui.Component{
		closeSettingMenuView,
	}, game_ui.ViewStyle{
		Width:            strconv.Itoa(settingMenuWidth),
		Height:           strconv.Itoa(settingMenuHeight),
		BackgroundColor:  "#2255aa88",
		BorderWidth:      "2",
		BorderColor:      "#ffffff88",
		Radius:           "11",
		PositionVertical: game_ui.Center,
	})}), selectedIndex, onExit, actionEffect, isDisabled}
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
	if m.isDisabled() {
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
	if m.isDisabled() {
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
	if m.isDisabled() {
		return
	}
	m.selectedIndex -= 1
	if m.selectedIndex < 0 {
		m.selectedIndex = len(settingMenuItems) - 1
	}
}

func (m *settingMenu) OnGamepadDown() {
	if m.isDisabled() {
		return
	}
	m.selectedIndex += 1
	if m.selectedIndex >= len(settingMenuItems) {
		m.selectedIndex = 0
	}
}

func (m *settingMenu) OnGamepadAction() {
	if m.isDisabled() {
		return
	}
	m.onAction(-0xffff, -0xffff)
}

func (m *settingMenu) onAction(x, y int) {
	switch m.selectedIndex {
	case 0:
		m.onExit()
	}
	m.actionEffect(x, y)
}

//// draw

func (m *settingMenu) Draw(screen *ebiten.Image, now int64) {
	// draw window
	for i := range settingMenuItems {
		if i == m.selectedIndex && !m.isDisabled() {
			settingMenuItems[i].ReplaceStyle(0, game_ui.ViewStyle{
				BorderColor: "#ffffffff",
			})
			ebitenutil.DebugPrintAt(screen, "on", 300, i*15)
		} else {
			settingMenuItems[i].PopStyle()
			ebitenutil.DebugPrintAt(screen, "off", 300, i*15)
		}
	}
	m.Component.Draw(screen, (640-settingMenuWidth-2)/2, (480-settingMenuHeight-2)/2)
	var min, max = closeSettingMenuView.GetActionArea()
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d,%d:%d,%d", min.X, min.Y, max.X, max.Y), 0, 100)
}

var settingMenuItemStyle = game_ui.ViewStyle{
	Margin:      "5 20",
	Width:       "50",
	Padding:     "2 5 0 5",
	BorderWidth: "0 0 1 0",
	BorderColor: "#00000000",
}

var closeSettingMenuText = game_ui.NewText("CLOSE")
var closeSettingMenuView = game_ui.NewView([]game_ui.Component{closeSettingMenuText}, settingMenuItemStyle)

var settingMenuItems = []game_ui.View{
	closeSettingMenuView,
}
