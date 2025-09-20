package gamepad

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

func ButtonToString(btn ebiten.StandardGamepadButton) string {
	switch btn {
	case StandardGamepadButtonLeftStickUp:
		return "L3↑"
	case StandardGamepadButtonLeftStickRight:
		return "L3↓"
	case StandardGamepadButtonLeftStickDown:
		return "L3←"
	case StandardGamepadButtonLeftStickLeft:
		return "L3→"
	case StandardGamepadButtonRightStickUp:
		return "R3↑"
	case StandardGamepadButtonRightStickRight:
		return "R3↓"
	case StandardGamepadButtonRightStickDown:
		return "R3←"
	case StandardGamepadButtonRightStickLeft:
		return "R3→"
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
	case ebiten.StandardGamepadButtonCenterRight:
		return "Opt1"
	case ebiten.StandardGamepadButtonCenterLeft:
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
