package menu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yiozio/game-ui"
	"image"
)

func FindHoveredCell(cells []game_ui.View) (int, game_ui.View) {
	var x, y = ebiten.CursorPosition()
	for j, c := range cells {
		var area = c.Area()
		if CheckPointInArea(x, y, area.Min, area.Max) {
			return j, c
		}
	}

	return -1, nil
}

func CheckPointInArea(px, py int, minPoint, maxPoint image.Point) bool {
	return minPoint.X <= px && px <= maxPoint.X && minPoint.Y <= py && py <= maxPoint.Y
}
