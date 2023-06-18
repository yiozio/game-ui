package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"strconv"
	"yioz.io/game-ui"
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
	tickCount uint16
	control   ControlMode
}

type ControlMode = int

const (
	Mouse   ControlMode = 0
	Gamepad ControlMode = 1
	Touch   ControlMode = 2
)

func NewGame() ebiten.Game {
	var g = &Game{
		tickCount: 0,
		control:   Mouse,
	}
	g.init()
	return g
}

func (g *Game) init() {

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func (g *Game) Update() error {
	if g.tickCount == 0xffff {
		g.tickCount = 0
	} else {
		g.tickCount += 1
	}

	if g.control == Mouse {
		var mouseX, mouseY = ebiten.CursorPosition()
		if mouseX != 0 || mouseY != 0 {
			var index = -1
			for i, view := range startMenuItems {
				var actionAreaMinPoint, actionAreaMaxPoint = view.GetActionArea()
				if actionAreaMinPoint.X <= mouseX && mouseX <= actionAreaMaxPoint.X &&
					actionAreaMinPoint.Y <= mouseY && mouseY <= actionAreaMaxPoint.Y {
					index = i
					break
				}
			}
			selectedStartMenuItemIndex = index
			if index >= 0 && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				switch index {
				case 0:
					debugMessage = "START"
				case 1:
					debugMessage = "SETTING"
				case 2:
					debugMessage = "EXIT"
				}
			}
		}
	}
	if g.control == Touch {
	}
	if g.control == Gamepad {
	}

	return nil
}

var debugMessage = ""

func (g *Game) Draw(screen *ebiten.Image) {
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

	var m = int64(0xe * ((g.tickCount / 2) % 16))
	if m > (0xe * 8) {
		m = 0xd*16 - m
	}
	m -= 1
	var bColor1 = "#ffffff" + strconv.FormatInt(0x190+m, 16)[1:]
	var bColor2 = "#ffffff00"
	var bgColor1 = "#5599cc" + strconv.FormatInt(0x150+m, 16)[1:]
	var bgColor2 = "#5599cc00"

	// draw window
	for i := range startMenuItems {
		if i == selectedStartMenuItemIndex {
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

	w := game_ui.NewWindow([]game_ui.Component{game_ui.NewView([]game_ui.Component{
		titleView,
		startView,
		settingView,
		exitView,
	}, game_ui.ViewStyle{
		Width:            "640",
		Height:           "480",
		PositionVertical: game_ui.Center,
	})})
	w.Draw(screen, 0, 0)
	for i, view := range startMenuItems {
		var actionAreaMinPoint, actionAreaMaxPoint = view.GetActionArea()
		ebitenutil.DebugPrintAt(screen,
			strconv.Itoa(actionAreaMinPoint.X)+","+
				strconv.Itoa(actionAreaMinPoint.Y)+":"+
				strconv.Itoa(actionAreaMaxPoint.X)+","+
				strconv.Itoa(actionAreaMaxPoint.Y), 100, i*15)
	}
	var mouseX, mouseY = ebiten.CursorPosition()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\nX:%d, Y:%d\n"+debugMessage, ebiten.ActualTPS(), ebiten.ActualFPS(), mouseX, mouseY))
}

func main() {
	ebiten.SetWindowSize(640*2, 480*2)
	if err := ebiten.RunGame(NewGame()); err != nil {
		panic(err)
	}
}
