package game_ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"strconv"
	"strings"
)

type viewComponent struct {
	components         []Component
	size               *image.Point
	style              ViewStyle
	extraStyles        []ViewStyle
	actionAreaMinPoint image.Point
	actionAreaMaxPoint image.Point
}
type View = *viewComponent
type ViewStyle struct {
	/* "#FFF" | "#FFFFFF" | "#FFFFFFFF" */
	BackgroundColor string
	/* "number" | "vertical horizontal" | "top horizontal bottom" | "top right bottom left" */
	Padding string
	/* "number" | "vertical horizontal" | "top horizontal bottom" | "top right bottom left" */
	Margin string
	/* "number" | "top_left top_right bottom_right bottom_left" */
	Radius string
	/* "number" | "vertical horizontal" | "top horizontal bottom" | "top right bottom left" */
	BorderWidth string
	/* "#FFF" | "#FFFFFF" | "#FFFFFFFF" | "top_left top_right bottom_right bottom_left" */
	BorderColor string
	/* "" | "number" */
	Width, Height      string
	Direction          DirectionType
	PositionHorizontal PositionType
	PositionVertical   PositionType
}

type DirectionType = string

const (
	DirectionUndefined DirectionType = ""
	Horizontal         DirectionType = "horizontal"
	Vertical           DirectionType = "vertical"
)

type PositionType = string

const (
	PositionUndefined PositionType = ""
	First             PositionType = "first"
	Center            PositionType = "center"
	Last              PositionType = "last"
)

func mergeViewStyle(target ViewStyle, styles []ViewStyle) ViewStyle {
	for i := range styles {
		if len(styles[i].BackgroundColor) > 0 {
			target.BackgroundColor = styles[i].BackgroundColor
		}
		if len(styles[i].Padding) > 0 {
			target.Padding = styles[i].Padding
		}
		if len(styles[i].Margin) > 0 {
			target.Margin = styles[i].Margin
		}
		if len(styles[i].Radius) > 0 {
			target.Radius = styles[i].Radius
		}
		if len(styles[i].BorderWidth) > 0 {
			target.BorderWidth = styles[i].BorderWidth
		}
		if len(styles[i].BorderColor) > 0 {
			target.BorderColor = styles[i].BorderColor
		}
		if len(styles[i].Width) > 0 {
			target.Width = styles[i].Width
		}
		if len(styles[i].Height) > 0 {
			target.Height = styles[i].Height
		}
		if len(styles[i].Direction) > 0 {
			target.Direction = styles[i].Direction
		}
		if len(styles[i].PositionHorizontal) > 0 {
			target.PositionHorizontal = styles[i].PositionHorizontal
		}
		if len(styles[i].PositionVertical) > 0 {
			target.PositionVertical = styles[i].PositionVertical
		}
	}
	return target
}

func getDefaultViewStyle() ViewStyle {
	return ViewStyle{
		BackgroundColor:    "#00000000",
		Padding:            "0",
		Margin:             "0",
		Radius:             "0",
		BorderWidth:        "0",
		BorderColor:        "#00000000",
		Width:              "0",
		Height:             "0",
		Direction:          Vertical,
		PositionHorizontal: First,
		PositionVertical:   First,
	}
}

func NewView(components []Component, styles ...ViewStyle) View {
	var style = mergeViewStyle(getDefaultViewStyle(), styles)
	return &viewComponent{components: components, style: style, extraStyles: []ViewStyle{}}
}

/* return (top, right, bottom, left) */
func strToArea(str string) (int, int, int, int) {
	var param = strings.Split(strings.TrimSpace(str), " ")
	if len(param) == 1 {
		var around, _ = strconv.Atoi(param[0])
		return around, around, around, around
	} else if len(param) == 2 {
		var vertical, _ = strconv.Atoi(param[0])
		var horizontal, _ = strconv.Atoi(param[1])
		return vertical, horizontal, vertical, horizontal
	} else if len(param) == 3 {
		var top, _ = strconv.Atoi(param[0])
		var horizontal, _ = strconv.Atoi(param[1])
		var bottom, _ = strconv.Atoi(param[2])
		return top, horizontal, bottom, horizontal
	} else if len(param) == 4 {
		var top, _ = strconv.Atoi(param[0])
		var right, _ = strconv.Atoi(param[1])
		var bottom, _ = strconv.Atoi(param[2])
		var left, _ = strconv.Atoi(param[3])
		return top, right, bottom, left
	}
	panic("illegal area param: " + str)
}

func strToColorCodes(str string) (string, string, string, string) {
	var strings = strings.Split(str, " ")
	if len(strings) == 4 {
		return strings[0], strings[1], strings[2], strings[3]
	} else {
		return strings[0], strings[0], strings[0], strings[0]
	}
}

func mixColorCode(r1, g1, b1, a1, r2, g2, b2, a2 uint32, rate float32) (uint32, uint32, uint32, uint32) {
	var r = uint32(int64(r2) - int64(float64(int64(r2)-int64(r1))*float64(rate)))
	var g = uint32(int64(g2) - int64(float64(int64(g2)-int64(g1))*float64(rate)))
	var b = uint32(int64(b2) - int64(float64(int64(b2)-int64(b1))*float64(rate)))
	var a = uint32(int64(a2) - int64(float64(int64(a2)-int64(a1))*float64(rate)))
	return r, g, b, a
}
func (v View) PushStyle(style ViewStyle) {
	v.extraStyles = append(v.extraStyles, style)
}
func (v View) PopStyle() {
	if len(v.extraStyles) > 0 {
		v.extraStyles = v.extraStyles[:len(v.extraStyles)-1]
	}
}
func (v View) ReplaceStyle(position int, style ViewStyle) {
	if position >= v.GetStylesCount() {
		v.extraStyles = append(v.extraStyles, style)
	} else {
		v.extraStyles[position] = style
	}
}
func (v View) GetStylesCount() int {
	return len(v.extraStyles)
}

func (v View) getContentSize() image.Point {
	var x, y = 0, 0
	var style = mergeViewStyle(v.style, v.extraStyles)
	for _, component := range v.components {
		var contentSize = component.GetSize()
		if style.Direction == Horizontal {
			if y <= contentSize.Y {
				y = contentSize.Y
			}
			x += contentSize.X
		} else {
			if x <= contentSize.X {
				x = contentSize.X
			}
			y += contentSize.Y
		}
	}
	for _, str := range []string{style.Padding, style.Margin, style.BorderWidth} {
		var top, right, bottom, left = strToArea(str)
		x += left + right
		y += top + bottom
	}
	return image.Point{X: x, Y: y}
}

func (v View) GetSize() image.Point {
	var point = v.getContentSize()
	var x = point.X
	var y = point.Y
	var style = mergeViewStyle(v.style, v.extraStyles)
	var width, _ = strconv.Atoi(style.Width)
	var height, _ = strconv.Atoi(style.Height)
	if x < width {
		x = width
	}
	if y < height {
		y = height
	}
	return image.Point{X: x, Y: y}
}

func (v View) GetActionArea() (image.Point, image.Point) {
	return v.actionAreaMinPoint, v.actionAreaMaxPoint
}

func (v View) Draw(screen *ebiten.Image, x, y int) {
	var style = mergeViewStyle(v.style, v.extraStyles)
	var marginTop, marginRight, marginBottom, marginLeft = strToArea(style.Margin)
	var borderTop, borderRight, borderBottom, borderLeft = strToArea(style.BorderWidth)
	var paddingTop, paddingRight, paddingBottom, paddingLeft = strToArea(style.Padding)
	var radiusTopLeft, radiusTopRight, radiusBottomRight, radiusBottomLeft = strToArea(style.Radius)

	var marginWidth = marginLeft + marginRight
	var borderWidth = borderLeft + borderRight
	var paddingWidth = paddingLeft + paddingRight
	var marginHeight = marginTop + marginBottom
	var borderHeight = borderTop + borderBottom
	var paddingHeight = paddingTop + paddingBottom

	var contentSize = v.getContentSize()
	var size = v.GetSize()

	var positionH float64 = 0
	var positionV float64 = 0

	var radiusMin = size.X - marginWidth - borderWidth
	if (size.Y - marginHeight - borderHeight) < radiusMin {
		radiusMin = size.Y - marginHeight - borderHeight
	}
	radiusMin /= 2
	if radiusMin < radiusTopLeft {
		radiusTopLeft = radiusMin
	}
	if radiusMin < radiusTopRight {
		radiusTopRight = radiusMin
	}
	if radiusMin < radiusBottomRight {
		radiusBottomRight = radiusMin
	}
	if radiusMin < radiusBottomLeft {
		radiusBottomLeft = radiusMin
	}

	if style.PositionHorizontal == Center {
		positionH = 0.5
	} else if style.PositionHorizontal == Last {
		positionH = 1
	}
	if style.PositionVertical == Center {
		positionV = 0.5
	} else if style.PositionVertical == Last {
		positionV = 1
	}

	var contentLeft = int(positionH * float64(size.X-contentSize.X))
	var contentTop = int(positionV * float64(size.Y-contentSize.Y))

	v.actionAreaMinPoint = image.Point{X: x + marginLeft, Y: y + marginTop}
	v.actionAreaMaxPoint = image.Point{X: v.actionAreaMinPoint.X + size.X - marginWidth, Y: v.actionAreaMinPoint.Y + size.Y - marginHeight}

	// draw border
	{
		var color1, color2, color3, color4 = strToColorCodes(style.BorderColor)
		var r1, g1, b1, a1 = colorCodeToColor(color1).RGBA()
		var r2, g2, b2, a2 = colorCodeToColor(color2).RGBA()
		var r3, g3, b3, a3 = colorCodeToColor(color3).RGBA()
		var r4, g4, b4, a4 = colorCodeToColor(color4).RGBA()
		if (borderWidth > 0 || borderHeight > 0) && (a1 > 0 || a2 > 0 || a3 > 0 || a4 > 0) {
			var path = vector.Path{}
			path.MoveTo(float32(radiusTopLeft), 0)
			path.LineTo(float32(size.X-marginWidth-radiusTopRight), 0)
			path.QuadTo(float32(size.X-marginWidth), 0, float32(size.X-marginWidth), float32(radiusTopRight))
			path.LineTo(float32(size.X-marginWidth), float32(size.Y-marginHeight-radiusBottomRight))
			path.QuadTo(float32(size.X-marginWidth), float32(size.Y-marginHeight), float32(size.X-marginWidth-radiusBottomRight), float32(size.Y-marginHeight))
			path.LineTo(float32(radiusBottomLeft), float32(size.Y-marginHeight))
			path.QuadTo(0, float32(size.Y-marginHeight), float32(0), float32(size.Y-marginHeight-radiusBottomLeft))
			path.LineTo(float32(0), float32(radiusTopLeft))
			path.QuadTo(0, 0, float32(radiusTopLeft), 0)

			path.MoveTo(float32(borderLeft+radiusTopLeft), float32(borderTop))
			path.LineTo(float32(size.X-marginWidth-borderRight-radiusTopRight), float32(borderTop))
			path.QuadTo(float32(size.X-marginWidth-borderRight), float32(borderTop), float32(size.X-marginWidth-borderRight), float32(borderTop+radiusTopRight))
			path.LineTo(float32(size.X-marginWidth-borderRight), float32(size.Y-marginHeight-borderBottom-radiusBottomRight))
			path.QuadTo(float32(size.X-marginWidth-borderRight), float32(size.Y-marginHeight-borderBottom), float32(size.X-marginWidth-borderRight-radiusBottomRight), float32(size.Y-marginHeight-borderBottom))
			path.LineTo(float32(borderLeft+radiusBottomLeft), float32(size.Y-marginHeight-borderBottom))
			path.QuadTo(float32(borderLeft), float32(size.Y-marginHeight-borderBottom), float32(borderLeft), float32(size.Y-marginHeight-borderBottom-radiusBottomLeft))
			path.LineTo(float32(borderLeft), float32(borderTop+radiusTopLeft))
			path.QuadTo(float32(borderLeft), float32(borderTop), float32(borderLeft+radiusTopLeft), float32(borderTop))
			path.Close()

			var vs, is = path.AppendVerticesAndIndicesForFilling(nil, nil)
			var minX, maxX, minY, maxY float32 = 0, 0, 0, 0
			for i := range vs {
				if minX < vs[i].DstX {
					minX = vs[i].DstX
				}
				if maxX > vs[i].DstX {
					maxX = vs[i].DstX
				}
				if minY < vs[i].DstY {
					minY = vs[i].DstY
				}
				if maxY > vs[i].DstY {
					maxY = vs[i].DstY
				}
			}
			for i := range vs {
				var rateX = (vs[i].DstX - minX) / (maxX - minX)
				var rateY = (vs[i].DstY - minY) / (maxY - minY)
				var rX1, gX1, bX1, aX1 = mixColorCode(r1, g1, b1, a1, r2, g2, b2, a2, rateX)
				var rX2, gX2, bX2, aX2 = mixColorCode(r4, g4, b4, a4, r3, g3, b3, a3, rateX)
				var r, g, b, a = mixColorCode(rX1, gX1, bX1, aX1, rX2, gX2, bX2, aX2, rateY)
				vs[i].DstX += float32(x + marginLeft)
				vs[i].DstY += float32(y + marginTop)
				vs[i].ColorR = float32(r) / 0xffff
				vs[i].ColorG = float32(g) / 0xffff
				vs[i].ColorB = float32(b) / 0xffff
				vs[i].ColorA = float32(a) / 0xffff
			}
			screen.DrawTriangles(vs, is, emptySubImage, &ebiten.DrawTrianglesOptions{
				FillRule: ebiten.EvenOdd,
			})
		}
	}

	// draw base
	{
		var color1, color2, color3, color4 = strToColorCodes(style.BackgroundColor)
		var r1, g1, b1, a1 = colorCodeToColor(color1).RGBA()
		var r2, g2, b2, a2 = colorCodeToColor(color2).RGBA()
		var r3, g3, b3, a3 = colorCodeToColor(color3).RGBA()
		var r4, g4, b4, a4 = colorCodeToColor(color4).RGBA()
		if a1 > 0 || a2 > 0 || a3 > 0 || a4 > 0 {
			var path = vector.Path{}
			path.MoveTo(float32(borderLeft+radiusTopLeft), float32(borderTop))
			path.LineTo(float32(size.X-marginWidth-borderRight-radiusTopRight), float32(borderTop))
			path.QuadTo(float32(size.X-marginWidth-borderRight), float32(borderTop), float32(size.X-marginWidth-borderRight), float32(borderTop+radiusTopRight))
			path.LineTo(float32(size.X-marginWidth-borderRight), float32(size.Y-marginHeight-borderBottom-radiusBottomRight))
			path.QuadTo(float32(size.X-marginWidth-borderRight), float32(size.Y-marginHeight-borderBottom), float32(size.X-marginWidth-borderRight-radiusBottomRight), float32(size.Y-marginHeight-borderBottom))
			path.LineTo(float32(borderLeft+radiusBottomLeft), float32(size.Y-marginHeight-borderBottom))
			path.QuadTo(float32(borderLeft), float32(size.Y-marginHeight-borderBottom), float32(borderLeft), float32(size.Y-marginHeight-borderBottom-radiusBottomLeft))
			path.LineTo(float32(borderLeft), float32(borderTop+radiusTopLeft))
			path.QuadTo(float32(borderLeft), float32(borderTop), float32(borderLeft+radiusTopLeft), float32(borderTop))
			path.Close()

			var vs, is = path.AppendVerticesAndIndicesForFilling(nil, nil)
			var maxX, minX, maxY, minY float32 = 0, 0, 0, 0
			for i := range vs {
				if maxX > vs[i].DstX {
					maxX = vs[i].DstX
				}
				if minX < vs[i].DstX {
					minX = vs[i].DstX
				}
				if maxY > vs[i].DstY {
					maxY = vs[i].DstY
				}
				if minY < vs[i].DstY {
					minY = vs[i].DstY
				}
			}
			for i := range vs {
				var rateX = (vs[i].DstX - minX) / (maxX - minX)
				var rateY = (vs[i].DstY - minY) / (maxY - minY)
				var rX1, gX1, bX1, aX1 = mixColorCode(r1, g1, b1, a1, r2, g2, b2, a2, rateX)
				var rX2, gX2, bX2, aX2 = mixColorCode(r4, g4, b4, a4, r3, g3, b3, a3, rateX)
				var r, g, b, a = mixColorCode(rX1, gX1, bX1, aX1, rX2, gX2, bX2, aX2, rateY)
				vs[i].DstX += float32(x + marginLeft)
				vs[i].DstY += float32(y + marginTop)
				vs[i].ColorR = float32(r) / 0xffff
				vs[i].ColorG = float32(g) / 0xffff
				vs[i].ColorB = float32(b) / 0xffff
				vs[i].ColorA = float32(a) / 0xffff
			}
			screen.DrawTriangles(vs, is, emptySubImage, &ebiten.DrawTrianglesOptions{
				FillRule: ebiten.EvenOdd,
			})
		}
	}

	var _x = marginLeft + borderLeft + paddingLeft + contentLeft
	var _y = marginTop + borderTop + paddingTop + contentTop
	for _, component := range v.components {
		var componentSize = component.GetSize()
		var __x = int(positionH * float64((contentSize.X-marginWidth-borderWidth-paddingWidth)-componentSize.X))
		var __y = int(positionV * float64((contentSize.Y-marginHeight-borderHeight-paddingHeight)-componentSize.Y))
		if style.Direction == Horizontal {
			__x = 0
		} else {
			__y = 0
		}
		component.Draw(screen, x+_x+__x, y+_y+__y)
		if style.Direction == Horizontal {
			_x += componentSize.X
		} else {
			_y += componentSize.Y
		}
	}
}
