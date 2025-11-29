# q2d

## Project Overview

`q2d` is a lightweight, pure Go 2D graphics library designed for software rendering. It provides a stateful drawing context with support for clipping, sub-images, basic shapes, and text rendering. It is particularly suitable for building custom UI frameworks or simple 2D games where hardware acceleration is not a primary requirement or where direct pixel manipulation is desired.

## Key Features

*   **Core Primitives:** `Point` (x, y), `Rectangle` (x, y, w, h), and `Color` (RGBA).
*   **Color Management:** Built-in support for RGBA colors and HSL conversions (Hue, Saturation, Lightness).
*   **Drawing Surface:** The `Image` struct is the central drawing area, backed by a raw byte slice (`Pix`).
*   **State Management:** Supports a stack-based approach for clipping regions and sub-images (`PushSubImage`, `PopSubImage`, `PushClip`, `PopClip`), allowing for hierarchical rendering (e.g., nested UI widgets).
*   **Text Rendering:** Integrated text drawing using `golang.org/x/image/font`, with support for word wrapping.
*   **Image Composition:** specific optimizations for scaling and alpha blending images.

## Directory Structure

*   `data.go`: (Assumed) Likely contains embedded data or resource loading logic.
*   `image.go`: Core implementation of the `Image` struct and drawing primitives (lines, borders, text, blitting).
*   `types.go`: Definitions of basic types: `Point`, `Rectangle`, `Color` and their associated methods.
*   `q2d_test.go`: Unit tests demonstrating usage and verifying logic (especially clipping).
*   `fonts/`: Directory containing font resources (e.g., `unscii-16.otf`).

## Building and Running

Since this is a standard Go library, use the standard Go toolchain.

### Build

```bash
go build ./...
```

### Test

The project includes unit tests, particularly for checking the clipping and coordinate translation logic.

```bash
go test ./...
```

## Development Conventions

*   **Coordinate System:** (0, 0) is at the top-left corner.
*   **Pixel Access:** Direct manipulation of the `Pix` slice (RGBA stride) is common for performance-critical operations.
*   **State Stack:** When implementing new drawing functions, respect the current clipping region defined in the `stateStack`. Always use `img.currentState()` to access the active origin and clip bounds.
*   **Dependencies:** Relies on `golang.org/x/image` for advanced image and font handling.
