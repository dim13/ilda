package ilda

import (
	"image"
	"image/color"
	"image/draw"
)

// A Point is an X, Y coordinate pair. The axes increase right and down.
//  X: Extreme left: 0, extreme right Inf
//  Y: Extreme top: 0, extreme bottom Inf
//
// ILDA Point
//  X: Extreme left: -32768, extreme right: +32767
//  Y: Extreme bottom: -32768, extreme top: +32767
func (p Point) normalize(r image.Rectangle) image.Point {
	dx := 65535 / float64(r.Max.X-r.Min.X)
	dy := 65535 / float64(r.Max.Y-r.Min.Y)
	x := 32768 + float64(p.X)
	y := 32767 - float64(p.Y)
	return image.Pt(int(x/dx), int(y/dy))
}

// Draw aligns r.Min in dst with sp in src and then replaces the
// rectangle r in dst with the result of drawing src on dst.
func (f *Frame) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	// copy background
	draw.Draw(dst, r, src, sp, draw.Src)
	var plt plot
	for _, pt := range f.Points {
		p := pt.normalize(r)
		plt.drawTo(dst, p.X, p.Y, pt.Color)
	}
}

type plot struct {
	x0, y0 int
}

func (m *plot) drawTo(dst draw.Image, x1, y1 int, c color.Color) {
	dx := x1 - m.x0
	if dx < 0 {
		dx = -dx
	}
	sx := -1
	if m.x0 < x1 {
		sx = 1
	}
	dy := m.y0 - y1
	if dy > 0 {
		dy = -dy
	}
	sy := -1
	if m.y0 < y1 {
		sy = 1
	}
	err := dx + dy // error value e_xy
	for {
		if c != color.Transparent {
			dst.Set(m.x0, m.y0, c)
		}
		if m.x0 == x1 && m.y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy { // e_xy + e_x > 0
			err += dy
			m.x0 += sx
		}
		if e2 <= dx { // e_xy + e_y < 0
			err += dx
			m.y0 += sy
		}
	}
}
