# q2d

This software was created with Gemini CLI and Antigravity using the Gemini 2.5
and 3 models. It is intended for use by other LLMs.

`q2d` is a lightweight, pure Go 2D graphics library designed for software
rendering. It provides a stateful drawing context with support for clipping,
sub-images, basic shapes, text rendering, and image composition. It is ideal for
building custom UI frameworks, simple 2D games, or applications requiring direct
pixel manipulation without heavy external dependencies.

## Installation

```bash
go get github.com/qbradq/q2d
```

## Features

*   **Core Primitives**: `Point`, `Rectangle`, and `Color` (RGBA & HSL support).
*   **Drawing Surface**: `Image` struct backed by a raw byte slice for direct
    access.
*   **State Management**: Stack-based clipping and sub-image support for
    hierarchical rendering.
*   **Text Rendering**: Built-in support for text drawing with word wrapping.
*   **Image Composition**: Scaling and alpha blending.
*   **Zero CGO**: Pure Go implementation.

## Usage Examples

### 1. Basic Setup and Pixel Manipulation

Create a new image and manipulate pixels directly.

```go
package main

import (
	"fmt"
	"github.com/qbradq/q2d"
)

func main() {
	// Create a 100x100 image
	img := q2d.NewImage(100, 100)

	// Define a color (Red)
	red := q2d.Color{255, 0, 0, 255}

	// Set a pixel at (10, 10)
	img.Set(q2d.Point{10, 10}, red)

	// Retrieve a pixel
	c := img.At(q2d.Point{10, 10})
	fmt.Printf("Color at (10, 10): %v\n", c)
}
```

**Result**: Creates an image in memory and sets a single red pixel.

### 2. Drawing Shapes

Draw filled rectangles, lines, and borders.

```go
package main

import "github.com/qbradq/q2d"

func main() {
	img := q2d.NewImage(200, 200)
	blue := q2d.Color{0, 0, 255, 255}
	green := q2d.Color{0, 255, 0, 255}

	// Fill the entire current clip/image with blue
	img.Fill(blue)

	// Draw a horizontal line
	img.HLine(50, 10, 190, 2, green) // y, x1, x2, height, color

	// Draw a vertical line
	img.VLine(50, 10, 190, 2, green) // x, y1, y2, width, color

	// Draw a border around the current bounds
	img.Border(q2d.Color{255, 255, 255, 255})
}
```

**Result**: An image filled with blue, containing a green cross and a white
border.

### 3. Clipping and Sub-Images

Manage drawing regions using a stack. This is useful for UI widgets where
children should not draw outside their parents.

```go
package main

import "github.com/qbradq/q2d"

func main() {
	img := q2d.NewImage(200, 200)
	red := q2d.Color{255, 0, 0, 255}

	// Create a sub-image (local coordinate system)
	// Origin (0,0) in the sub-image is (50, 50) in the main image
	img.PushSubImage(q2d.Rectangle{50, 50, 100, 100})

	// Restrict drawing to a smaller region within the sub-image
	img.PushClip(q2d.Rectangle{10, 10, 80, 80})

	// This point is relative to the sub-image origin
	// It will be drawn at absolute (60, 60)
	img.Set(q2d.Point{10, 10}, red)

	// This point is outside the clipping region, so it won't be drawn
	img.Set(q2d.Point{0, 0}, red)

	// Restore state
	img.PopClip()
	img.PopSubImage()
}
```

**Result**: Only pixels within the defined clipping region are modified.
Coordinate systems are automatically translated.

### 4. Color Manipulation (HSL)

Work with colors in the HSL space for easier theming (lightening, darkening, hue
shifts).

```go
package main

import (
	"fmt"
	"github.com/qbradq/q2d"
)

func main() {
	red := q2d.Color{255, 0, 0, 255}

	// Darken by 20%
	darkRed := red.Darken(0.2)

	// Lighten by 20%
	lightRed := red.Lighten(0.2)

	// Shift Hue by 180 degrees (Complementary color)
	cyan := red.AdjustHue(180)

	fmt.Printf("Original: %v, Darker: %v, Cyan: %v\n", red, darkRed, cyan)
}
```

**Result**: Generates new color variations based on the input color.

### 5. Text Rendering

Draw text using standard Go font faces.

```go
package main

import (
	"github.com/qbradq/q2d"
	"golang.org/x/image/font/basicfont"
)

func main() {
	img := q2d.NewImage(200, 100)
	white := q2d.Color{255, 255, 255, 255}

	// Use a basic font face (or load a TTF/OTF)
	face := basicfont.Face7x13

	// Draw text at (10, 20) with wrapping enabled
	img.Text(q2d.Point{10, 20}, white, face, true, "Hello, World!")
}
```

**Result**: Renders "Hello, World!" onto the image.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
