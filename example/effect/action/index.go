package action

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var x, y int = 0, 0
var effectAt int64 = 0

func StartEffect(now int64) {
	x, y = ebiten.CursorPosition()
	effectAt = now
}

func DrawEffect(screen *ebiten.Image, now int64) {
	const msec = 200
	if (now - effectAt) < msec {
		const maxR = 30
		var rate = float32(now-effectAt) / msec
		var r = rate * maxR
		vector.StrokeCircle(screen, float32(x), float32(y), r, 4*(1-rate), color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: uint8(0xff * rate)}, false)
	}
}
