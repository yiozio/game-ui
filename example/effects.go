package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"math"
)

func drawClickEffect(screen *ebiten.Image, g *Game) {
	const msec = 200
	if (g.now - g.effectAt) < msec {
		var path = &vector.Path{}
		const maxR = 30
		var rate = float32(g.now-g.effectAt) / msec
		var r = rate * maxR
		path.MoveTo(0, 0)
		path.Arc(0, 0, r, 0, math.Pi*2+1, vector.Clockwise)
		path.Close()
		var vs, is = path.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{Width: 4 * (1 - rate)})
		for i := range vs {
			vs[i].DstX += float32(g.effectPos.X)
			vs[i].DstY += float32(g.effectPos.Y)
			vs[i].ColorR = 1
			vs[i].ColorG = 1
			vs[i].ColorB = 1
			vs[i].ColorA = 1 - rate
		}
		screen.DrawTriangles(vs, is, emptySubImage, &ebiten.DrawTrianglesOptions{
			FillRule: ebiten.EvenOdd,
		})
	}
}
