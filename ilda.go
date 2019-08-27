// Package ilda implements ILDA Image Data Transfer Format
package ilda

import "image/color"

// Magic identifying an ILDA format header
const Magic = "ILDA"

const (
	FormatCode0 = 0 // 3D Coordinates with Indexed Color
	FormatCode1 = 1 // 2D Coordinates with Indexed Color
	FormatCode2 = 2 // Color Palette
	FormatCode4 = 4 // 3D Coordinates with True Color
	FormatCode5 = 5 // 2D Coordinates with True Color
)

var Palette = DefaultPalette

type ILDA struct {
	Frames []Frame
}

type Data interface {
	Point() (x, y, z int)
	Color() color.Color
}

type Frame struct {
	Header Header
	Data   []Data
}

type Header struct {
	Magic           [4]byte  // "ILDA"
	_               [3]uint8 // Reserved
	FormatCode      uint8    // Format Code
	FrameName       [8]byte  // Frame of Color Palette Name
	CompanyName     [8]byte  // Company Name
	RecordsNumber   uint16   // Number of Records
	FrameNumber     uint16   // Frame or Color Palette Number
	TotalFrames     uint16   // Total Frames in Sequence or 0
	ProjectorNumber uint8    // Projector Number
	_               uint8    // Reserved
}

// Status Codes
const (
	LastPoint = 1 << 7
	Blanking  = 1 << 6
)

// Format0 – 3D Coordinates with Indexed Color
type Format0 struct {
	X, Y, Z    int16
	StatusCode uint8
	ColorIndex uint8
}

func (f Format0) Point() (x, y, z int) {
	return int(f.X), int(f.Y), int(f.Z)
}

func (f Format0) Color() color.Color {
	if f.StatusCode&Blanking != 0 {
		return color.Black
	}
	return Palette[int(f.ColorIndex)]
}

// Format1 – 2D Coordinates with Indexed Color
type Format1 struct {
	X, Y       int16
	StatusCode uint8
	ColorIndex uint8
}

func (f Format1) Point() (x, y, z int) {
	return int(f.X), int(f.Y), 0
}

func (f Format1) Color() color.Color {
	if f.StatusCode&Blanking != 0 {
		return color.Black
	}
	return Palette[int(f.ColorIndex)]
}

// Format2 – Color Palette
type Format2 struct {
	R, G, B uint8
}

func (f Format2) Point() (x, y, z int) {
	return 0, 0, 0
}

func (f Format2) Color() color.Color {
	return color.RGBA{f.R, f.G, f.B, 255}
}

// Format4 – 3D Coordinates with True Color
type Format4 struct {
	X, Y, Z    int16
	StatusCode uint8
	B, G, R    uint8
}

func (f Format4) Point() (x, y, z int) {
	return int(f.X), int(f.Y), int(f.Z)
}

func (f Format4) Color() color.Color {
	if f.StatusCode&Blanking != 0 {
		return color.Black
	}
	return color.RGBA{f.R, f.G, f.B, 255}
}

// Format5 – 2D Coordinates with True Color
type Format5 struct {
	X, Y       int16
	StatusCode uint8
	B, G, R    uint8
}

func (f Format5) Point() (x, y, z int) {
	return int(f.X), int(f.Y), 0
}

func (f Format5) Color() color.Color {
	if f.StatusCode&Blanking != 0 {
		return color.Black
	}
	return color.RGBA{f.R, f.G, f.B, 255}
}
