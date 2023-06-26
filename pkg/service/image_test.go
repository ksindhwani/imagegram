package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ksindhwani/imagegram/pkg/config"
	"github.com/ksindhwani/imagegram/pkg/internal/converter"
	"github.com/ksindhwani/imagegram/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUpdateConvertedLocationsForImages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		Name                                  string
		Input                                 []converter.ImageConversionResponse
		ExpectedUpdateImageConvertedDataError error
		ExpectedUpdateImageConvertedDataCalls int
		ExpectedError                         error
	}{
		{
			Name: "Test All Valid",
			Input: []converter.ImageConversionResponse{
				{
					ImageId:                1,
					ImageName:              "Test Image",
					ImageLocation:          "Test location",
					ConvertedImageName:     "Test converted name",
					ConvertedImageLocation: "Test converted image location",
					ConversionStatus:       true,
					Error:                  nil,
				},
				{
					ImageId:                2,
					ImageName:              "Test Imagem 2",
					ImageLocation:          "Test location 2",
					ConvertedImageName:     "Test converted name 2",
					ConvertedImageLocation: "Test converted image location 2",
					ConversionStatus:       true,
					Error:                  nil,
				},
			},
			ExpectedUpdateImageConvertedDataError: nil,
			ExpectedUpdateImageConvertedDataCalls: 2,
			ExpectedError:                         nil,
		},
		{
			Name: "Test when image gets update query error",
			Input: []converter.ImageConversionResponse{
				{
					ImageId:                1,
					ImageName:              "Test Image",
					ImageLocation:          "Test location",
					ConvertedImageName:     "Test converted name",
					ConvertedImageLocation: "Test converted image location",
					ConversionStatus:       true,
					Error:                  nil,
				},
				{
					ImageId:                2,
					ImageName:              "Test Imagem 2",
					ImageLocation:          "Test location 2",
					ConvertedImageName:     "Test converted name 2",
					ConvertedImageLocation: "Test converted image location 2",
					ConversionStatus:       true,
					Error:                  nil,
				},
			},
			ExpectedUpdateImageConvertedDataError: errors.New("error in db query"),
			ExpectedUpdateImageConvertedDataCalls: 2,
			ExpectedError:                         errors.New("error in db query"),
		},
	}

	any := gomock.Any()
	config := config.Config{
		HostImageDirectory:  "test host directory",
		LocalImageDirectory: "test local directory",
	}
	database := mocks.NewMockDatabase(ctrl)
	imageService := NewImageService(&config, database)
	for _, test := range tests {
		database.EXPECT().UpdateImageConvertedData(any).
			Return(test.ExpectedUpdateImageConvertedDataError).
			Times(test.ExpectedUpdateImageConvertedDataCalls)
		err := imageService.UpdateConvertedLocationsForImages(test.Input)
		assert.Equal(t, test.ExpectedError, err, test.Name)
	}
}
