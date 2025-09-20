package action

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yiozio/game-ui/example/def/common"
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
		var path = &vector.Path{}
		const maxR = 30
		var rate = float32(now-effectAt) / msec
		var r = rate * maxR
		path.MoveTo(0, 0)
		path.Arc(0, 0, r, 0, math.Pi*2+1, vector.Clockwise)
		path.Close()
		var vs, is = path.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{Width: 4 * (1 - rate)})
		for i := range vs {
			vs[i].DstX += float32(x)
			vs[i].DstY += float32(y)
			vs[i].ColorR = 1
			vs[i].ColorG = 1
			vs[i].ColorB = 1
			vs[i].ColorA = 1 - rate
		}
		screen.DrawTriangles(vs, is, common.EmptySubImage, &ebiten.DrawTrianglesOptions{
			FillRule: ebiten.FillRuleEvenOdd,
		})
	}
}
