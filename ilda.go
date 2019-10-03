// Package ilda implements ILDA Image Data Transfer Format
package ilda

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"io"
)

const (
	indexedColor3D = iota // 3D Coordinates with Indexed Color
	indexedColor2D        // 2D Coordinates with Indexed Color
	colorPalette          // Color Palette
	_                     // not used
	trueColor3D           // 3D Coordinates with True Color
	trueColor2D           // 2D Coordinates with True Color
)

// magic identifying an ILDA format header
var magic = [4]byte{'I', 'L', 'D', 'A'}

type ILDA struct {
	Frames  []Frame
	Palette color.Palette
}

var (
	ErrMagic  = errors.New("bad magic")
	ErrFormat = errors.New("invalid format")
)

func read(r io.Reader, v interface{}) error {
	return binary.Read(r, binary.BigEndian, v)
}

func readHeader(r io.Reader) (Header, error) {
	var h Header
	if err := read(r, &h); err != nil {
		return Header{}, err
	}
	if h.Magic != magic {
		return Header{}, ErrMagic
	}
	return h, nil
}

func readData(r io.Reader, code uint8) (Data, error) {
	switch code {
	case indexedColor3D:
		var d IndexedColor3D
		if err := read(r, &d); err != nil {
			return nil, err
		}
		return d, nil
	case indexedColor2D:
		var d IndexedColor2D
		if err := read(r, &d); err != nil {
			return nil, err
		}
		return d, nil
	case colorPalette:
		var d ColorPalette
		if err := read(r, &d); err != nil {
			return nil, err
		}
		return d, nil
	case trueColor3D:
		var d TrueColor3D
		if err := read(r, &d); err != nil {
			return nil, err
		}
		return d, nil
	case trueColor2D:
		var d TrueColor2D
		if err := read(r, &d); err != nil {
			return nil, err
		}
		return d, nil
	default:
		return nil, ErrFormat
	}
}

func Read(r io.Reader) (ILDA, error) {
	var l ILDA
	l.Palette = DefaultPalette
	for {
		var f Frame
		h, err := readHeader(r)
		if err != nil {
			return ILDA{}, err
		}
		if h.RecordsNumber == 0 {
			break
		}
		f.Header = h
		for i := 0; i < int(h.RecordsNumber); i++ {
			d, err := readData(r, h.FormatCode)
			if err != nil {
				return ILDA{}, err
			}
			f.Data = append(f.Data, d)
		}
		if h.FormatCode == colorPalette {
			for _, d := range f.Data {
				l.Palette = append(l.Palette, d.Color(nil))
			}
		} else {
			l.Frames = append(l.Frames, f)
		}
	}
	return l, nil
}

type Flags uint8

// Status Codes
const (
	LastPoint Flags = 1 << 7
	Blanking  Flags = 1 << 6
)

// Coordinate Data
//
//  X: Extreme left: -32768, extreme right: +32767
//  Y: Extreme bottom: -32768, extreme top: +32767
//  Z: Extreme rear: -32768, extreme front: +32767
//  rear: away from viewer, behind screen
//  front: towards viewer, in front of screen
//
type Data interface {
	Point() image.Point
	Depth() int
	Color(color.Palette) color.Color
	Status(Flags) bool
}

type Frame struct {
	Header Header
	Data   []Data
}

type Header struct {
	Magic           [4]byte  // "ILDA"
	_               [3]uint8 // Reserved, all zero
	FormatCode      uint8    // Format Code
	FrameName       [8]byte  // Frame of Color Palette Name
	CompanyName     [8]byte  // Company Name
	RecordsNumber   uint16   // Number of Records
	FrameNumber     uint16   // Frame or Color Palette Number
	TotalFrames     uint16   // Total Frames in Sequence or 0
	ProjectorNumber uint8    // Projector Number
	_               uint8    // Reserved for future, zero
}

func trimZero(b []byte) string {
	n := bytes.IndexByte(b, 0)
	if n < 0 {
		n = len(b)
	}
	return string(b[:n])
}

func (h Header) Name() string {
	return trimZero(h.FrameName[:])
}

func (h Header) Company() string {
	return trimZero(h.CompanyName[:])
}

// IndexedColor3D – 3D Coordinates with Indexed Color
type IndexedColor3D struct {
	X, Y, Z    int16
	StatusCode uint8
	ColorIndex uint8
}

func (f IndexedColor3D) Point() image.Point                { return image.Pt(int(f.X), int(f.Y)) }
func (f IndexedColor3D) Depth() int                        { return int(f.Z) }
func (f IndexedColor3D) Color(p color.Palette) color.Color { return p[int(f.ColorIndex)] }
func (f IndexedColor3D) Status(v Flags) bool               { return f.StatusCode&uint8(v) != 0 }

// IndexedColor2D – 2D Coordinates with Indexed Color
type IndexedColor2D struct {
	X, Y       int16
	StatusCode uint8
	ColorIndex uint8
}

func (f IndexedColor2D) Point() image.Point                { return image.Pt(int(f.X), int(f.Y)) }
func (f IndexedColor2D) Depth() int                        { return 0 }
func (f IndexedColor2D) Color(p color.Palette) color.Color { return p[int(f.ColorIndex)] }
func (f IndexedColor2D) Status(v Flags) bool               { return f.StatusCode&uint8(v) != 0 }

// ColorPalette – Color Palette
type ColorPalette struct {
	R, G, B uint8
}

func (f ColorPalette) Point() image.Point              { return image.Pt(0, 0) }
func (f ColorPalette) Depth() int                      { return 0 }
func (f ColorPalette) Color(color.Palette) color.Color { return color.RGBA{f.R, f.G, f.B, 255} }
func (f ColorPalette) Status(v Flags) bool             { return false }

// TrueColor3D – 3D Coordinates with True Color
type TrueColor3D struct {
	X, Y, Z    int16
	StatusCode uint8
	B, G, R    uint8
}

func (f TrueColor3D) Point() image.Point              { return image.Pt(int(f.X), int(f.Y)) }
func (f TrueColor3D) Depth() int                      { return int(f.Z) }
func (f TrueColor3D) Color(color.Palette) color.Color { return color.RGBA{f.R, f.G, f.B, 255} }
func (f TrueColor3D) Status(v Flags) bool             { return f.StatusCode&uint8(v) != 0 }

// TrueColor2D – 2D Coordinates with True Color
type TrueColor2D struct {
	X, Y       int16
	StatusCode uint8
	B, G, R    uint8
}

func (f TrueColor2D) Point() image.Point              { return image.Pt(int(f.X), int(f.Y)) }
func (f TrueColor2D) Depth() int                      { return 0 }
func (f TrueColor2D) Color(color.Palette) color.Color { return color.RGBA{f.R, f.G, f.B, 255} }
func (f TrueColor2D) Status(v Flags) bool             { return f.StatusCode&uint8(v) != 0 }
