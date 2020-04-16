// Package ilda implements decoding of ILDA Image Data Transfer Format
package ilda

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image/color"
	"io"
)

// magic identifying an ILDA format header
var magic = [4]byte{'I', 'L', 'D', 'A'}

// Point coordinate
//
//  X: Extreme left: -32768, extreme right: +32767
//  Y: Extreme bottom: -32768, extreme top: +32767
//  Z: Extreme rear: -32768, extreme front: +32767
//  rear: away from viewer, behind screen
//  front: towards viewer, in front of screen
//
type Point struct {
	X, Y, Z int
	color.Color
}

var (
	ErrMagic  = errors.New("bad magic")
	ErrFormat = errors.New("invalid format")
)

func New(r io.Reader) *Decoder {
	return &Decoder{r: r, pal: DefaultPalette}
}

type Decoder struct {
	r    io.Reader
	last bool
	pal  color.Palette
	err  error
}

func (d *Decoder) Next() bool {
	return d.err == nil && !d.last
}

func (d *Decoder) Err() error {
	if d.err == io.EOF {
		return nil
	}
	return d.err
}

func (d *Decoder) AllFrames() ([]Frame, error) {
	var frames []Frame
	for d.Next() {
		frames = append(frames, d.Frame())
	}
	return frames, d.Err()
}

func (d *Decoder) Frame() Frame {
	if d.err != nil {
		return Frame{}
	}
	hdr, err := d.readFrameHeader()
	if err != nil {
		d.err = err
		return Frame{}
	}
	frame := Frame{
		Name:      trimZero(hdr.FrameName[:]),
		Company:   trimZero(hdr.CompanyName[:]),
		Number:    int(hdr.FrameNumber),
		Total:     int(hdr.TotalFrames),
		Projector: int(hdr.ProjectorNumber),
	}
	d.last = frame.Number >= frame.Total-1
	switch hdr.FormatCode {
	case formatIndexedColor3D:
		frame.Points, err = d.readIndexedColor3D(hdr.RecordsNumber)
	case formatIndexedColor2D:
		frame.Points, err = d.readIndexedColor2D(hdr.RecordsNumber)
	case formatColorPalette:
		d.pal, err = d.readPaletteColor(hdr.RecordsNumber)
	case formatTrueColor3D:
		frame.Points, err = d.readTrueColor3D(hdr.RecordsNumber)
	case formatTrueColor2D:
		frame.Points, err = d.readTrueColor2D(hdr.RecordsNumber)
	default:
		err = ErrFormat
	}
	if err != nil {
		d.err = err
		return Frame{}
	}
	// read next frame after palette change
	if hdr.FormatCode == formatColorPalette {
		return d.Frame()
	}
	return frame
}

func trimZero(b []byte) string {
	n := bytes.IndexByte(b, 0)
	if n < 0 {
		n = len(b)
	}
	return string(b[:n])
}

func (d *Decoder) readFrameHeader() (frameHeader, error) {
	var hdr frameHeader
	if err := binary.Read(d.r, binary.BigEndian, &hdr); err != nil {
		return frameHeader{}, err
	}
	if hdr.Magic != magic {
		return frameHeader{}, ErrMagic
	}
	return hdr, nil
}

func (d *Decoder) readIndexedColor3D(n uint16) ([]Point, error) {
	frame := make([]indexedColor3D, n)
	if err := binary.Read(d.r, binary.BigEndian, &frame); err != nil {
		return nil, err
	}
	points := make([]Point, n)
	for i, v := range frame {
		points[i] = Point{int(v.X), int(v.Y), int(v.Z), v.Color(d.pal)}
	}
	return points, nil
}

func (d *Decoder) readIndexedColor2D(n uint16) ([]Point, error) {
	frame := make([]indexedColor2D, n)
	if err := binary.Read(d.r, binary.BigEndian, &frame); err != nil {
		return nil, err
	}
	points := make([]Point, n)
	for i, v := range frame {
		points[i] = Point{int(v.X), int(v.Y), 0, v.Color(d.pal)}
	}
	return points, nil
}

func (d *Decoder) readPaletteColor(n uint16) (color.Palette, error) {
	frame := make([]paletteColor, n)
	if err := binary.Read(d.r, binary.BigEndian, &frame); err != nil {
		return nil, err
	}
	palette := make([]color.Color, n)
	for i, v := range frame {
		palette[i] = v.Color()
	}
	return palette, nil
}

func (d *Decoder) readTrueColor3D(n uint16) ([]Point, error) {
	frame := make([]trueColor3D, n)
	if err := binary.Read(d.r, binary.BigEndian, &frame); err != nil {
		return nil, err
	}
	points := make([]Point, n)
	for i, v := range frame {
		points[i] = Point{int(v.X), int(v.Y), int(v.Z), v.Color()}
	}
	return points, nil
}

func (d *Decoder) readTrueColor2D(n uint16) ([]Point, error) {
	frame := make([]trueColor2D, n)
	if err := binary.Read(d.r, binary.BigEndian, &frame); err != nil {
		return nil, err
	}
	points := make([]Point, n)
	for i, v := range frame {
		points[i] = Point{int(v.X), int(v.Y), 0, v.Color()}
	}
	return points, nil
}

type Frame struct {
	Name      string
	Company   string
	Number    int
	Total     int
	Projector int
	Points    []Point
}

// format codes
const (
	formatIndexedColor3D = iota // 3D Coordinates with Indexed Color
	formatIndexedColor2D        // 2D Coordinates with Indexed Color
	formatColorPalette          // Color Palette
	_                           // Not used
	formatTrueColor3D           // 3D Coordinates with True Color
	formatTrueColor2D           // 2D Coordinates with True Color
)

// status codes
const (
	statusLastPoint = 1 << (7 - iota)
	statusBlanking
)

type frameHeader struct {
	Magic           [4]byte  // "ILDA"
	_               [3]uint8 // Reserved, all zero
	FormatCode      uint8    // Format Code
	FrameName       [8]byte  // Frame or Color Palette Name
	CompanyName     [8]byte  // Company Name
	RecordsNumber   uint16   // Number of Records
	FrameNumber     uint16   // Frame or Color Palette Number
	TotalFrames     uint16   // Total Frames in Sequence or zero
	ProjectorNumber uint8    // Projector Number
	_               uint8    // Reserved for future, zero
}

// indexedColor3D – 3D Coordinates with Indexed Color
type indexedColor3D struct {
	X, Y, Z    int16
	StatusCode uint8
	ColorIndex uint8
}

func (v indexedColor3D) Color(pal color.Palette) color.Color {
	if v.StatusCode&statusBlanking != 0 {
		return color.Transparent
	}
	return pal[v.ColorIndex]
}

// indexedColor2D – 2D Coordinates with Indexed Color
type indexedColor2D struct {
	X, Y       int16
	StatusCode uint8
	ColorIndex uint8
}

func (v indexedColor2D) Color(pal color.Palette) color.Color {
	if v.StatusCode&statusBlanking != 0 {
		return color.Transparent
	}
	return pal[v.ColorIndex]
}

// paletteColor – Color Palette
type paletteColor struct {
	R, G, B uint8
}

func (v paletteColor) Color() color.Color {
	return color.RGBA{v.R, v.G, v.B, 0xff}
}

// trueColor3D – 3D Coordinates with True Color
type trueColor3D struct {
	X, Y, Z    int16
	StatusCode uint8
	B, G, R    uint8
}

func (v trueColor3D) Color() color.Color {
	if v.StatusCode&statusBlanking != 0 {
		return color.Transparent
	}
	return color.RGBA{v.R, v.G, v.B, 0xff}
}

// trueColor2D – 2D Coordinates with True Color
type trueColor2D struct {
	X, Y       int16
	StatusCode uint8
	B, G, R    uint8
}

func (v trueColor2D) Color() color.Color {
	if v.StatusCode&statusBlanking != 0 {
		return color.Transparent
	}
	return color.RGBA{v.R, v.G, v.B, 0xff}
}
