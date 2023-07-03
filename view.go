package game_ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
)

type viewComponent struct {
	components  []Component
	size        *image.Point
	style       ViewStyle
	extraStyles []ViewStyle
	drawnArea   image.Rectangle
	screenSize  image.Point
}
type View = *viewComponent
type ViewStyle struct {
	/* top_left top_right bottom_right bottom_left */
	BackgroundColor, BorderColor *[4]color.Color
	/* top right bottom left */
	Padding, Margin, BorderWidth *[4]sizeSeg
	/* top_left top_right bottom_right bottom_left */
	Radius             *[4]int
	Width, Height      *sizeSeg
	Direction          *DirectionType
	PositionHorizontal *PositionType
	PositionVertical   *PositionType
	IsFloating         bool
}

type DirectionType = string

const (
	Horizontal DirectionType = "horizontal"
	Vertical   DirectionType = "vertical"
)

type PositionType = string

const (
	First  PositionType = "first"
	Center PositionType = "center"
	Last   PositionType = "last"
)

func mergeViewStyle(target ViewStyle, styles []ViewStyle) ViewStyle {
	for i := range styles {
		if styles[i].BackgroundColor != nil {
			target.BackgroundColor = styles[i].BackgroundColor
		}
		if styles[i].Padding != nil {
			target.Padding = styles[i].Padding
		}
		if styles[i].Margin != nil {
			target.Margin = styles[i].Margin
		}
		if styles[i].Radius != nil {
			target.Radius = styles[i].Radius
		}
		if styles[i].BorderWidth != nil {
			target.BorderWidth = styles[i].BorderWidth
		}
		if styles[i].BorderColor != nil {
			target.BorderColor = styles[i].BorderColor
		}
		if styles[i].Width != nil {
			target.Width = styles[i].Width
		}
		if styles[i].Height != nil {
			target.Height = styles[i].Height
		}
		if styles[i].Direction != nil {
			target.Direction = styles[i].Direction
		}
		if styles[i].PositionHorizontal != nil {
			target.PositionHorizontal = styles[i].PositionHorizontal
		}
		if styles[i].PositionVertical != nil {
			target.PositionVertical = styles[i].PositionVertical
		}
		target.IsFloating = target.IsFloating || styles[i].IsFloating
	}
	return target
}

func NewView(components []Component, styles ...ViewStyle) View {
	var style = mergeViewStyle(ViewStyle{}, styles)
	return &viewComponent{components: components, style: style, extraStyles: []ViewStyle{}}
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

func getSizePx(screenSize image.Point, size [4]sizeSeg) (int, int, int, int) {
	return calcSize(screenSize, size[0]), calcSize(screenSize, size[1]), calcSize(screenSize, size[2]), calcSize(screenSize, size[3])
}

func (v View) getContentSize() image.Point {
	var x, y = 0, 0
	var style = mergeViewStyle(v.style, v.extraStyles)
	for _, component := range v.components {
		if component.IsFloating() {
			continue
		}
		var contentSize = component.GetSize()
		if style.Direction == nil || *style.Direction == Vertical {
			if x <= contentSize.X {
				x = contentSize.X
			}
			y += contentSize.Y
		} else {
			if y <= contentSize.Y {
				y = contentSize.Y
			}
			x += contentSize.X
		}
	}
	for _, size := range []*[4]sizeSeg{style.Padding, style.Margin, style.BorderWidth} {
		if size == nil {
			continue
		}
		var top, right, bottom, left = getSizePx(v.screenSize, *size)
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
	var width = 0
	if style.Width != nil {
		width = calcSize(v.screenSize, *style.Width)
	}
	var height = 0
	if style.Height != nil {
		height = calcSize(v.screenSize, *style.Height)
	}
	if x < width {
		x = width
	}
	if y < height {
		y = height
	}
	return image.Point{X: x, Y: y}
}

func (v View) Draw(screen *ebiten.Image, x, y int) {
	var style = mergeViewStyle(v.style, v.extraStyles)
	var marginTop, marginRight, marginBottom, marginLeft int
	var borderTop, borderRight, borderBottom, borderLeft int
	var paddingTop, paddingRight, paddingBottom, paddingLeft int

	v.screenSize = screen.Bounds().Size()

	if style.Margin != nil {
		marginTop, marginRight, marginBottom, marginLeft = getSizePx(v.screenSize, *style.Margin)
	}
	if style.BorderWidth != nil {
		borderTop, borderRight, borderBottom, borderLeft = getSizePx(v.screenSize, *style.BorderWidth)
	}
	if style.Padding != nil {
		paddingTop, paddingRight, paddingBottom, paddingLeft = getSizePx(v.screenSize, *style.Padding)
	}
	var radiusTopLeft, radiusTopRight, radiusBottomRight, radiusBottomLeft int
	if style.Radius != nil {
		radiusTopLeft, radiusTopRight, radiusBottomRight, radiusBottomLeft = style.Radius[0], style.Radius[1], style.Radius[2], style.Radius[3]
	}

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

	if style.PositionHorizontal != nil {
		if *style.PositionHorizontal == Center {
			positionH = 0.5
		} else if *style.PositionHorizontal == Last {
			positionH = 1
		}
	}
	if style.PositionVertical != nil {
		if *style.PositionVertical == Center {
			positionV = 0.5
		} else if *style.PositionVertical == Last {
			positionV = 1
		}
	}

	var contentLeft = int(positionH * float64(size.X-contentSize.X))
	var contentTop = int(positionV * float64(size.Y-contentSize.Y))

	var minX, minY = x + marginLeft, y + marginTop
	v.drawnArea = image.Rect(minX, minY, minX+size.X-marginWidth, minY+size.Y-marginHeight)

	// draw border
	if style.BorderColor != nil {
		var r1, g1, b1, a1 = style.BorderColor[0].RGBA()
		var r2, g2, b2, a2 = style.BorderColor[1].RGBA()
		var r3, g3, b3, a3 = style.BorderColor[2].RGBA()
		var r4, g4, b4, a4 = style.BorderColor[3].RGBA()
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
	if style.BackgroundColor != nil {
		var r1, g1, b1, a1 = style.BackgroundColor[0].RGBA()
		var r2, g2, b2, a2 = style.BackgroundColor[1].RGBA()
		var r3, g3, b3, a3 = style.BackgroundColor[2].RGBA()
		var r4, g4, b4, a4 = style.BackgroundColor[3].RGBA()
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
		if style.Direction == nil || *style.Direction == Vertical {
			__y = 0
		} else {
			__x = 0
		}
		component.Draw(screen, x+_x+__x, y+_y+__y)
		if component.IsFloating() {
			continue
		}
		if style.Direction == nil || *style.Direction == Vertical {
			_y += componentSize.Y
		} else {
			_x += componentSize.X
		}
	}
}

func (v View) IsFloating() bool {
	var style = mergeViewStyle(v.style, v.extraStyles)
	return style.IsFloating
}

func (v View) Components() []Component {
	return v.components
}

func (v View) Area() image.Rectangle {
	return v.drawnArea
}
