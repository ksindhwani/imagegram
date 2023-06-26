package decoder

import (
	"image"
	"os"
)

var typeDecoderMap = map[string]ImageDecoder{
	".jpg":  NewJpgImageDecoder(),
	".jpeg": NewJpgImageDecoder(),
	".png":  NewPngImageDecoder(),
	".bmp":  NewBMPImageDecoder(),
}

type ImageDecoder interface {
	Decode(file *os.File) (image.Image, error)
}

func New(imageExtension string) ImageDecoder {
	return typeDecoderMap[imageExtension]
}
