package service

import (
	"log"

	"github.com/ksindhwani/imagegram/pkg/config"
	"github.com/ksindhwani/imagegram/pkg/database"
	"github.com/ksindhwani/imagegram/pkg/filesystem"
	"github.com/ksindhwani/imagegram/pkg/internal/converter"
)

type ImageService struct {
	Config     config.Config
	Database   database.Database
	FileSystem filesystem.FileSystem
}

func NewImageService(
	Config *config.Config,
	database database.Database,
) *ImageService {
	return &ImageService{
		Config:   *Config,
		Database: database,
	}
}

func (is *ImageService) UpdateConvertedLocationsForImages(convertedImages []converter.ImageConversionResponse) error {
	var err error
	for _, image := range convertedImages {
		err = is.Database.UpdateImageConvertedData(image)
		if err != nil {
			log.Printf("unable to save image in database: %s", err.Error())
			continue
		}
	}
	return err
}
