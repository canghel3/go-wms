package render

import "image"

type Renderer interface {
	Render() (image.Image, error)
}
