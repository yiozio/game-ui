package setting

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yiozio/game-ui/example/control"
)

type Menu interface {
	Update(now int64, screenSize image.Point, mode control.Mode, enable bool)
	Draw(screen *ebiten.Image, now int64, screenSize image.Point, mode control.Mode)
}

var Opened Menu = nil
