package render

import (
	"image"
	"image/color"
)

type RGBRenderer struct {
	width  int
	height int
	data   []float64
}

func RGB(data []float64, width, height int) *RGBRenderer {
	return &RGBRenderer{
		width:  width,
		height: height,
		data:   data,
	}
}

func (r *RGBRenderer) Render() (image.Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, r.width, r.height))
	min, max := findMinMax(r.data)

	// Normalize and apply the color map
	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			value := r.data[y*r.width+x]
			//TODO: extract min and max from whole raster once and use that for normalization,
			// otherwise the raster changes color as we zoom in

			// Normalize the value between 0 and 255
			normalized := uint8((value - min) / (max - min) * 255)
			col := applyColorMap(normalized)
			img.Set(x, y, col)
		}
	}
	return img, nil
}

// Define the color map as per the provided style
var colorMap = []colorEntry{
	{color.RGBA{0x00, 0x00, 0x00, 0xFF}, 0, 10},    // #000000
	{color.RGBA{0x09, 0x3a, 0x7f, 0xFF}, 10, 20},   // #093a7f
	{color.RGBA{0x00, 0x2b, 0x65, 0xFF}, 20, 30},   // #002b65
	{color.RGBA{0x05, 0x3e, 0x88, 0xFF}, 30, 40},   // #053e88
	{color.RGBA{0x11, 0x4e, 0xa3, 0xFF}, 40, 50},   // #114ea3
	{color.RGBA{0x27, 0x5a, 0xb6, 0xFF}, 50, 75},   // #275ab6
	{color.RGBA{0x72, 0x83, 0xad, 0xFF}, 75, 100},  // #7283ad
	{color.RGBA{0x86, 0x8c, 0x9e, 0xFF}, 100, 150}, // #868c9e
	{color.RGBA{0xb1, 0xa9, 0x6a, 0xFF}, 150, 200}, // #b1a96a
	{color.RGBA{0xff, 0xea, 0x49, 0xFF}, 200, 230}, // #ffea49
	{color.RGBA{0xff, 0xea, 0x49, 0xFF}, 230, 255}, // #ffea49 (Extreme)
}

type colorEntry struct {
	color color.RGBA
	min   uint8
	max   uint8
}

func applyColorMap(value uint8) color.RGBA {
	for _, entry := range colorMap {
		if value >= entry.min && value <= entry.max {
			return entry.color
		}
	}
	// Default to black if no range matches
	return color.RGBA{0, 0, 0, 0}
}
