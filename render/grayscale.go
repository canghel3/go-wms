package render

import (
	"image"
	"image/color"
)

type GrayScale struct {
	width  int
	height int
	data   []float64
	bbox   [4]float64
}

func NewGrayScale(data []float64, width, height int, bbox [4]float64) *GrayScale {
	return &GrayScale{
		data:   data,
		bbox:   bbox,
		width:  width,
		height: height,
	}
}

func (gs *GrayScale) Render() image.Image {
	img := image.NewGray(image.Rect(0, 0, gs.width, gs.height))

	// Convert the TIFF data to a grayscale image
	for y := 0; y < gs.height; y++ {
		for x := 0; x < gs.width; x++ {
			value := gs.data[y*gs.width+x]
			grayValue := uint8(value * 255)
			img.SetGray(x, y, color.Gray{Y: grayValue})
		}
	}

	return img
}
