package q2d

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type subImageState struct {
	origin        Point
	bounds        Rectangle
	clip          Rectangle
	clippingStack []Rectangle
}

type Image struct {
	Pix    []uint8
	Stride int
	Rect   Rectangle

	stateStack []subImageState
}

func NewImage(w, h int) *Image {
	return &Image{
		Pix:    make([]uint8, 4*w*h),
		Stride: 4 * w,
		Rect:   Rectangle{0, 0, w, h},
		stateStack: []subImageState{
			{
				origin: Point{0, 0},
				bounds: Rectangle{0, 0, w, h},
				clip:   Rectangle{0, 0, w, h},
			},
		},
	}
}

func (img *Image) currentState() *subImageState {
	return &img.stateStack[len(img.stateStack)-1]
}

func (img *Image) PushSubImage(r Rectangle) {
	curr := img.currentState()
	absOrigin := curr.origin.Add(Point{r.X(), r.Y()})
	absRect := Rectangle{absOrigin.X(), absOrigin.Y(), r.Width(), r.Height()}
	newClip := curr.clip.Overlap(absRect)

	img.stateStack = append(img.stateStack, subImageState{
		origin: absOrigin,
		bounds: absRect,
		clip:   newClip,
	})
}

func (img *Image) PopSubImage() {
	if len(img.stateStack) > 1 {
		img.stateStack = img.stateStack[:len(img.stateStack)-1]
	}
}

func (img *Image) PushClip(r Rectangle) {
	curr := img.currentState()
	absRect := r.Add(curr.origin)
	newClip := curr.clip.Overlap(absRect)
	curr.clippingStack = append(curr.clippingStack, curr.clip)
	curr.clip = newClip
}

func (img *Image) PopClip() {
	curr := img.currentState()
	if len(curr.clippingStack) > 0 {
		curr.clip = curr.clippingStack[len(curr.clippingStack)-1]
		curr.clippingStack = curr.clippingStack[:len(curr.clippingStack)-1]
	}
}

func (img *Image) Set(p Point, c Color) {
	curr := img.currentState()
	absP := curr.origin.Add(p)

	if !curr.clip.Contains(absP) {
		return
	}

	offset := absP.Y()*img.Stride + absP.X()*4
	img.Pix[offset] = c.R()
	img.Pix[offset+1] = c.G()
	img.Pix[offset+2] = c.B()
	img.Pix[offset+3] = c.A()
}

func (img *Image) At(p Point) Color {
	curr := img.currentState()
	absP := curr.origin.Add(p)

	if !curr.clip.Contains(absP) {
		return Color{}
	}

	offset := absP.Y()*img.Stride + absP.X()*4
	return Color{
		img.Pix[offset],
		img.Pix[offset+1],
		img.Pix[offset+2],
		img.Pix[offset+3],
	}
}

func (img *Image) Fill(c Color) {
	curr := img.currentState()
	clip := curr.clip

	for y := clip.Y(); y < clip.Y()+clip.Height(); y++ {
		for x := clip.X(); x < clip.X()+clip.Width(); x++ {
			offset := y*img.Stride + x*4
			img.Pix[offset] = c.R()
			img.Pix[offset+1] = c.G()
			img.Pix[offset+2] = c.B()
			img.Pix[offset+3] = c.A()
		}
	}
}

func (img *Image) HLine(y, x1, x2, h int, c Color) {
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	img.fillAbsRect(Rectangle{x1, y, x2 - x1, h}, c)
}

func (img *Image) VLine(x, y1, y2, w int, c Color) {
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	img.fillAbsRect(Rectangle{x, y1, w, y2 - y1}, c)
}

func (img *Image) fillAbsRect(r Rectangle, c Color) {
	curr := img.currentState()
	absRect := r.Add(curr.origin)
	drawRect := curr.clip.Overlap(absRect)

	if drawRect.Width() <= 0 || drawRect.Height() <= 0 {
		return
	}

	for y := drawRect.Y(); y < drawRect.Y()+drawRect.Height(); y++ {
		for x := drawRect.X(); x < drawRect.X()+drawRect.Width(); x++ {
			offset := y*img.Stride + x*4
			img.Pix[offset] = c.R()
			img.Pix[offset+1] = c.G()
			img.Pix[offset+2] = c.B()
			img.Pix[offset+3] = c.A()
		}
	}
}

func (img *Image) Border(c Color) {
	curr := img.currentState()
	w, h := curr.bounds.Width(), curr.bounds.Height()

	img.HLine(0, 0, w, 1, c)
	img.HLine(h-1, 0, w, 1, c)
	img.VLine(0, 0, h, 1, c)
	img.VLine(w-1, 0, h, 1, c)
}

func (img *Image) Text(p Point, c Color, f font.Face, wrap bool, s string, args ...any) {
	msg := fmt.Sprintf(s, args...)
	curr := img.currentState()

	relClip := curr.clip.Sub(curr.origin)
	maxWidth := relClip.X() + relClip.Width() - p.X()

	if maxWidth <= 0 {
		return
	}

	var lines []string
	if wrap {
		lines = wrapText(msg, f, maxWidth)
	} else {
		lines = strings.Split(msg, "\n")
	}

	metrics := f.Metrics()
	lineHeight := (metrics.Ascent + metrics.Descent).Ceil()
	totalHeight := len(lines) * lineHeight

	tempImg := image.NewRGBA(image.Rect(0, 0, maxWidth, totalHeight))

	d := &font.Drawer{
		Dst:  tempImg,
		Src:  image.NewUniform(color.RGBA{c.R(), c.G(), c.B(), c.A()}),
		Face: f,
		Dot:  fixed.P(0, metrics.Ascent.Ceil()),
	}

	for _, line := range lines {
		d.DrawString(line)
		d.Dot.X = 0
		d.Dot.Y += fixed.I(lineHeight)
	}

	bounds := tempImg.Bounds()
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			offset := tempImg.PixOffset(x, y)
			sa := tempImg.Pix[offset+3]

			if sa == 0 {
				continue
			}

			targetP := p.Add(Point{x, y})
			absP := curr.origin.Add(targetP)

			if !curr.clip.Contains(absP) {
				continue
			}

			sr := tempImg.Pix[offset]
			sg := tempImg.Pix[offset+1]
			sb := tempImg.Pix[offset+2]

			bg := img.At(targetP)

			a := float64(sa) / 255.0
			outR := uint8(float64(sr)*a + float64(bg.R())*(1-a))
			outG := uint8(float64(sg)*a + float64(bg.G())*(1-a))
			outB := uint8(float64(sb)*a + float64(bg.B())*(1-a))
			outA := uint8(float64(sa) + float64(bg.A())*(1-a))

			img.Set(targetP, Color{outR, outG, outB, outA})
		}
	}
}

func (img *Image) DrawImageScaled(src image.Image, p Point, scale int) {
	curr := img.currentState()
	bounds := src.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()
	dstW, dstH := srcW*scale, srcH*scale

	// Calculate intersection of destination rect and clip rect
	dstRect := Rectangle{p.X(), p.Y(), dstW, dstH}
	absDstRect := dstRect.Add(curr.origin)
	drawRect := curr.clip.Overlap(absDstRect)

	if drawRect.Width() <= 0 || drawRect.Height() <= 0 {
		return
	}

	// Fast path for *image.RGBA (common case)
	srcRGBA, ok := src.(*image.RGBA)

	for y := drawRect.Y(); y < drawRect.Y()+drawRect.Height(); y++ {
		// Map back to source Y
		relY := y - absDstRect.Y()
		srcY := relY / scale

		// Bounds check for source (should be covered by logic, but safe is good)
		if srcY < 0 || srcY >= srcH {
			continue
		}

		for x := drawRect.X(); x < drawRect.X()+drawRect.Width(); x++ {
			// Map back to source X
			relX := x - absDstRect.X()
			srcX := relX / scale

			if srcX < 0 || srcX >= srcW {
				continue
			}

			var r, g, b, a uint8
			if ok {
				// Direct access
				offset := (bounds.Min.Y+srcY)*srcRGBA.Stride + (bounds.Min.X+srcX)*4
				r = srcRGBA.Pix[offset]
				g = srcRGBA.Pix[offset+1]
				b = srcRGBA.Pix[offset+2]
				a = srcRGBA.Pix[offset+3]
			} else {
				// General case
				c := src.At(bounds.Min.X+srcX, bounds.Min.Y+srcY)
				r16, g16, b16, a16 := c.RGBA()
				r, g, b, a = uint8(r16>>8), uint8(g16>>8), uint8(b16>>8), uint8(a16>>8)
			}

			if a == 0 {
				continue
			}

			dstOffset := y*img.Stride + x*4

			// Alpha blending
			if a == 255 {
				img.Pix[dstOffset] = r
				img.Pix[dstOffset+1] = g
				img.Pix[dstOffset+2] = b
				img.Pix[dstOffset+3] = a
			} else {
				bgR := img.Pix[dstOffset]
				bgG := img.Pix[dstOffset+1]
				bgB := img.Pix[dstOffset+2]
				bgA := img.Pix[dstOffset+3]

				fa := float64(a) / 255.0
				outR := uint8(float64(r)*fa + float64(bgR)*(1-fa))
				outG := uint8(float64(g)*fa + float64(bgG)*(1-fa))
				outB := uint8(float64(b)*fa + float64(bgB)*(1-fa))
				outA := uint8(float64(a) + float64(bgA)*(1-fa))

				img.Pix[dstOffset] = outR
				img.Pix[dstOffset+1] = outG
				img.Pix[dstOffset+2] = outB
				img.Pix[dstOffset+3] = outA
			}
		}
	}
}

func wrapText(text string, f font.Face, maxWidth int) []string {
	var lines []string
	for _, paragraph := range strings.Split(text, "\n") {
		words := strings.Fields(paragraph)
		if len(words) == 0 {
			lines = append(lines, "")
			continue
		}

		currentLine := words[0]
		for _, word := range words[1:] {
			width := font.MeasureString(f, currentLine+" "+word).Ceil()
			if width <= maxWidth {
				currentLine += " " + word
			} else {
				lines = append(lines, currentLine)
				currentLine = word
			}
		}
		lines = append(lines, currentLine)
	}
	return lines
}
