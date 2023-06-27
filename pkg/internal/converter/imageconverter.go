package converter

import (
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ksindhwani/imagegram/pkg/config"
	"github.com/ksindhwani/imagegram/pkg/internal/decoder"
	"github.com/ksindhwani/imagegram/pkg/internal/tables"
	"github.com/nfnt/resize"
)

const CONVERTED_IMAGE_SUBDIRECTORY = "converted"

type ImageConversionResponse struct {
	ImageId                int64
	ImageName              string
	ImageLocation          string
	ConvertedImageName     string
	ConvertedImageLocation string
	ConversionStatus       bool
	Error                  error
}

func ConvertImagesIntoJpgAndSize(
	images []tables.ImageTable,
	config *config.Config,
	length int,
	width int,
) ([]ImageConversionResponse, []ImageConversionResponse, error) {

	destinationDirectory, err := createDestinationDirectory(config.HostImageDirectory)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create subdirectory for images - %w", err)
	}

	return convertImages(images, config.HostImageDirectory, destinationDirectory, length, width)
}

func convertImages(
	files []tables.ImageTable,
	sourceDir string,
	destinationDir string,
	length int,
	width int,
) ([]ImageConversionResponse, []ImageConversionResponse, error) {

	var successfullConversions []ImageConversionResponse
	var failedConversions []ImageConversionResponse

	// Process each file in the source directory
	for _, file := range files {
		imageFilePath := filepath.Join(sourceDir, file.ImageFileName)

		// Check if the file is an image
		if isImage(imageFilePath) {

			// Open the image file
			imageFile, err := os.Open(imageFilePath)
			if err != nil {
				failedConversions = addToFailedConversion(failedConversions, file, fmt.Errorf("error opening image: %w", err))
				continue
			}
			defer imageFile.Close()

			// Decode the image
			imageExtension := strings.ToLower(filepath.Ext(imageFilePath))
			fileWithoutExt := file.ImageFileName[:len(file.ImageFileName)-len(imageExtension)]
			img, err := decoder.New(imageExtension).Decode(imageFile)
			if err != nil {
				failedConversions = addToFailedConversion(failedConversions, file, fmt.Errorf("error decoding image: %w", err))
				continue
			}

			// Resize the image to length * width pixels
			resizedImg := resize.Resize(uint(length), uint(width), img, resize.Lanczos3)
			convertedFileName := strconv.FormatInt(file.ImageId, 10) + "converted" + fileWithoutExt + ".jpg"

			// Create the destination file path
			destinationFilePath := filepath.Join(destinationDir, convertedFileName)

			// Create the destination file
			destinationFile, err := createConvertedFile(destinationFilePath)
			if err != nil {
				failedConversions = addToFailedConversion(failedConversions, file, fmt.Errorf("error creating destination file: %w", err))
				continue
			}
			defer destinationFile.Close()

			// Encode the resized image as JPEG and save it to the destination file
			err = jpeg.Encode(destinationFile, resizedImg, nil)
			if err != nil {
				failedConversions = addToFailedConversion(failedConversions, file, fmt.Errorf("error encoding image: %w", err))
				continue
			}

			successfullConversions = addToSuccessfulConversion(successfullConversions, file, convertedFileName, destinationFilePath)
		}
	}

	return successfullConversions, failedConversions, nil
}

func createConvertedFile(filePath string) (*os.File, error) {
	// Check if the file already exists
	_, err := os.Stat(filePath)
	if err == nil {
		// File exists, remove it
		err := os.Remove(filePath)
		if err != nil {
			return nil, fmt.Errorf("error removing existing file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		// Error occurred while checking file existence
		return nil, fmt.Errorf("error checking file existence: %w", err)
	}
	// Create the new file
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("error creating file: %w", err)
	}
	return file, nil
}

func addToSuccessfulConversion(
	successfullConversions []ImageConversionResponse,
	file tables.ImageTable,
	convertedFileName,
	destinationFilePath string,
) []ImageConversionResponse {
	response := ImageConversionResponse{
		ImageId:                file.ImageId,
		ImageName:              file.ImageFileName,
		ImageLocation:          file.Location,
		ConvertedImageName:     convertedFileName,
		ConvertedImageLocation: destinationFilePath,
		ConversionStatus:       true,
		Error:                  nil,
	}
	successfullConversions = append(successfullConversions, response)
	return successfullConversions
}

func addToFailedConversion(failedConversions []ImageConversionResponse, file tables.ImageTable, err error) []ImageConversionResponse {
	response := ImageConversionResponse{
		ImageId:          file.ImageId,
		ImageName:        file.ImageFileName,
		ImageLocation:    file.Location,
		ConversionStatus: false,
		Error:            err,
	}
	failedConversions = append(failedConversions, response)
	return failedConversions
}

func createDestinationDirectory(hostDirectory string) (string, error) {
	// Create the full path to the subdirectory
	subDirPath := filepath.Join(hostDirectory, CONVERTED_IMAGE_SUBDIRECTORY)

	// Check if the subdirectory already exists
	if _, err := os.Stat(subDirPath); os.IsNotExist(err) {
		// Create the subdirectory and its parent directories
		err := os.MkdirAll(subDirPath, 0755) // 0755 sets the directory's permission
		if err != nil {
			return "", err
		}
	}
	return subDirPath, nil
}

func (icr ImageConversionResponse) ToString() string {
	return fmt.Sprintf("%s: %s", icr.ImageLocation, icr.Error.Error())
}

// Helper function to check if a file has an image extension
func isImage(filename string) bool {
	extensions := []string{".jpg", ".jpeg", ".png", ".bmp", ".gif"}
	ext := strings.ToLower(filepath.Ext(filename))
	for _, e := range extensions {
		if ext == e {
			return true
		}
	}
	return false
}
