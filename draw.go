package ilda

import (
	"image"
	"image/color"
	"image/draw"
)

// A Point is an X, Y coordinate pair. The axes increase right and down.
//
//	X: Extreme left: 0, extreme right max.X
//	Y: Extreme top: 0, extreme bottom max.Y
//
// ILDA Point
//
//	X: Extreme left: -32768, extreme right: +32767
//	Y: Extreme bottom: -32768, extreme top: +32767
//
//	[-32768, +32767] -> [min.X, min.Y] (0, 0)
//	[+32767, -32768] -> [max.X, max.Y] (65535, 65535)
func (p Point) Point(r image.Rectangle) image.Point {
	pp := image.Pt(int(p.X), -int(p.Y)).Add(image.Pt(0x8000, 0x8000))
	pp.X, pp.Y = pp.X*r.Dx(), pp.Y*r.Dy()
	return pp.Div(0x10000).Add(image.Pt(1, -1))
}

// Draw aligns r.Min in dst with sp in src and then replaces the
// rectangle r in dst with the result of drawing src on dst.
func (f *Frame) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	// copy background
	draw.Draw(dst, r, src, sp, draw.Src)
	var plt plot
	for _, pt := range f.Points {
		switch pt.Color {
		case color.Transparent:
			plt.move(pt.Point(r))
		default:
			plt.draw(dst, pt.Point(r), pt.Color)
		}
	}
}

// DrawServo renders like Draw, but first tracks the frame's points with a
// simulated galvo pair (x, y) at the given scan rate (points per second), so
// corners show realistic overshoot and settling instead of instant jumps.
func (f *Frame) DrawServo(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point, pps float64, x, y *Servo) {
	moved := f.Move(pps, x, y)
	moved.Draw(dst, r, src, sp)
}

type plot image.Point

func (p0 *plot) move(p1 image.Point) {
	p0.X, p0.Y = p1.X, p1.Y
}

func abs(x int) (int, int) {
	if x < 0 {
		return -x, -1
	}
	return x, 1
}

func (p0 *plot) draw(img draw.Image, p1 image.Point, c color.Color) {
	dx, sx := abs(p1.X - p0.X)
	dy, sy := abs(p1.Y - p0.Y)
	dy = -dy
	err := dx + dy // error value e_xy
	for {
		img.Set(p0.X, p0.Y, c)
		e2 := 2 * err
		if e2 >= dy { // e_xy + e_x > 0
			if p0.X == p1.X {
				break
			}
			err += dy
			p0.X += sx
		}
		if e2 <= dx { // e_xy + e_y < 0
			if p0.Y == p1.Y {
				break
			}
			err += dx
			p0.Y += sy
		}
	}
}
