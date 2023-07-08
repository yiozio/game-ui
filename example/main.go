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
	gamepadJustPressedButtons buttons
	startMenu                 *startMenu
	effectPos                 image.Point
	effectAt                  int64
	screenSize                image.Point
}

func NewGame() ebiten.Game {
	var g = &Game{}
	g.init()
	return g
}

func (g *Game) init() {
	g.now = time.Now().UnixMilli()
	g.screenSize = image.Point{X: 640, Y: 480}
	g.startMenu = NewStartMenu(g.actionEffect)
}

func (g *Game) actionEffect(x, y int) {
	g.effectPos = image.Point{X: x, Y: y}
	g.effectAt = g.now
}

func (g *Game) Update() error {
	g.now = time.Now().UnixMilli()

	var isMouseClicked = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	g.gamepadJustPressedButtons = []ebiten.StandardGamepadButton{}
	var gamepadIds = inpututil.AppendJustConnectedGamepadIDs(nil)
	if len(gamepadIds) > 0 {
		gamepadId = &gamepadIds[0]
	}
	if gamepadId != nil && inpututil.IsGamepadJustDisconnected(*gamepadId) {
		gamepadId = nil
		control = Mouse
	}
	if gamepadId != nil {
		g.gamepadJustPressedButtons = inpututil.AppendJustPressedStandardGamepadButtons(*gamepadId, nil)
	}
	g.justPressedTouchIds = inpututil.AppendJustPressedTouchIDs(nil)

	if control == Mouse {
		var mouseX, mouseY = ebiten.CursorPosition()
		if mouseX != 0 || mouseY != 0 {
			g.startMenu.OnMouseMove(mouseX, mouseY, isMouseClicked)
		}
	}
	if control == Touch {
		if len(g.justPressedTouchIds) > 0 {
			var touchX, touchY = ebiten.TouchPosition(g.justPressedTouchIds[0])
			g.startMenu.OnTouch(touchX, touchY)
		}
	}
	if control == Gamepad {
		if g.gamepadJustPressedButtons.findIndex(buttonSetting.Up) >= 0 {
			g.startMenu.OnGamepadUp()
		} else if g.gamepadJustPressedButtons.findIndex(buttonSetting.Down) >= 0 {
			g.startMenu.OnGamepadDown()
		}
		if g.gamepadJustPressedButtons.findIndex(buttonSetting.Action) >= 0 {
			g.startMenu.OnGamepadAction()
		}
	}

	var _control = control
	if control != Mouse && isMouseClicked {
		_control = Mouse
	} else if control != Touch && len(g.justPressedTouchIds) > 0 {
		_control = Touch
	} else if control != Gamepad &&
		(len(g.gamepadJustPressedButtons) > 0) {
		_control = Gamepad
	}
	if _control != control {
		control = _control
		g.startMenu.ChangeControlMode(control)
	}

	g.startMenu.Update(g.now)

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
	g.startMenu.Draw(screen, g.now, g.screenSize.X, g.screenSize.Y)
	drawClickEffect(screen, g)
	var mouseX, mouseY = ebiten.CursorPosition()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\nX:%d, Y:%d\n"+debugMessage+"\n"+"%d:%d:"+buttonInput, ebiten.ActualTPS(), ebiten.ActualFPS(), mouseX, mouseY, control, g.startMenu.selectedIndex))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenSize.X, g.screenSize.Y
}

func main() {
	ebiten.SetWindowSize(640*2, 480*2)
	if err := ebiten.RunGame(NewGame()); err != nil {
		panic(err)
	}
}
