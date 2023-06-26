package decoder

import (
	"image"
	"os"

	"golang.org/x/image/bmp"
)

type BmpImageDecoder struct{}

func NewBMPImageDecoder() *BmpImageDecoder {
	return &BmpImageDecoder{}
}

func (jpg *BmpImageDecoder) Decode(file *os.File) (image.Image, error) {
	img, err := bmp.Decode(file)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	return rgba, nil
}
