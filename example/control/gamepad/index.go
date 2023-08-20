package gamepad

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type ButtonMapping struct {
	Up     ebiten.StandardGamepadButton
	Down   ebiten.StandardGamepadButton
	Action ebiten.StandardGamepadButton
}

var defaultButtonMapping = ButtonMapping{
	Up:     ebiten.StandardGamepadButtonLeftTop,
	Down:   ebiten.StandardGamepadButtonLeftBottom,
	Action: ebiten.StandardGamepadButtonRightBottom,
}

var CurrentButtonMapping = defaultButtonMapping

const (
	StandardGamepadButtonLeftStickUp     ebiten.StandardGamepadButton = -1
	StandardGamepadButtonLeftStickRight  ebiten.StandardGamepadButton = -2
	StandardGamepadButtonLeftStickDown   ebiten.StandardGamepadButton = -3
	StandardGamepadButtonLeftStickLeft   ebiten.StandardGamepadButton = -4
	StandardGamepadButtonRightStickUp    ebiten.StandardGamepadButton = -5
	StandardGamepadButtonRightStickRight ebiten.StandardGamepadButton = -6
	StandardGamepadButtonRightStickDown  ebiten.StandardGamepadButton = -7
	StandardGamepadButtonRightStickLeft  ebiten.StandardGamepadButton = -8
)
