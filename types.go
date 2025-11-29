package q2d

import "math"

// Point represents an X and Y position.
type Point [2]int

// X returns the x coordinate.
func (p Point) X() int { return p[0] }

// Y returns the y coordinate.
func (p Point) Y() int { return p[1] }

// Add returns the sum of p and q.
func (p Point) Add(q Point) Point {
	return Point{p[0] + q[0], p[1] + q[1]}
}

// Sub returns the difference of p - q.
func (p Point) Sub(q Point) Point {
	return Point{p[0] - q[0], p[1] - q[1]}
}

// Mul returns p multiplied by k.
func (p Point) Mul(k int) Point {
	return Point{p[0] * k, p[1] * k}
}

// Div returns p divided by k.
func (p Point) Div(k int) Point {
	return Point{p[0] / k, p[1] / k}
}

// Rectangle represents an X and Y position plus Width and Height.
type Rectangle [4]int

// X returns the x coordinate.
func (r Rectangle) X() int { return r[0] }

// Y returns the y coordinate.
func (r Rectangle) Y() int { return r[1] }

// Width returns the width.
func (r Rectangle) Width() int { return r[2] }

// Height returns the height.
func (r Rectangle) Height() int { return r[3] }

// Add returns the rectangle translated by p.
func (r Rectangle) Add(p Point) Rectangle {
	return Rectangle{r[0] + p[0], r[1] + p[1], r[2], r[3]}
}

// Sub returns the rectangle translated by -p.
func (r Rectangle) Sub(p Point) Rectangle {
	return Rectangle{r[0] - p[0], r[1] - p[1], r[2], r[3]}
}

// Contains returns true if p is within r.
func (r Rectangle) Contains(p Point) bool {
	return p[0] >= r[0] && p[0] < r[0]+r[2] &&
		p[1] >= r[1] && p[1] < r[1]+r[3]
}

// Overlap returns the overlapping area of r and s.
// If there is no overlap, the returned rectangle will have 0 or negative width/height.
func (r Rectangle) Overlap(s Rectangle) Rectangle {
	x1 := max(r[0], s[0])
	y1 := max(r[1], s[1])
	x2 := min(r[0]+r[2], s[0]+s[2])
	y2 := min(r[1]+r[3], s[1]+s[3])

	return Rectangle{x1, y1, max(0, x2-x1), max(0, y2-y1)}
}

// Color represents R, G, B, and A values.
type Color [4]uint8

// R returns the red component.
func (c Color) R() uint8 { return c[0] }

// G returns the green component.
func (c Color) G() uint8 { return c[1] }

// B returns the blue component.
func (c Color) B() uint8 { return c[2] }

// A returns the alpha component.
func (c Color) A() uint8 { return c[3] }

// ToHSL converts the color to Hue, Saturation, Lightness.
// h is in [0, 360), s and l are in [0, 1].
func (c Color) ToHSL() (h, s, l float64) {
	r := float64(c[0]) / 255.0
	g := float64(c[1]) / 255.0
	b := float64(c[2]) / 255.0

	maxVal := math.Max(r, math.Max(g, b))
	minVal := math.Min(r, math.Min(g, b))

	l = (maxVal + minVal) / 2

	if maxVal == minVal {
		h = 0
		s = 0
	} else {
		d := maxVal - minVal
		if l > 0.5 {
			s = d / (2 - maxVal - minVal)
		} else {
			s = d / (maxVal + minVal)
		}

		switch maxVal {
		case r:
			h = (g - b) / d
			if g < b {
				h += 6
			}
		case g:
			h = (b-r)/d + 2
		case b:
			h = (r-g)/d + 4
		}
		h /= 6
	}

	return h * 360, s, l
}

// FromHSL creates a Color from Hue, Saturation, Lightness and Alpha.
func FromHSL(h, s, l float64, a uint8) Color {
	var r, g, b float64

	if s == 0 {
		r, g, b = l, l, l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q

		r = hueToRGB(p, q, h/360+1.0/3.0)
		g = hueToRGB(p, q, h/360)
		b = hueToRGB(p, q, h/360-1.0/3.0)
	}

	return Color{
		uint8(math.Round(r * 255)),
		uint8(math.Round(g * 255)),
		uint8(math.Round(b * 255)),
		a,
	}
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

// Darken returns a darker version of the color using HSL.
// factor should be between 0 and 1.
func (c Color) Darken(factor float64) Color {
	h, s, l := c.ToHSL()
	l = math.Max(0, l*(1-factor))
	return FromHSL(h, s, l, c[3])
}

// Lighten returns a lighter version of the color using HSL.
// factor should be between 0 and 1.
func (c Color) Lighten(factor float64) Color {
	h, s, l := c.ToHSL()
	l = math.Min(1, l+(1-l)*factor)
	return FromHSL(h, s, l, c[3])
}

// AdjustHue shifts the hue by degrees.
func (c Color) AdjustHue(degrees float64) Color {
	h, s, l := c.ToHSL()
	h += degrees
	for h < 0 {
		h += 360
	}
	for h >= 360 {
		h -= 360
	}
	return FromHSL(h, s, l, c[3])
}

// AdjustSaturation multiplies the saturation by factor.
func (c Color) AdjustSaturation(factor float64) Color {
	h, s, l := c.ToHSL()
	s = math.Max(0, math.Min(1, s*factor))
	return FromHSL(h, s, l, c[3])
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
