package ilda

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestRead(t *testing.T) {
	files, err := filepath.Glob("testdata/*.ild")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			f, err := os.Open(file)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			d := NewDecoder(f)
			for d.Next() {
				frame := d.Frame()
				if !testing.Short() {
					t.Log("Frame", frame.Name, frame.Company, frame.Number, frame.Total)
					for _, point := range frame.Points {
						t.Log("Point", point)
					}
				}
			}
			if err := d.Err(); err != nil {
				t.Error("Err", d.Err())
			}
		})
	}
}

func ExampleDecoder() {
	fd, err := os.Open("testdata/ildatest.ild")
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()
	d := NewDecoder(fd)
	for d.Next() {
		f := d.Frame()
		fmt.Println("Name", f.Name)
		fmt.Println("Company", f.Company)
		fmt.Println("Number", f.Number)
		fmt.Println("Total", f.Total)
		fmt.Println("Projector", f.Projector)
		fmt.Println("Points", len(f.Points))
	}
	if err := d.Err(); err != nil {
		log.Fatal(err)
	}
	// Output:
	// Name ILDA Tes
	// Company t patter
	// Number 0
	// Total 1
	// Projector 0
	// Points 1191
}
