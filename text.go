package game_ui

import (
	"github.com/hajimehoshi/bitmapfont/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"image"
	"strings"
)

type textComponent struct {
	str   string
	size  *image.Point
	style TextStyle
}
type Text = *textComponent
type TextStyle struct {
	/* #FFF | #FFFFFF | #FFFFFFFF */
	Color      string
	LineHeight int
	Font       *TextFont
	Width      int
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
		if len(styles[i].Color) > 0 {
			target.Color = styles[i].Color
		}
		if styles[i].LineHeight > 0 {
			target.LineHeight = styles[i].LineHeight
		}
		if styles[i].Font != nil {
			target.Font = styles[i].Font
		}
		if styles[i].Width > 0 {
			target.Width = styles[i].Width
		}
	}
	return target
}

func getDefaultTextStyle() TextStyle {
	return TextStyle{
		Color:      "#ffffffff",
		LineHeight: 12,
		Font:       &defaultTextFont,
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
	if t.style.Width > 0 {
		var str = ""
		var line = ""
		for _, _char := range t.str {
			var char = string(_char)
			if text.BoundString(t.style.Font.face, line+char).Size().X > t.style.Width {
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

	var size = text.BoundString(text.FaceWithLineHeight(t.style.Font.face, float64(t.style.LineHeight)), t.str).Size()
	var lineCount = strings.Count(t.str, "\n") + 1
	t.size = &image.Point{
		X: size.X - 1,
		Y: lineCount * t.style.LineHeight,
	}
	return *t.size
}

func (t Text) Draw(screen *ebiten.Image, x, y int) {
	// m := bitmapfont.Face.Metrics()
	// bounds := text.BoundString(bitmapfont.Face, "a")
	// x += bounds.Dx() - 2
	// y += int(float64(bounds.Dy()) * float64(m.Ascent) / float64(m.Height))
	t.GetSize()
	var font = t.style.Font
	x += font.xAdjustment
	y += t.style.LineHeight - t.style.LineHeight/2 + font.yAdjustment
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(colorCodeToColor(t.style.Color))
	text.DrawWithOptions(screen, t.str, text.FaceWithLineHeight(font.face, float64(t.style.LineHeight)), op)
}
