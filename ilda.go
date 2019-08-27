// Package ilda implements ILDA Image Data Transfer Format
package ilda

type Header struct {
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

// Format1 – 2D Coordinates with Indexed Color
type Format1 struct {
	X, Y       int16
	StatusCode uint8
	ColorIndex uint8
}

// Format2 – Color Palette
type Format2 struct {
	R, G, B uint8
}

// Format4 – 3D Coordinates with True Color
type Format4 struct {
	X, Y, Z    int16
	StatusCode uint8
	B, G, R    uint8
}

// Format5 – 2D Coordinates with True Color
type Format5 struct {
	X, Y       int16
	StatusCode uint8
	B, G, R    uint8
}
