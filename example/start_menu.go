package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"strconv"
	"github.com/yiozio/game-ui"
)

type startMenu struct {
	game_ui.Component
	selectedIndex int
	settingMenu   *settingMenu
	actionEffect  func(x, y int)
}

func NewStartMenu(actionEffect func(x, y int)) *startMenu {
	return &startMenu{game_ui.NewWindow([]game_ui.Component{game_ui.NewView([]game_ui.Component{
		titleView,
		startView,
		settingView,
		exitView,
	}, game_ui.ViewStyle{
		Width:            "640",
		Height:           "480",
		PositionVertical: game_ui.Center,
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
		var actionAreaMinPoint, actionAreaMaxPoint = view.GetActionArea()
		if actionAreaMinPoint.X <= mouseX && mouseX <= actionAreaMaxPoint.X &&
			actionAreaMinPoint.Y <= mouseY && mouseY <= actionAreaMaxPoint.Y {
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
		var actionAreaMinPoint, actionAreaMaxPoint = view.GetActionArea()
		if actionAreaMinPoint.X <= touchX && touchX <= actionAreaMaxPoint.X &&
			actionAreaMinPoint.Y <= touchY && touchY <= actionAreaMaxPoint.Y {
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

	var n = 0xe * ((now / 2) % 16)
	if n > (0xe * 8) {
		n = 0xd*16 - n
	}
	n -= 1
	var bColor1 = "#ffffff" + strconv.FormatInt(0x190+n, 16)[1:]
	var bColor2 = "#ffffff00"
	var bgColor1 = "#5599cc" + strconv.FormatInt(0x150+n, 16)[1:]
	var bgColor2 = "#5599cc00"

	// draw window
	for i := range startMenuItems {
		if i == m.selectedIndex && m.settingMenu == nil {
			startMenuItems[i].ReplaceStyle(0, game_ui.ViewStyle{
				BorderColor:     bColor1 + " " + bColor2 + " " + bColor2 + " " + bColor1,
				BackgroundColor: bgColor1 + " " + bgColor2 + " " + bgColor2 + " " + bgColor1,
			})
			ebitenutil.DebugPrintAt(screen, "on", 300, i*15)
		} else {
			startMenuItems[i].PopStyle()
			ebitenutil.DebugPrintAt(screen, "off", 300, i*15)
		}
	}

	m.Component.Draw(screen, 0, 0)
	for i, view := range startMenuItems {
		var actionAreaMinPoint, actionAreaMaxPoint = view.GetActionArea()
		ebitenutil.DebugPrintAt(screen,
			strconv.Itoa(actionAreaMinPoint.X)+","+
				strconv.Itoa(actionAreaMinPoint.Y)+":"+
				strconv.Itoa(actionAreaMaxPoint.X)+","+
				strconv.Itoa(actionAreaMaxPoint.Y), 100, i*15)
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
var titleView = game_ui.NewView([]game_ui.Component{titleText}, game_ui.ViewStyle{Margin: "10 50 20"})

var startMenuItemStyle = game_ui.ViewStyle{
	Margin:          "5 45",
	Width:           "200",
	Padding:         "2 40 1 10",
	BorderWidth:     "1 0 1 1",
	BorderColor:     "#00000000",
	BackgroundColor: "#00000000",
	Radius:          "20 0 0 20",
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
