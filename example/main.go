package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"image/color"
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
	justPressedTouchIds       []ebiten.TouchID
	gamepadIds                []ebiten.GamepadID
	gamepadJustPressedButtons buttons
	startMenu                 *startMenu
	settingMenu               *settingMenu
	selectedStartMenuIndex    int
	effectPos                 image.Point
	effectAt                  int64
}

func NewGame() ebiten.Game {
	var g = &Game{}
	g.init()
	return g
}

func (g *Game) init() {
	g.selectedStartMenuIndex = -1
	g.now = time.Now().UnixMilli()
	g.startMenu = NewStartMenu(
		func() {
			debugMessage = "START"
		},
		func() {
			g.selectedStartMenuIndex = 1
			g.settingMenu = NewSettingMenu(
				func() { g.selectedStartMenuIndex = -1; g.settingMenu = nil },
				g.menuEffect,
				func() bool { return g.selectedStartMenuIndex != 1 },
				control,
			)
			debugMessage = "SETTING"
		},
		func() {
			debugMessage = "EXIT"
		},
		g.menuEffect,
		func() bool {
			return g.selectedStartMenuIndex >= 0
		})
}

func (g *Game) menuEffect(x, y int) {
	g.effectPos = image.Point{X: x, Y: y}
	g.effectAt = g.now
}

func (g *Game) Update() error {
	g.now = time.Now().UnixMilli()
	actioned = false

	// get input
	var isMouseClicked = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	g.gamepadJustPressedButtons = []ebiten.GamepadButton{}
	g.gamepadIds = inpututil.AppendJustConnectedGamepadIDs(g.gamepadIds)
	if len(g.gamepadIds) > 0 {
		g.gamepadJustPressedButtons = inpututil.AppendJustPressedGamepadButtons(g.gamepadIds[0], nil)
	}
	g.justPressedTouchIds = inpututil.AppendJustPressedTouchIDs(nil)

	// set input
	if control == Mouse {
		var mouseX, mouseY = ebiten.CursorPosition()
		if mouseX != 0 || mouseY != 0 {
			g.startMenu.OnMouseMove(mouseX, mouseY, isMouseClicked)
			if g.settingMenu != nil {
				g.settingMenu.OnMouseMove(mouseX, mouseY, isMouseClicked)
			}
		}
	}
	if control == Touch {
		if len(g.justPressedTouchIds) > 0 {
			var touchX, touchY = ebiten.TouchPosition(g.justPressedTouchIds[0])
			g.startMenu.OnTouch(touchX, touchY)
			if g.settingMenu != nil {
				g.settingMenu.OnTouch(touchX, touchY)
			}
		}
	}
	if control == Gamepad {
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
			actioned = true
			if g.settingMenu == nil {
				g.startMenu.OnGamepadAction()
			} else if g.settingMenu != nil {
				g.settingMenu.OnGamepadAction()
			}
		}
	}

	// switch control mode
	if control != Mouse && isMouseClicked {
		control = Mouse
		g.startMenu.ChangeControlMode(control)
		if g.settingMenu != nil {
			g.settingMenu.ChangeControlMode(control)
		}
	} else if control != Touch && len(g.justPressedTouchIds) > 0 {
		control = Touch
		g.startMenu.ChangeControlMode(control)
		if g.settingMenu != nil {
			g.settingMenu.ChangeControlMode(control)
		}
	} else if control != Gamepad &&
		(g.gamepadJustPressedButtons.findIndex(buttonSetting.Up) >= 0 ||
			g.gamepadJustPressedButtons.findIndex(buttonSetting.Down) >= 0 ||
			g.gamepadJustPressedButtons.findIndex(buttonSetting.Left) >= 0 ||
			g.gamepadJustPressedButtons.findIndex(buttonSetting.Right) >= 0 ||
			g.gamepadJustPressedButtons.findIndex(buttonSetting.Action) >= 0) {
		control = Gamepad
		g.startMenu.ChangeControlMode(control)
		if g.settingMenu != nil {
			g.settingMenu.ChangeControlMode(control)
		}
	}

	if g.settingMenu != nil {
		var gamepadId *ebiten.GamepadID = nil
		if len(g.gamepadIds) > 0 {
			gamepadId = &g.gamepadIds[0]
		}
		g.settingMenu.Update(g.now, gamepadId)
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

	drawClickEffect(screen, g)

	var mouseX, mouseY = ebiten.CursorPosition()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\nX:%d, Y:%d\n"+debugMessage+"\n"+"%d:%d:"+buttonInput, ebiten.ActualTPS(), ebiten.ActualFPS(), mouseX, mouseY, control, g.startMenu.selectedIndex))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640*2, 480*2)
	if err := ebiten.RunGame(NewGame()); err != nil {
		panic(err)
	}
}
