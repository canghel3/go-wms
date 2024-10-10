package render

import (
	"image"
	"image/color"
)

type GrayScale struct {
	width  int
	height int
	data   []float64
	invert bool
}

func Gray(data []float64, width, height int, invert bool) *GrayScale {
	return &GrayScale{
		data:   data,
		width:  width,
		height: height,
		invert: invert,
	}
}

func (gs *GrayScale) Render() (image.Image, error) {
	img := image.NewGray(image.Rect(0, 0, gs.width, gs.height))

	min, max := findMinMax(gs.data)
	// Convert the TIFF data to a grayscale image
	for y := 0; y < gs.height; y++ {
		for x := 0; x < gs.width; x++ {
			value := gs.data[y*gs.width+x]
			normalized := (value - min) / (max - min)
			var grayValue uint8
			if gs.invert {
				grayValue = uint8((1 - normalized) * 255)
			} else {
				grayValue = uint8(normalized * 255)
			}

			img.SetGray(x, y, color.Gray{Y: grayValue})
		}
	}

	return img, nil
}

func findMinMax(data []float64) (min, max float64) {
	min, max = data[0], data[0]
	for _, v := range data {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return
}
