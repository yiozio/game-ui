package start

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yiozio/game-ui/example/control"
	"github.com/yiozio/game-ui/example/menu/start"
	"image"
)

type Scene struct {
	instance *start.Menu
}

func NewScene() *Scene {
	return &Scene{
		instance: start.NewStartMenu(),
	}
}

func (s *Scene) Update(now int64, screenSize image.Point, mode control.Mode) {
	s.instance.Update(now, screenSize, mode, true)
}

func (s *Scene) Draw(screen *ebiten.Image, now int64, screenSize image.Point, mode control.Mode) {
	s.instance.Draw(screen, now, screenSize, mode)
}
