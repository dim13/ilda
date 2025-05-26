package ilda

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDraw(t *testing.T) {
	files, err := filepath.Glob("testdata/*.ild")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			ild, err := os.Open(file)
			if err != nil {
				t.Fatal(err)
			}
			defer ild.Close()

			frames, err := NewDecoder(ild).AllFrames()
			if err != nil {
				t.Fatal(err)
			}
			if testing.Short() && len(frames) > 1 {
				t.Skip("skip animation")
			}

			bg := image.NewUniform(color.Black)
			pal := append(DefaultPalette, color.Black, color.Transparent)
			r := image.Rect(0, 0, 640, 640)
			images := make([]*image.Paletted, len(frames))
			delays := make([]int, len(frames))
			for i, frame := range frames {
				images[i] = image.NewPaletted(r, pal)
				delays[i] = 4 // 25Hz
				frame.Draw(images[i], r, bg, image.Point{})
			}

			out, err := os.Create(strings.TrimSuffix(file, filepath.Ext(file)) + ".gif")
			if err != nil {
				t.Fatal(err)
			}
			defer out.Close()
			gif.EncodeAll(out, &gif.GIF{Image: images, Delay: delays})
		})
	}
}
