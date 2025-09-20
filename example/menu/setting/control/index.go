package control

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yiozio/game-ui"
	"github.com/yiozio/game-ui/example/control"
	"github.com/yiozio/game-ui/example/control/gamepad"
	actionEffect "github.com/yiozio/game-ui/example/effect/action"
	"github.com/yiozio/game-ui/example/menu"
	"github.com/yiozio/game-ui/example/menu/setting"
)

type Menu struct {
	game_ui.Window
	buttonMapping      gamepad.ButtonMapping
	selectedMenuIndex  int
	inputWaitMenuIndex int
	initialized        bool
}

func toP[T string](str T) *T {
	return &str
}

func NewSettingMenu() *Menu {
	return &Menu{game_ui.NewWindow([]game_ui.Component{settingWindow}), gamepad.CurrentButtonMapping, -1, -1, false}
}

func (m *Menu) Update(now int64, screenSize image.Point, mode control.Mode, enable bool) {
	if mode == control.Gamepad && m.selectedMenuIndex < 0 {
		m.selectedMenuIndex = 0
	}
	if !m.initialized {
		gamepadUpSettingMenuKeyText.ChangeText(control.ButtonLabel.Up)
		gamepadUpSettingMenuValueText.ChangeText(gamepad.ButtonToString(m.buttonMapping.Up))
		gamepadDownSettingMenuKeyText.ChangeText(control.ButtonLabel.Down)
		gamepadDownSettingMenuValueText.ChangeText(gamepad.ButtonToString(m.buttonMapping.Down))
		gamepadActionSettingMenuKeyText.ChangeText(control.ButtonLabel.Action)
		gamepadActionSettingMenuValueText.ChangeText(gamepad.ButtonToString(m.buttonMapping.Action))
	}

	if gid := control.GetGamepadId(); gid != nil {
		var name = ebiten.GamepadName(*gid)
		if len(name) > 17 {
			name = name[0:17] + "…"
		}
		gamepadSettingMenuValueText.ChangeText(name)
	} else {
		gamepadSettingMenuValueText.ChangeText("None")
	}

	if m.inputWaitMenuIndex >= 0 {
		inputWait(m)
	} else {
		controlMenu(m, mode, now)
	}

	m.initialized = true
}

func (m *Menu) Draw(screen *ebiten.Image, now int64, screenSize image.Point, mode control.Mode) {
	// draw window
	for i := range settingMenuItems {
		if i == m.inputWaitMenuIndex {
			settingMenuItems[i].ReplaceStyle(0, game_ui.ViewStyle{
				BorderColor: game_ui.ColorCode1(0xffffffff),
				BorderWidth: game_ui.Size4(game_ui.Px(0), game_ui.Px(0), game_ui.Px(1), game_ui.Px(10)),
			})
		} else if i == m.selectedMenuIndex {
			settingMenuItems[i].ReplaceStyle(0, selectedMenuItemStyle)
		} else {
			settingMenuItems[i].ReplaceStyle(0, game_ui.ViewStyle{})
		}
	}
	var size = m.Window.GetSize()
	m.Window.Draw(screen, (screenSize.X-size.X-2)/2, (screenSize.Y-size.Y-2)/2)
}

func inputWait(m *Menu) {
	if gid := control.GetGamepadId(); gid != nil && m.initialized {
		var buttons = inpututil.AppendJustPressedStandardGamepadButtons(*gid, nil)
		if len(buttons) > 0 {
			switch m.selectedMenuIndex {
			case 0:
				m.buttonMapping.Up = buttons[0]
				gamepadUpSettingMenuValueText.ChangeText(gamepad.ButtonToString(m.buttonMapping.Up))
			case 1:
				m.buttonMapping.Down = buttons[0]
				gamepadDownSettingMenuValueText.ChangeText(gamepad.ButtonToString(m.buttonMapping.Down))
			case 2:
				m.buttonMapping.Action = buttons[0]
				gamepadActionSettingMenuValueText.ChangeText(gamepad.ButtonToString(m.buttonMapping.Action))
			}
			m.inputWaitMenuIndex = -1
		}
	} else {
		m.inputWaitMenuIndex = -1
	}
}

func controlMenu(m *Menu, mode control.Mode, now int64) {
	if mode != control.Gamepad {
		if m.selectedMenuIndex >= 0 {
			settingMenuItems[m.selectedMenuIndex].ReplaceStyle(0, game_ui.ViewStyle{})
			m.selectedMenuIndex = -1
		}
	}
	var action = false
	if mode == control.Mouse {
		m.selectedMenuIndex, _ = menu.FindHoveredCell(settingMenuItems)
		action = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	} else if mode == control.Touch {
	} else if mode == control.Gamepad {
		if m.selectedMenuIndex < 0 {
			m.selectedMenuIndex = 0
		}

		if gid := control.GetGamepadId(); gid != nil && m.initialized {
			if inpututil.IsStandardGamepadButtonJustPressed(*gid, m.buttonMapping.Up) {
				settingMenuItems[m.selectedMenuIndex].ReplaceStyle(0, game_ui.ViewStyle{})
				m.selectedMenuIndex += len(settingMenuItems) - 1
			} else if inpututil.IsStandardGamepadButtonJustPressed(*gid, m.buttonMapping.Down) {
				settingMenuItems[m.selectedMenuIndex].ReplaceStyle(0, game_ui.ViewStyle{})
				m.selectedMenuIndex += 1
			}
			m.selectedMenuIndex = m.selectedMenuIndex % len(settingMenuItems)

			action = inpututil.IsStandardGamepadButtonJustPressed(*gid, m.buttonMapping.Action)
		}
	}

	if m.selectedMenuIndex >= 0 {
		var view = settingMenuItems[m.selectedMenuIndex]
		view.ReplaceStyle(0, selectedMenuItemStyle)
		if action {
			if mode == control.Mouse {
				actionEffect.StartEffect(now)
			}
			if m.selectedMenuIndex < 3 {
				m.inputWaitMenuIndex = m.selectedMenuIndex
			}
			switch m.selectedMenuIndex {
			case 0:
				gamepadUpSettingMenuValueText.ChangeText("...")
			case 1:
				gamepadDownSettingMenuValueText.ChangeText("...")
			case 2:
				gamepadActionSettingMenuValueText.ChangeText("...")
			case 3:
				gamepad.CurrentButtonMapping = m.buttonMapping
				m.initialized = false
				setting.Opened = nil
			case 16:
				m.initialized = false
				setting.Opened = nil
			}
		}
	}
}
