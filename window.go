package game_ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type windowComponent struct {
	components []Component
}
type Window = *windowComponent

func NewWindow(components []Component) Window {
	return &windowComponent{components}
}

func (w Window) GetSize() image.Point {
	var x, y = 0, 0
	for _, component := range w.components {
		var contentSize = component.GetSize()
		if x <= contentSize.X {
			x = contentSize.X
		}
		y += contentSize.Y
	}
	return image.Point{X: x, Y: y}
}

func (w Window) Draw(screen *ebiten.Image, x, y int) {
	var _y = 0
	for _, component := range w.components {
		component.Draw(screen, x, y+_y)
		_y += component.GetSize().Y
	}
}

func (w Window) IsFloating() bool {
	return false
}

func (w Window) Components() []Component {
	return w.components
}
