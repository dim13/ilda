package ilda

import (
	"image"
	"image/color"
	"image/draw"
)

// A Point is an X, Y coordinate pair. The axes increase right and down.
//
//	X: Extreme left: 0, extreme right Inf
//	Y: Extreme top: 0, extreme bottom Inf
//
// ILDA Point
//
//	X: Extreme left: -32768, extreme right: +32767
//	Y: Extreme bottom: -32768, extreme top: +32767
func (p Point) normalize(r image.Rectangle) image.Point {
	dx := 65535 / r.Dx()
	dy := 65535 / r.Dy()
	x := 32768 + int(p.X)
	y := 32767 - int(p.Y)
	return r.Min.Add(image.Pt(x/dx, y/dy))
}

// Draw aligns r.Min in dst with sp in src and then replaces the
// rectangle r in dst with the result of drawing src on dst.
func (f *Frame) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	// copy background
	draw.Draw(dst, r, src, sp, draw.Src)
	var plt plot
	for _, pt := range f.Points {
		plt.drawTo(dst, pt.normalize(r), pt.Color)
	}
}

type plot image.Point

func (m *plot) drawTo(dst draw.Image, p image.Point, c color.Color) {
	dx := p.X - m.X
	if dx < 0 {
		dx = -dx
	}
	sx := -1
	if m.X < p.X {
		sx = 1
	}
	dy := m.Y - p.Y
	if dy > 0 {
		dy = -dy
	}
	sy := -1
	if m.Y < p.Y {
		sy = 1
	}
	err := dx + dy // error value e_xy
	for {
		if c != color.Transparent {
			dst.Set(m.X, m.Y, c)
		}
		if m.X == p.X && m.Y == p.Y {
			break
		}
		e2 := 2 * err
		if e2 >= dy { // e_xy + e_x > 0
			err += dy
			m.X += sx
		}
		if e2 <= dx { // e_xy + e_y < 0
			err += dx
			m.Y += sy
		}
	}
}
