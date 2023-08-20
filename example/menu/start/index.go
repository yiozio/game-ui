package start

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yiozio/game-ui"
	"github.com/yiozio/game-ui/example/control"
	"github.com/yiozio/game-ui/example/control/gamepad"
	"github.com/yiozio/game-ui/example/def/common"
	actionEffect "github.com/yiozio/game-ui/example/effect/action"
	"github.com/yiozio/game-ui/example/menu"
	"github.com/yiozio/game-ui/example/menu/setting"
	controlMenu "github.com/yiozio/game-ui/example/menu/setting/control"
	"image"
	"strconv"
)

type Menu struct {
	game_ui.Window
	selectedMenuIndex int
	initialized       bool
}

func NewStartMenu() *Menu {
	var pos = game_ui.Center
	return &Menu{game_ui.NewWindow([]game_ui.Component{game_ui.NewView([]game_ui.Component{
		titleView,
		startView,
		settingView,
		exitView,
	}, game_ui.ViewStyle{
		Width:            game_ui.Vw(1),
		Height:           game_ui.Vh(1),
		PositionVertical: &pos,
	})}), -1, false}
}

func (m *Menu) Update(now int64, screenSize image.Point, mode control.Mode, enable bool) {
	if setting.Opened != nil {
		setting.Opened.Update(now, screenSize, mode, enable)
		return
	}

	if mode == control.Gamepad && m.selectedMenuIndex < 0 {
		m.selectedMenuIndex = 0
	}

	for _, v := range startMenuItems {
		v.ReplaceStyle(0, game_ui.ViewStyle{})
	}

	if mode != control.Gamepad {
		m.selectedMenuIndex = -1
	}
	var action = false
	if mode == control.Mouse {
		m.selectedMenuIndex, _ = menu.FindHoveredCell(startMenuItems)
		action = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	} else if mode == control.Touch {
	} else if mode == control.Gamepad {
		if m.selectedMenuIndex < 0 {
			m.selectedMenuIndex = 0
		}

		if gid := control.GetGamepadId(); gid != nil && m.initialized {
			if inpututil.IsStandardGamepadButtonJustPressed(*gid, gamepad.CurrentButtonMapping.Up) {
				startMenuItems[m.selectedMenuIndex].ReplaceStyle(0, game_ui.ViewStyle{})
				m.selectedMenuIndex += len(startMenuItems) - 1
			} else if inpututil.IsStandardGamepadButtonJustPressed(*gid, gamepad.CurrentButtonMapping.Down) {
				startMenuItems[m.selectedMenuIndex].ReplaceStyle(0, game_ui.ViewStyle{})
				m.selectedMenuIndex += 1
			}
			m.selectedMenuIndex = m.selectedMenuIndex % len(startMenuItems)

			action = inpututil.IsStandardGamepadButtonJustPressed(*gid, gamepad.CurrentButtonMapping.Action)
		}
	}

	if m.selectedMenuIndex >= 0 {
		startMenuItems[m.selectedMenuIndex].ReplaceStyle(0, game_ui.ViewStyle{BorderColor: game_ui.ColorCode1(0xffffffff)})
		if action {
			if mode == control.Mouse {
				actionEffect.StartEffect(now)
			}
			switch m.selectedMenuIndex {
			case 0:
			case 1:
				setting.Opened = controlMenu.NewSettingMenu()
			case 2:
			}
		}
	}
	m.initialized = true
}

func (m *Menu) Draw(screen *ebiten.Image, now int64, screenSize image.Point, mode control.Mode) {
	// draw background
	{
		var path = &vector.Path{}
		path.MoveTo(0, 0)
		path.LineTo(640, 0)
		path.LineTo(640, 480)
		path.LineTo(0, 480)
		path.Close()

		var vs, is = path.AppendVerticesAndIndicesForFilling(nil, nil)
		for i := range vs {
			vs[i].ColorR = 0x1 / float32(0xf)
			vs[i].ColorG = 0x1 / float32(0xf)
			vs[i].ColorB = 0x1 / float32(0xf)
			vs[i].ColorA = 1
		}
		screen.DrawTriangles(vs, is, common.EmptySubImage, &ebiten.DrawTrianglesOptions{
			FillRule: ebiten.EvenOdd,
		})
	}

	var n = 0xe * uint32((now/2)%16)
	if n > (0xe * 8) {
		n = 0xd*16 - n
	}
	n -= 1
	var bColor1 uint32 = 0xffffff90 + n
	var bColor2 uint32 = 0xffffff00
	var bgColor1 uint32 = 0x5599cc50 + n
	var bgColor2 uint32 = 0x5599cc00

	// draw window
	for i := range startMenuItems {
		if i == m.selectedMenuIndex && setting.Opened == nil {
			startMenuItems[i].ReplaceStyle(0, game_ui.ViewStyle{
				BorderColor:     game_ui.ColorCodeHorizontal(bColor1, bColor2),
				BackgroundColor: game_ui.ColorCodeHorizontal(bgColor1, bgColor2),
			})
			ebitenutil.DebugPrintAt(screen, "on", 300, i*15)
		} else {
			startMenuItems[i].PopStyle()
			ebitenutil.DebugPrintAt(screen, "off", 300, i*15)
		}
	}

	m.Window.Draw(screen, 0, 0)
	for i, view := range startMenuItems {
		var actionArea = view.Area()
		ebitenutil.DebugPrintAt(screen,
			strconv.Itoa(actionArea.Min.X)+","+
				strconv.Itoa(actionArea.Min.Y)+":"+
				strconv.Itoa(actionArea.Max.X)+","+
				strconv.Itoa(actionArea.Max.Y), 100, i*15)
	}

	if setting.Opened != nil {
		setting.Opened.Draw(screen, now, screenSize, mode)
	}
}
