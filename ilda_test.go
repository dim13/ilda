package ilda

import (
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	testCases := []string{
		"testdata/ildatest.ild",
		"testdata/ildatstb.ild",
		"testdata/30k.ild",
		"testdata/barney.ild",
		"testdata/biker.ild",
		//"testdata/theriddle.ild",
	}
	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			f, err := os.Open(tc)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			d := New(f)
			for d.Next() {
				frame := d.Frame()
				t.Log("Frame", frame.Name, frame.Company, frame.Number, frame.Total)
				for _, point := range frame.Points {
					t.Log("Point", point)
				}
			}
			if err := d.Err(); err != nil {
				t.Error("Err", d.Err())
			}
		})
	}
}
