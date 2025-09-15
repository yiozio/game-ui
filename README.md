# Game UI

A flexible, component-based UI framework for [Ebiten](https://github.com/hajimehoshi/ebiten) game development in Go. This library provides a CSS-like styling system with support for responsive layouts, gradients, borders, and modern UI patterns.

## Features

- **Component-based Architecture**: Build complex UIs by composing simple components
- **CSS-like Styling**: Familiar styling concepts including margins, padding, borders, and box model
- **Responsive Design**: Support for viewport-relative units (vw, vh) and pixel values
- **Gradient Support**: Multi-corner gradients for backgrounds and borders
- **Text Rendering**: Flexible text components with custom fonts and line wrapping
- **Layout System**: Horizontal and vertical layout directions with alignment options
- **Floating Components**: Position components outside normal document flow

## Quick Start

```go
package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    gameui "github.com/yiozio/game-ui"
)

func main() {
    // Create a simple UI with a styled view containing text
    text := gameui.NewText("Hello, World!", gameui.TextStyle{
        Color: gameui.Color(0xFFFFFFFF), // White text
    })

    view := gameui.NewView([]gameui.Component{text}, gameui.ViewStyle{
        BackgroundColor: gameui.ColorCode1(0x333333FF), // Dark gray background
        Padding: gameui.Size1(gameui.Px(20)),           // 20px padding on all sides
        Margin: gameui.Size1(gameui.Px(10)),            // 10px margin on all sides
        Radius: gameui.Radius1(8),                      // 8px border radius
    })

    window := gameui.NewWindow([]gameui.Component{view})

    // Use window in your Ebiten game...
}
```

## Core Components

### View
The primary container component that supports all styling features:

```go
view := gameui.NewView(components, gameui.ViewStyle{
    // Background with gradient
    BackgroundColor: gameui.ColorCodeGradation(
        0xFF0000FF, // Top-left: Red
        0x00FF00FF, // Top-right: Green
        0x0000FFFF, // Bottom-right: Blue
        0xFFFF00FF, // Bottom-left: Yellow
    ),

    // Padding and margins
    Padding: gameui.Size4(gameui.Px(10), gameui.Px(20), gameui.Px(10), gameui.Px(20)),
    Margin: gameui.Size2(gameui.Px(5), gameui.Px(10)),

    // Borders
    BorderWidth: gameui.Size1(gameui.Px(2)),
    BorderColor: gameui.ColorCode1(0x000000FF),

    // Border radius
    Radius: gameui.Radius4(10, 5, 10, 5),

    // Layout
    Direction: &gameui.Horizontal,
    PositionHorizontal: &gameui.Center,
    PositionVertical: &gameui.Center,
})
```

### Text
Text components with customizable styling:

```go
text := gameui.NewText("Sample Text", gameui.TextStyle{
    Color: gameui.Color(0xFFFFFFFF),
    LineHeight: gameui.Px(16),
    Width: gameui.Vw(0.8), // 80% of viewport width for text wrapping
})
```

### Window
Root container for organizing multiple components:

```go
window := gameui.NewWindow([]gameui.Component{
    header,
    content,
    footer,
})
```

## Sizing System

The library supports three types of size units:

- **Pixels**: `gameui.Px(100)` - Fixed pixel values
- **Viewport Width**: `gameui.Vw(0.5)` - 50% of screen width
- **Viewport Height**: `gameui.Vh(0.3)` - 30% of screen height

Size values can be combined:
```go
// 50% viewport width plus 20 pixels
size := gameui.Vw(0.5, gameui.Px(20))
```

### Size Utilities

Convenient functions for common size patterns:

```go
// Same value for all sides (top, right, bottom, left)
gameui.Size1(gameui.Px(10))

// Vertical and horizontal values
gameui.Size2(gameui.Px(5), gameui.Px(10)) // 5px top/bottom, 10px left/right

// Individual values for top, sides, bottom
gameui.Size3(gameui.Px(5), gameui.Px(10), gameui.Px(15))

// Individual values for each side
gameui.Size4(gameui.Px(5), gameui.Px(10), gameui.Px(15), gameui.Px(20))
```

## Color System

Colors are specified as 32-bit hex values in RRGGBBAA format:

```go
// Solid colors
red := gameui.Color(0xFF0000FF)
blue := gameui.Color(0x0000FFFF)
transparent := gameui.Color(0xFF000000) // Red with 0 alpha

// Color arrays for gradients
solidRed := gameui.ColorCode1(0xFF0000FF)
horizontalGradient := gameui.ColorCodeHorizontal(0xFF0000FF, 0x0000FFFF)
verticalGradient := gameui.ColorCodeVertical(0xFF0000FF, 0x0000FFFF)
fourCornerGradient := gameui.ColorCodeGradation(0xFF0000FF, 0x00FF00FF, 0x0000FFFF, 0xFFFF00FF)
```

## Layout System

### Direction
Control how child components are arranged:

```go
// Stack components vertically (default)
view := gameui.NewView(components, gameui.ViewStyle{
    Direction: &gameui.Vertical,
})

// Arrange components horizontally
view := gameui.NewView(components, gameui.ViewStyle{
    Direction: &gameui.Horizontal,
})
```

### Positioning
Control component alignment within their container:

```go
view := gameui.NewView(components, gameui.ViewStyle{
    PositionHorizontal: &gameui.Center, // first, center, last
    PositionVertical: &gameui.Center,   // first, center, last
})
```

### Floating Components
Components can be positioned outside the normal layout flow:

```go
floatingView := gameui.NewView(components, gameui.ViewStyle{
    IsFloating: true,
})
```

## Dynamic Styling

Views support dynamic style changes with a stack-based system:

```go
view := gameui.NewView(components)

// Add temporary styling
view.PushStyle(gameui.ViewStyle{
    BackgroundColor: gameui.ColorCode1(0xFF0000FF),
})

// Remove the temporary style
view.PopStyle()

// Replace or add style at specific position
view.ReplaceStyle(0, gameui.ViewStyle{
    BackgroundColor: gameui.ColorCode1(0x00FF00FF),
})
```

## Custom Fonts

Create custom text fonts with position adjustments:

```go
customFont := gameui.NewTextFont(yourFontFace, 2, 4) // x adjustment, y adjustment

text := gameui.NewText("Custom Font", gameui.TextStyle{
    Font: &customFont,
})
```

## Integration with Ebiten

Implement the `Component` interface in your Ebiten game:

```go
type Game struct {
    ui gameui.Window
}

func (g *Game) Update() error {
    // Update game state
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    // Draw your UI
    g.ui.Draw(screen, 0, 0)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return 640, 480
}
```

## Examples

Check out the `example/` directory for complete working examples including:
- Menu systems
- Settings panels
- Game scenes
- Control handling

### Running the Example

You can run the example application in several ways:

**Native (Desktop):**
```bash
go run ./example/main.go
```

**Web Browser:**
```bash
go run github.com/hajimehoshi/wasmserve@latest ./example/main.go
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
