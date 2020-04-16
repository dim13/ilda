package ilda

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func TestDraw(t *testing.T) {
	ild, err := os.Open("testdata/ildatest.ild")
	if err != nil {
		t.Fatal(err)
	}
	defer ild.Close()

	frames, err := New(ild).AllFrames()
	if err != nil {
		t.Fatal(err)
	}

	bg := image.Image(image.NewUniform(color.Black))
	dst := image.NewRGBA(image.Rect(0, 0, 1024, 1024))
	for _, frame := range frames {
		frame.Draw(dst, dst.Bounds(), bg, image.ZP)
		bg = dst
	}

	out, err := os.Create("test.png")
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()
	if err := png.Encode(out, dst); err != nil {
		t.Error(err)
	}
}