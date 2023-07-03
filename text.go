package game_ui

import (
	"github.com/hajimehoshi/bitmapfont/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"strings"
)

type textComponent struct {
	str        string
	size       *image.Point
	style      TextStyle
	screenSize image.Point
	drawnArea  image.Rectangle
}
type Text = *textComponent
type TextStyle struct {
	Color      *color.Color
	LineHeight *sizeSeg
	Font       *TextFont
	Width      *sizeSeg
}
type TextFont struct {
	face        font.Face
	xAdjustment int
	yAdjustment int
}

var defaultTextFont = TextFont{
	face:        bitmapfont.Face,
	xAdjustment: 4,
	yAdjustment: 4,
}

func NewTextFont(face font.Face, xAdjustment int, yAdjustment int) TextFont {
	return TextFont{
		face,
		xAdjustment,
		yAdjustment,
	}
}

func mergeTextStyle(target TextStyle, styles []TextStyle) TextStyle {
	for i := range styles {
		if styles[i].Color != nil {
			target.Color = styles[i].Color
		}
		if styles[i].LineHeight != nil {
			target.LineHeight = styles[i].LineHeight
		}
		if styles[i].Font != nil {
			target.Font = styles[i].Font
		}
		if styles[i].Width != nil {
			target.Width = styles[i].Width
		}
	}
	return target
}

func getDefaultTextStyle() TextStyle {
	var c color.Color = color.White
	return TextStyle{
		Color:      &c,
		Font:       &defaultTextFont,
		LineHeight: Px(12),
	}
}

func NewText(str string, styles ...TextStyle) Text {
	var style = mergeTextStyle(getDefaultTextStyle(), styles)
	return &textComponent{str: str, size: nil, style: style}
}

func (t Text) GetSize() image.Point {
	if t.size != nil {
		return *t.size
	}
	if t.style.Width != nil {
		var str = ""
		var line = ""
		for _, _char := range t.str {
			var char = string(_char)
			if text.BoundString(t.style.Font.face, line+char).Size().X > calcSize(t.screenSize, *t.style.Width) {
				if len(line) == 0 {
					str += line + char + "\n"
					line = ""
					continue
				} else {
					str += line + "\n"
					line = char
					continue
				}
			}
			line += char
		}
		if len(line) > 0 {
			str += line
		}
		t.str = str
	}

	var lineHeightPx = calcSize(t.screenSize, *t.style.LineHeight)
	var size = text.BoundString(text.FaceWithLineHeight(t.style.Font.face, float64(lineHeightPx)), t.str).Size()
	var lineCount = strings.Count(t.str, "\n") + 1
	t.size = &image.Point{
		X: size.X - 1,
		Y: lineCount * lineHeightPx,
	}
	return *t.size
}

func (t Text) Draw(screen *ebiten.Image, x, y int) {
	t.screenSize = screen.Bounds().Size()
	var size = t.GetSize()
	t.drawnArea = image.Rect(x, y, x+size.X, y+size.Y)
	var font = t.style.Font
	var lineHeight = t.style.LineHeight
	var lineHeightPx = calcSize(t.screenSize, *lineHeight)
	x += font.xAdjustment
	y += lineHeightPx - lineHeightPx/2 + font.yAdjustment
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(*t.style.Color)
	text.DrawWithOptions(screen, t.str, text.FaceWithLineHeight(font.face, float64(lineHeightPx)), op)
}

func (t Text) ChangeText(text string) {
	t.str = text
	t.size = nil
}

func (t Text) IsFloating() bool {
	return false
}

func (t Text) Components() []Component {
	return []Component{}
}

func (t Text) Area() image.Rectangle {
	return t.drawnArea
}
