package decoder

import (
	"image"
	"image/jpeg"
	"os"
)

type JpgImageDecoder struct{}

func NewJpgImageDecoder() *JpgImageDecoder {
	return &JpgImageDecoder{}
}

func (jpg *JpgImageDecoder) Decode(file *os.File) (image.Image, error) {
	return jpeg.Decode(file)
}
