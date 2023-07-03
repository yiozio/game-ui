package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yiozio/game-ui"
	"strconv"
)

type startMenu struct {
	game_ui.Window
	selectedIndex int
	settingMenu   *settingMenu
	actionEffect  func(x, y int)
}

func NewStartMenu(actionEffect func(x, y int)) *startMenu {
	var pos = game_ui.Center
	return &startMenu{game_ui.NewWindow([]game_ui.Component{game_ui.NewView([]game_ui.Component{
		titleView,
		startView,
		settingView,
		exitView,
	}, game_ui.ViewStyle{
		Width:            game_ui.Vw(1),
		Height:           game_ui.Vh(1),
		PositionVertical: &pos,
	})}), 0, nil, actionEffect}
}

func (m *startMenu) start() {
	debugMessage = "START"
}

func (m *startMenu) setting() {
	m.selectedIndex = 1
	m.settingMenu = NewSettingMenu(
		func() { m.selectedIndex = -1; m.settingMenu = nil },
		m.actionEffect,
		func() bool { return m.selectedIndex != 1 },
		control,
	)
	debugMessage = "SETTING"
}

func (m *startMenu) exit() {
	debugMessage = "EXIT"
}

//// controls

func (m *startMenu) ChangeControlMode(mode ControlMode) {
	switch mode {
	case Mouse:
		m.selectedIndex = -1
	case Gamepad:
		m.selectedIndex = 0
	case Touch:
		m.selectedIndex = -1
	}
	if m.settingMenu != nil {
		m.settingMenu.ChangeControlMode(mode)
	}
}

func (m *startMenu) OnMouseMove(mouseX, mouseY int, justClick bool) {
	if m.settingMenu != nil {
		m.settingMenu.OnMouseMove(mouseX, mouseY, justClick)
		return
	}
	var hovered = false
	for i, view := range startMenuItems {
		var actionArea = view.Area()
		if actionArea.Min.X <= mouseX && mouseX <= actionArea.Max.X &&
			actionArea.Min.Y <= mouseY && mouseY <= actionArea.Max.Y {
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

func (m *startMenu) OnTouch(touchX, touchY int) {
	if m.settingMenu != nil {
		m.settingMenu.OnTouch(touchX, touchY)
		return
	}
	var action = false
	for i, view := range startMenuItems {
		var actionArea = view.Area()
		if actionArea.Min.X <= touchX && touchX <= actionArea.Max.X &&
			actionArea.Min.Y <= touchY && touchY <= actionArea.Max.Y {
			m.selectedIndex = i
			action = true
			break
		}
	}

	if action {
		m.onAction(touchX, touchY)
	}
}

func (m *startMenu) OnGamepadUp() {
	if m.settingMenu != nil {
		m.settingMenu.OnGamepadUp()
		return
	}
	m.selectedIndex -= 1
	if m.selectedIndex < 0 {
		m.selectedIndex = len(startMenuItems) - 1
	}
}

func (m *startMenu) OnGamepadDown() {
	if m.settingMenu != nil {
		m.settingMenu.OnGamepadDown()
		return
	}
	m.selectedIndex += 1
	if m.selectedIndex >= len(startMenuItems) {
		m.selectedIndex = 0
	}
}

func (m *startMenu) OnGamepadAction() {
	if m.settingMenu != nil {
		m.settingMenu.OnGamepadAction()
		return
	}
	m.onAction(-0xffff, -0xffff)
}

func (m *startMenu) onAction(x, y int) {
	switch m.selectedIndex {
	case 0:
		m.start()
	case 1:
		m.setting()
	case 2:
		m.exit()
	}
	m.actionEffect(x, y)
}

func (m *startMenu) Update(now int64) {
	if m.settingMenu != nil {
		m.settingMenu.Update(now)
	}
}

//// draw

func (m *startMenu) Draw(screen *ebiten.Image, now int64, screenSizeX, screenSizeY int) {
	// draw background
	{
		var path = &vector.Path{}
		path.MoveTo(0, 0)
		path.LineTo(640, 0)
		path.LineTo(640, 480)
		path.LineTo(0, 480)
		path.Close()

		var vs, is = path.AppendVerticesAndIndicesForFilling(nil, nil)
		for i := range vs {
			vs[i].ColorR = 0x1 / float32(0xf)
			vs[i].ColorG = 0x1 / float32(0xf)
			vs[i].ColorB = 0x1 / float32(0xf)
			vs[i].ColorA = 1
		}
		screen.DrawTriangles(vs, is, emptySubImage, &ebiten.DrawTrianglesOptions{
			FillRule: ebiten.EvenOdd,
		})
	}

	var n = 0xe * uint32((now/2)%16)
	if n > (0xe * 8) {
		n = 0xd*16 - n
	}
	n -= 1
	var bColor1 uint32 = 0xffffff90 + n
	var bColor2 uint32 = 0xffffff00
	var bgColor1 uint32 = 0x5599cc50 + n
	var bgColor2 uint32 = 0x5599cc00

	// draw window
	for i := range startMenuItems {
		if i == m.selectedIndex && m.settingMenu == nil {
			startMenuItems[i].ReplaceStyle(0, game_ui.ViewStyle{
				BorderColor:     game_ui.ColorCodeHorizontal(bColor1, bColor2),
				BackgroundColor: game_ui.ColorCodeHorizontal(bgColor1, bgColor2),
			})
			ebitenutil.DebugPrintAt(screen, "on", 300, i*15)
		} else {
			startMenuItems[i].PopStyle()
			ebitenutil.DebugPrintAt(screen, "off", 300, i*15)
		}
	}

	m.Window.Draw(screen, 0, 0)
	for i, view := range startMenuItems {
		var actionArea = view.Area()
		ebitenutil.DebugPrintAt(screen,
			strconv.Itoa(actionArea.Min.X)+","+
				strconv.Itoa(actionArea.Min.Y)+":"+
				strconv.Itoa(actionArea.Max.X)+","+
				strconv.Itoa(actionArea.Max.Y), 100, i*15)
	}

	if m.settingMenu != nil {
		m.settingMenu.Draw(screen, now, screenSizeX, screenSizeY)
	}
}

//// components

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
