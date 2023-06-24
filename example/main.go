package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"math"
	"strconv"
	"time"
)

var (
	emptyImage = ebiten.NewImage(3, 3)

	// emptySubImage is an internal sub image of emptyImage.
	// Use emptySubImage at DrawTriangles instead of emptyImage in order to avoid bleeding edges.
	emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	emptyImage.Fill(color.White)
}

type Game struct {
	now                       int64
	control                   ControlMode
	justPressedTouchIds       []ebiten.TouchID
	gamepadIds                []ebiten.GamepadID
	gamepadJustPressedButtons buttons
	startMenu                 *startMenu
	settingMenu               *settingMenu
	selectedStartMenuIndex    int
	effectPos                 image.Point
	effectAt                  int64
}

type ControlMode = int

const (
	Mouse   ControlMode = 0
	Gamepad ControlMode = 1
	Touch   ControlMode = 2
)

func NewGame() ebiten.Game {
	var g = &Game{}
	g.init()
	return g
}

func (g *Game) init() {
	g.selectedStartMenuIndex = -1
	g.now = time.Now().UnixMilli()
	g.control = Mouse
	g.startMenu = NewStartMenu(
		g.Start,
		g.Setting,
		g.Exit,
		g.menuEffect,
		func() bool {
			return g.selectedStartMenuIndex >= 0
		})
}
func (g *Game) Start() {
	debugMessage = "START"
}
func (g *Game) Setting() {
	g.selectedStartMenuIndex = 1
	g.settingMenu = NewSettingMenu(func() { g.selectedStartMenuIndex = -1; g.settingMenu = nil }, g.menuEffect, func() bool { return g.selectedStartMenuIndex != 1 }, g.control)
	debugMessage = "SETTING"
}
func (g *Game) Exit() {
	debugMessage = "EXIT"
}

func (g *Game) menuEffect(x, y int) {
	g.effectPos = image.Point{X: x, Y: y}
	g.effectAt = g.now
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func (g *Game) Update() error {
	g.now = time.Now().UnixMilli()

	// get input
	var isMouseClicked = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	g.gamepadJustPressedButtons = []ebiten.GamepadButton{}
	g.gamepadIds = inpututil.AppendJustConnectedGamepadIDs(g.gamepadIds)
	if len(g.gamepadIds) > 0 {
		g.gamepadJustPressedButtons = inpututil.AppendJustPressedGamepadButtons(g.gamepadIds[0], nil)
	}
	g.justPressedTouchIds = inpututil.AppendJustPressedTouchIDs(nil)

	// set input
	if g.control == Mouse {
		var mouseX, mouseY = ebiten.CursorPosition()
		if mouseX != 0 || mouseY != 0 {
			g.startMenu.OnMouseMove(mouseX, mouseY, isMouseClicked)
			if g.settingMenu != nil {
				g.settingMenu.OnMouseMove(mouseX, mouseY, isMouseClicked)
			}
		}
	}
	if g.control == Touch {
		if len(g.justPressedTouchIds) > 0 {
			var touchX, touchY = ebiten.TouchPosition(g.justPressedTouchIds[0])
			g.startMenu.OnTouch(touchX, touchY)
			if g.settingMenu != nil {
				g.settingMenu.OnTouch(touchX, touchY)
			}
		}
	}
	if g.control == Gamepad {
		if g.gamepadJustPressedButtons.findIndex(buttonSetting.Up) >= 0 {
			g.startMenu.OnGamepadUp()
			if g.settingMenu != nil {
				g.settingMenu.OnGamepadUp()
			}
		} else if g.gamepadJustPressedButtons.findIndex(buttonSetting.Down) >= 0 {
			g.startMenu.OnGamepadDown()
			if g.settingMenu != nil {
				g.settingMenu.OnGamepadDown()
			}
		}
		if g.gamepadJustPressedButtons.findIndex(buttonSetting.Action) >= 0 {
			g.startMenu.OnGamepadAction()
			if g.settingMenu != nil {
				g.settingMenu.OnGamepadAction()
			}
		}
	}

	// switch control mode
	if g.control != Mouse && isMouseClicked {
		g.control = Mouse
		g.startMenu.ChangeControlMode(g.control)
		if g.settingMenu != nil {
			g.settingMenu.ChangeControlMode(g.control)
		}
	} else if g.control != Touch && len(g.justPressedTouchIds) > 0 {
		g.control = Touch
		g.startMenu.ChangeControlMode(g.control)
		if g.settingMenu != nil {
			g.settingMenu.ChangeControlMode(g.control)
		}
	} else if g.control != Gamepad &&
		(g.gamepadJustPressedButtons.findIndex(buttonSetting.Up) >= 0 ||
			g.gamepadJustPressedButtons.findIndex(buttonSetting.Down) >= 0 ||
			g.gamepadJustPressedButtons.findIndex(buttonSetting.Left) >= 0 ||
			g.gamepadJustPressedButtons.findIndex(buttonSetting.Right) >= 0 ||
			g.gamepadJustPressedButtons.findIndex(buttonSetting.Action) >= 0) {
		g.control = Gamepad
		g.startMenu.ChangeControlMode(g.control)
		if g.settingMenu != nil {
			g.settingMenu.ChangeControlMode(g.control)
		}
	}

	{ // debug
		if len(g.gamepadJustPressedButtons) > 0 {
			buttonInput = strconv.Itoa(int(g.gamepadJustPressedButtons[0]))
		}
	}

	return nil
}

var debugMessage = ""
var buttonInput = ""

func (g *Game) Draw(screen *ebiten.Image) {
	g.startMenu.Draw(screen, g.now)
	if g.selectedStartMenuIndex == 1 && g.settingMenu != nil {
		g.settingMenu.Draw(screen, g.now)
	}

	const msec = 200
	if (g.now - g.effectAt) < msec {
		var path = &vector.Path{}
		const maxR = 30
		var rate = float32(g.now-g.effectAt) / msec
		var r = rate * maxR
		path.MoveTo(0, 0)
		path.Arc(0, 0, r, 0, math.Pi*2+1, vector.Clockwise)
		path.Close()
		var vs, is = path.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{Width: 4 * (1 - rate)})
		for i := range vs {
			vs[i].DstX += float32(g.effectPos.X)
			vs[i].DstY += float32(g.effectPos.Y)
			vs[i].ColorR = 1
			vs[i].ColorG = 1
			vs[i].ColorB = 1
			vs[i].ColorA = 1 - rate
		}
		screen.DrawTriangles(vs, is, emptySubImage, &ebiten.DrawTrianglesOptions{
			FillRule: ebiten.EvenOdd,
		})
	}

	var mouseX, mouseY = ebiten.CursorPosition()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\nX:%d, Y:%d\n"+debugMessage+"\n"+"%d:%d:"+buttonInput, ebiten.ActualTPS(), ebiten.ActualFPS(), mouseX, mouseY, g.control, g.startMenu.selectedIndex))
}

func main() {
	ebiten.SetWindowSize(640*2, 480*2)
	if err := ebiten.RunGame(NewGame()); err != nil {
		panic(err)
	}
}
