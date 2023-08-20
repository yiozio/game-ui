package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yiozio/game-ui/example/control"
	"image"
)

type Scene interface {
	Update(now int64, screenSize image.Point, mode control.Mode)
	Draw(screen *ebiten.Image, now int64, screenSize image.Point, mode control.Mode)
}

var Instance Scene = nil
