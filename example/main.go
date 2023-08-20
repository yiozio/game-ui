package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yiozio/game-ui/example/control"
	"github.com/yiozio/game-ui/example/effect/action"
	"github.com/yiozio/game-ui/example/scene"
	"github.com/yiozio/game-ui/example/scene/start"
	"image"
	"time"
)

type Game struct {
	now        int64
	screenSize image.Point
	mode       control.Mode
}

func NewGame() ebiten.Game {
	var g = &Game{
		now:        time.Now().UnixMilli(),
		screenSize: image.Point{X: 640, Y: 480},
		mode:       control.Mouse,
	}
	return g
}

func (g *Game) Update() error {
	g.now = time.Now().UnixMilli()
	control.UpdateGamepadIds()
	g.mode = control.UpdateControlMode(g.mode)
	if scene.Instance == nil {
		scene.Instance = start.NewScene()
	}
	scene.Instance.Update(g.now, g.screenSize, g.mode)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	scene.Instance.Draw(screen, g.now, g.screenSize, g.mode)
	action.DrawEffect(screen, g.now)
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
