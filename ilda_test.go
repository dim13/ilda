package ilda

import (
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	f, err := os.Open("testdata/ildatest.ild")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	l, err := Read(f)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range l.Frames {
		h := f.Header
		t.Log("frame", h, h.Name(), h.Company())
		for _, d := range f.Data {
			pt := d.Point()
			dp := d.Depth()
			c := d.Color(l.Palette)
			b := d.Status(Blanking)
			t.Log("data", pt, dp, c, b)
		}
	}
}
