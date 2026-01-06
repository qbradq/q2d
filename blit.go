package q2d

func (img *Image) Blit(src *Image, p Point) {
	dstCurr := img.currentState()
	srcCurr := src.currentState()

	w, h := srcCurr.bounds.Width(), srcCurr.bounds.Height()

	// Destination (logical) rect
	dstRect := Rectangle{p.X(), p.Y(), w, h}
	// Destination (absolute) rect
	absDstRect := dstRect.Add(dstCurr.origin)
	// Clip to destination
	drawRect := dstCurr.clip.Overlap(absDstRect)

	if drawRect.Width() <= 0 || drawRect.Height() <= 0 {
		return
	}

	for y := drawRect.Y(); y < drawRect.Y()+drawRect.Height(); y++ {
		// Map back to source logical coordinates
		relY := y - absDstRect.Y()

		// Map to source absolute coordinates
		srcAbsY := srcCurr.origin.Y() + relY

		// Check source clipping
		if srcAbsY < srcCurr.clip.Y() || srcAbsY >= srcCurr.clip.Y()+srcCurr.clip.Height() {
			continue
		}

		for x := drawRect.X(); x < drawRect.X()+drawRect.Width(); x++ {
			// Map back to source logical coordinates
			relX := x - absDstRect.X()

			// Map to source absolute coordinates
			srcAbsX := srcCurr.origin.X() + relX

			// Check source clipping
			// Note: srcCurr.clip is absolute
			if srcAbsX < srcCurr.clip.X() || srcAbsX >= srcCurr.clip.X()+srcCurr.clip.Width() {
				continue
			}

			// Read Source
			srcOffset := srcAbsY*src.Stride + srcAbsX*4
			sr := src.Pix[srcOffset]
			sg := src.Pix[srcOffset+1]
			sb := src.Pix[srcOffset+2]
			sa := src.Pix[srcOffset+3]

			if sa == 0 {
				continue
			}

			// Destination Offset
			dstOffset := y*img.Stride + x*4

			if sa == 255 {
				img.Pix[dstOffset] = sr
				img.Pix[dstOffset+1] = sg
				img.Pix[dstOffset+2] = sb
				img.Pix[dstOffset+3] = sa
			} else {
				bgR := img.Pix[dstOffset]
				bgG := img.Pix[dstOffset+1]
				bgB := img.Pix[dstOffset+2]
				bgA := img.Pix[dstOffset+3]

				a := float64(sa) / 255.0
				outR := uint8(float64(sr)*a + float64(bgR)*(1-a))
				outG := uint8(float64(sg)*a + float64(bgG)*(1-a))
				outB := uint8(float64(sb)*a + float64(bgB)*(1-a))
				outA := uint8(float64(sa) + float64(bgA)*(1-a))

				img.Pix[dstOffset] = outR
				img.Pix[dstOffset+1] = outG
				img.Pix[dstOffset+2] = outB
				img.Pix[dstOffset+3] = outA
			}
		}
	}
}
