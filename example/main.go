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

type Game struct{}

func NewGame() ebiten.Game {
	var g = &Game{}
	g.init()
	return g
}

func (g *Game) init() {

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func (g *Game) Update() error {
	return nil
}

var n = 0

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

	var m = n
	if m > 0x6f {
		m = (0x6f * 2) - n
	}
	var bColor1 = "#ffffff" + strconv.FormatInt(0x190+int64(m), 16)[1:]
	var bColor2 = "#ffffff00"
	var bgColor1 = "#5599cc" + strconv.FormatInt(0x150+int64(m), 16)[1:]
	var bgColor2 = "#5599cc00"
	if n >= 0xc0 {
		n = 0
	} else {
		n += 0xd
	}
	if startView.GetStylesCount() > 0 {
		startView.PopStyle()
	} else {
		startView.PushStyle(game_ui.ViewStyle{
			Margin:          "5 45",
			Padding:         "2 40 1 10",
			BorderWidth:     "1 0 1 1",
			BorderColor:     "#ffffff99 #ffffff00 #ffffff00 #ffffff99",
			BackgroundColor: "#5599cc55 #5599cc00 #5599cc00 #5599cc55",
			Radius:          "20 0 0 20",
		})
	}
	startView.PushStyle(game_ui.ViewStyle{
		BorderColor:     bColor1 + " " + bColor2 + " " + bColor2 + " " + bColor1,
		BackgroundColor: bgColor1 + " " + bgColor2 + " " + bgColor2 + " " + bgColor1,
	})

	// draw window
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
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func main() {
	ebiten.SetWindowSize(640*2, 480*2)
	if err := ebiten.RunGame(NewGame()); err != nil {
		panic(err)
	}
}
