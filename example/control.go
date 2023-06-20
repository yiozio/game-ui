package main

import "github.com/hajimehoshi/ebiten/v2"

type KeySetting struct {
	Up     ebiten.Key
	Right  ebiten.Key
	Down   ebiten.Key
	Left   ebiten.Key
	Action ebiten.Key
}

type ButtonSetting struct {
	Up     ebiten.GamepadButton
	Right  ebiten.GamepadButton
	Down   ebiten.GamepadButton
	Left   ebiten.GamepadButton
	Action ebiten.GamepadButton
}

var (
	keySettings = KeySetting{
		Up:     ebiten.KeyW,
		Left:   ebiten.KeyA,
		Down:   ebiten.KeyS,
		Right:  ebiten.KeyD,
		Action: ebiten.KeySpace,
	}
	buttonSetting = ButtonSetting{
		Up:     ebiten.GamepadButton12,
		Left:   ebiten.GamepadButton14,
		Down:   ebiten.GamepadButton13,
		Right:  ebiten.GamepadButton15,
		Action: ebiten.GamepadButton0,
	}
)

type buttons []ebiten.GamepadButton

func (bts buttons) findIndex(button ebiten.GamepadButton) int {
	for i, b := range bts {
		if b == button {
			return i
		}
	}
	return -1
}
