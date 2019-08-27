package ilda

import "image/color"

// DefaultPalette used by most ILDA files that do not contain a color palette
var DefaultPalette = color.Palette{
	color.RGBA{255, 0, 0, 255}, // Red
	color.RGBA{255, 16, 0, 255},
	color.RGBA{255, 32, 0, 255},
	color.RGBA{255, 48, 0, 255},
	color.RGBA{255, 64, 0, 255},
	color.RGBA{255, 80, 0, 255},
	color.RGBA{255, 96, 0, 255},
	color.RGBA{255, 112, 0, 255},
	color.RGBA{255, 128, 0, 255},
	color.RGBA{255, 144, 0, 255},
	color.RGBA{255, 160, 0, 255},
	color.RGBA{255, 176, 0, 255},
	color.RGBA{255, 176, 0, 255},
	color.RGBA{255, 192, 0, 255},
	color.RGBA{255, 208, 0, 255},
	color.RGBA{255, 224, 0, 255},
	color.RGBA{255, 240, 0, 255},
	color.RGBA{255, 255, 0, 255}, // Yellow
	color.RGBA{192, 255, 0, 255},
	color.RGBA{160, 255, 0, 255},
	color.RGBA{128, 255, 0, 255},
	color.RGBA{96, 255, 0, 255},
	color.RGBA{64, 255, 0, 255},
	color.RGBA{32, 255, 0, 255},
	color.RGBA{0, 255, 0, 255}, // Green
	color.RGBA{0, 255, 36, 255},
	color.RGBA{0, 255, 73, 255},
	color.RGBA{0, 255, 109, 255},
	color.RGBA{0, 255, 146, 255},
	color.RGBA{0, 255, 182, 255},
	color.RGBA{0, 255, 219, 255},
	color.RGBA{0, 255, 255, 255}, // Cyan
	color.RGBA{0, 227, 255, 255},
	color.RGBA{0, 198, 255, 255},
	color.RGBA{0, 170, 255, 255},
	color.RGBA{0, 142, 255, 255},
	color.RGBA{0, 113, 255, 255},
	color.RGBA{0, 85, 255, 255},
	color.RGBA{0, 56, 255, 255},
	color.RGBA{0, 28, 255, 255},
	color.RGBA{0, 0, 255, 255}, // Blue
	color.RGBA{32, 0, 255, 255},
	color.RGBA{64, 0, 255, 255},
	color.RGBA{96, 0, 255, 255},
	color.RGBA{128, 0, 255, 255},
	color.RGBA{160, 0, 255, 255},
	color.RGBA{192, 0, 255, 255},
	color.RGBA{224, 0, 255, 255},
	color.RGBA{255, 0, 255, 255}, // Magenta
	color.RGBA{255, 32, 255, 255},
	color.RGBA{255, 64, 255, 255},
	color.RGBA{255, 96, 255, 255},
	color.RGBA{255, 128, 255, 255},
	color.RGBA{255, 160, 255, 255},
	color.RGBA{255, 192, 255, 255},
	color.RGBA{255, 224, 255, 255},
	color.RGBA{255, 255, 255, 255}, // White
	color.RGBA{255, 224, 224, 255},
	color.RGBA{255, 192, 192, 255},
	color.RGBA{255, 160, 160, 255},
	color.RGBA{255, 128, 128, 255},
	color.RGBA{255, 96, 96, 255},
	color.RGBA{255, 64, 64, 255},
	color.RGBA{255, 32, 32, 255},
}
