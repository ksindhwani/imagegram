package service

import (
	"fmt"

	"github.com/ksindhwani/imagegram/pkg/config"
	"github.com/ksindhwani/imagegram/pkg/database"
	"github.com/ksindhwani/imagegram/pkg/internal/converter"
)

const (
	LENGTH600 = 600
	WIDTH600  = 600
)

type ImageConvertorService struct {
	Config   *config.Config
	Database database.Database
}

func NewImageConvertorService(
	Config *config.Config,
	database database.Database,
) *ImageConvertorService {
	return &ImageConvertorService{
		Config:   Config,
		Database: database,
	}
}

func (ics *ImageConvertorService) ConvertImages() (
	[]converter.ImageConversionResponse, []converter.ImageConversionResponse, error) {
	response, err := ics.Database.GetAllImages()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to fetch images from database - %w", err)
	}
	return converter.ConvertImagesIntoJpgAndSize(response, ics.Config, LENGTH600, WIDTH600)
}
