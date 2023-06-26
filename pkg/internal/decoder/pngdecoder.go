package decoder

import (
	"image"
	"image/png"
	"os"
)

type PngImageDecoder struct{}

func NewPngImageDecoder() *PngImageDecoder {
	return &PngImageDecoder{}
}

func (jpg *PngImageDecoder) Decode(file *os.File) (image.Image, error) {
	return png.Decode(file)
}
