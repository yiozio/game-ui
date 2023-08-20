package control

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Mode = int

const (
	Mouse Mode = iota
	Touch
	Gamepad
)

func UpdateControlMode(current Mode) Mode {
	if current != Mouse && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return Mouse
	} else if touchedIDs := inpututil.AppendJustPressedTouchIDs(nil); len(touchedIDs) > 0 && current != Touch {
		return Touch
	} else if gid := GetGamepadId(); gid != nil && current != Gamepad && len(inpututil.AppendJustPressedStandardGamepadButtons(*gid, nil)) > 0 {
		return Gamepad
	}
	return current
}

var gamepadIds []ebiten.GamepadID

func UpdateGamepadIds() {
	for i, gid := range gamepadIds {
		if inpututil.IsGamepadJustDisconnected(gid) {
			gamepadIds = append(gamepadIds[:i], gamepadIds[i+1:]...)
		}
	}
	gamepadIds = inpututil.AppendJustConnectedGamepadIDs(gamepadIds)
}

func GetGamepadId() *ebiten.GamepadID {
	if len(gamepadIds) == 0 {
		return nil
	}
	return &gamepadIds[0]
}

type buttonLabel struct {
	Up     string
	Down   string
	Action string
}

var ButtonLabel = ButtonLabelEn

var ButtonLabelEn = buttonLabel{
	Up:     "Up",
	Down:   "Down",
	Action: "Action",
}

var ButtonLabelJp = buttonLabel{
	Up:     "上",
	Down:   "下",
	Action: "アクション",
}
