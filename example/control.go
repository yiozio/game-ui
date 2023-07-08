package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"strconv"
)

type KeySetting struct {
	Up     ebiten.Key
	Right  ebiten.Key
	Down   ebiten.Key
	Left   ebiten.Key
	Action ebiten.Key
}

type ButtonSetting struct {
	Up     ebiten.StandardGamepadButton
	Right  ebiten.StandardGamepadButton
	Down   ebiten.StandardGamepadButton
	Left   ebiten.StandardGamepadButton
	Action ebiten.StandardGamepadButton
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
		Up:     ebiten.StandardGamepadButtonLeftTop,
		Left:   ebiten.StandardGamepadButtonLeftLeft,
		Down:   ebiten.StandardGamepadButtonLeftBottom,
		Right:  ebiten.StandardGamepadButtonLeftRight,
		Action: ebiten.StandardGamepadButtonRightBottom,
	}
	control ControlMode = Mouse
)

type ControlMode = int

const (
	Mouse   ControlMode = 0
	Gamepad ControlMode = 1
	Touch   ControlMode = 2
)

var gamepadId *ebiten.GamepadID = nil

type buttons []ebiten.StandardGamepadButton

func (bts buttons) findIndex(button ebiten.StandardGamepadButton) int {
	for i, b := range bts {
		if b == button {
			return i
		}
	}
	return -1
}

func buttonToString(btn ebiten.StandardGamepadButton) string {
	if !ebiten.IsStandardGamepadButtonAvailable(*gamepadId, btn) {
		return "None"
	}

	switch btn {
	case ebiten.StandardGamepadButtonRightBottom:
		return "A"
	case ebiten.StandardGamepadButtonRightRight:
		return "B"
	case ebiten.StandardGamepadButtonRightLeft:
		return "X"
	case ebiten.StandardGamepadButtonRightTop:
		return "Y"
	case ebiten.StandardGamepadButtonFrontTopLeft:
		return "L1"
	case ebiten.StandardGamepadButtonFrontTopRight:
		return "R1"
	case ebiten.StandardGamepadButtonFrontBottomLeft:
		return "L2"
	case ebiten.StandardGamepadButtonFrontBottomRight:
		return "R2"
	case ebiten.StandardGamepadButtonCenterLeft:
		return "Opt1"
	case ebiten.StandardGamepadButtonCenterRight:
		return "Opt2"
	case ebiten.StandardGamepadButtonLeftStick:
		return "L3"
	case ebiten.StandardGamepadButtonRightStick:
		return "R3"
	case ebiten.StandardGamepadButtonLeftTop:
		return "↑"
	case ebiten.StandardGamepadButtonLeftBottom:
		return "↓"
	case ebiten.StandardGamepadButtonLeftLeft:
		return "←"
	case ebiten.StandardGamepadButtonLeftRight:
		return "→"
	case ebiten.StandardGamepadButtonCenterCenter:
		return "Opt3"
	}
	return strconv.Itoa(int(btn))
}
