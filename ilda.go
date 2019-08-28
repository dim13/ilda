// Package ilda implements ILDA Image Data Transfer Format
package ilda

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image/color"
	"io"
)

const (
	format0 = 0 // 3D Coordinates with Indexed Color
	format1 = 1 // 2D Coordinates with Indexed Color
	format2 = 2 // Color Palette
	format4 = 4 // 3D Coordinates with True Color
	format5 = 5 // 2D Coordinates with True Color
)

// magic identifying an ILDA format header
var magic = [4]byte{'I', 'L', 'D', 'A'}

var (
	Palette = DefaultPalette
	Off     = color.Black
)

type ILDA struct {
	Frames []Frame
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
	case format0:
		var d Format0
		if err := read(r, &d); err != nil {
			return nil, err
		}
		return d, nil
	case format1:
		var d Format1
		if err := read(r, &d); err != nil {
			return nil, err
		}
		return d, nil
	case format2:
		var d Format2
		if err := read(r, &d); err != nil {
			return nil, err
		}
		return d, nil
	case format4:
		var d Format4
		if err := read(r, &d); err != nil {
			return nil, err
		}
		return d, nil
	case format5:
		var d Format5
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
		l.Frames = append(l.Frames, f)
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
	Point() (x, y, z int)
	Color() color.Color
	Flag(Flags) bool
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

// Format0 – 3D Coordinates with Indexed Color
type Format0 struct {
	X, Y, Z    int16
	StatusCode uint8
	ColorIndex uint8
}

func (f Format0) Point() (x, y, z int) { return int(f.X), int(f.Y), int(f.Z) }
func (f Format0) Color() color.Color   { return Palette[int(f.ColorIndex)] }
func (f Format0) Flag(v Flags) bool    { return f.StatusCode&uint8(v) != 0 }

// Format1 – 2D Coordinates with Indexed Color
type Format1 struct {
	X, Y       int16
	StatusCode uint8
	ColorIndex uint8
}

func (f Format1) Point() (x, y, z int) { return int(f.X), int(f.Y), 0 }
func (f Format1) Color() color.Color   { return Palette[int(f.ColorIndex)] }
func (f Format1) Flag(v Flags) bool    { return f.StatusCode&uint8(v) != 0 }

// Format2 – Color Palette
type Format2 struct {
	R, G, B uint8
}

func (f Format2) Point() (x, y, z int) { return 0, 0, 0 }
func (f Format2) Color() color.Color   { return color.RGBA{f.R, f.G, f.B, 255} }
func (f Format2) Flag(v Flags) bool    { return false }

// Format4 – 3D Coordinates with True Color
type Format4 struct {
	X, Y, Z    int16
	StatusCode uint8
	B, G, R    uint8
}

func (f Format4) Point() (x, y, z int) { return int(f.X), int(f.Y), int(f.Z) }
func (f Format4) Color() color.Color   { return color.RGBA{f.R, f.G, f.B, 255} }
func (f Format4) Flag(v Flags) bool    { return f.StatusCode&uint8(v) != 0 }

// Format5 – 2D Coordinates with True Color
type Format5 struct {
	X, Y       int16
	StatusCode uint8
	B, G, R    uint8
}

func (f Format5) Point() (x, y, z int) { return int(f.X), int(f.Y), 0 }
func (f Format5) Color() color.Color   { return color.RGBA{f.R, f.G, f.B, 255} }
func (f Format5) Flag(v Flags) bool    { return f.StatusCode&uint8(v) != 0 }
