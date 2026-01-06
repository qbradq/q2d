package q2d

import (
	"testing"
)

func TestImage(t *testing.T) {
	img := NewImage(100, 100)

	// Test Set/At
	red := Color{255, 0, 0, 255}
	img.Set(Point{10, 10}, red)
	c := img.At(Point{10, 10})
	if c != red {
		t.Errorf("Expected %v, got %v", red, c)
	}

	// Test SubImage
	img.PushSubImage(Rectangle{20, 20, 50, 50})
	// Origin is now (20, 20). Clip is (20, 20, 50, 50).

	// Set (0,0) -> Absolute (20, 20)
	blue := Color{0, 0, 255, 255}
	img.Set(Point{0, 0}, blue)

	// Check absolute (20, 20)
	// We need to pop to check absolute? Or check via internal Pix?
	// Let's pop.
	img.PopSubImage()
	c = img.At(Point{20, 20})
	if c != blue {
		t.Errorf("Expected %v, got %v", blue, c)
	}

	// Test Clipping
	img.PushSubImage(Rectangle{0, 0, 100, 100})
	img.PushClip(Rectangle{10, 10, 10, 10}) // Clip is (10, 10, 10, 10)

	// Set at (5, 5) -> Absolute (5, 5). Outside clip.
	img.Set(Point{5, 5}, red)
	c = img.At(Point{5, 5})
	if c == red {
		t.Errorf("Should be clipped")
	}

	// Set at (15, 15) -> Absolute (15, 15). Inside clip.
	img.Set(Point{15, 15}, red)
	c = img.At(Point{15, 15})
	if c != red {
		t.Errorf("Should be drawn")
	}

	img.PopClip()
	// Clip restored to (0, 0, 100, 100)
	img.Set(Point{5, 5}, red)
	c = img.At(Point{5, 5})
	if c != red {
		t.Errorf("Should be drawn after pop clip")
	}
}

func TestBlit(t *testing.T) {
	src := NewImage(20, 20)
	red := Color{255, 0, 0, 255}
	src.Fill(red)

	dst := NewImage(50, 50)

	// Blit at (10, 10)
	dst.Blit(src, Point{10, 10})

	// Check (10, 10) -> should be red
	c := dst.At(Point{10, 10})
	if c != red {
		t.Errorf("Expected red at 10,10, got %v", c)
	}

	// Check (9, 9) -> should be transparent/black
	c = dst.At(Point{9, 9})
	empty := Color{}
	if c != empty {
		t.Errorf("Expected empty at 9,9, got %v", c)
	}

	// Check (29, 29) -> red (10+19, 10+19)
	c = dst.At(Point{29, 29})
	if c != red {
		t.Errorf("Expected red at 29,29, got %v", c)
	}

	// Test Clipping in Blit
	dst.PushClip(Rectangle{0, 0, 15, 15})
	dst.Fill(empty) // Clear

	dst.Blit(src, Point{10, 10})

	// (14, 14) inside clip -> red
	if dst.At(Point{14, 14}) != red {
		t.Errorf("Expected red at 14,14 (clipped)")
	}
	// (15, 15) outside clip -> empty
	if dst.At(Point{15, 15}) != empty {
		t.Errorf("Expected empty at 15,15 (clipped)")
	}
}
