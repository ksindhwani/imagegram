package local

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

type LocalFileSystem struct {
	HostDirectory  string
	LocalDirectory string
}

// Getting instance of local file system
func New(hostDirectory string, localDirectory string) *LocalFileSystem {
	return &LocalFileSystem{
		HostDirectory:  hostDirectory,
		LocalDirectory: localDirectory,
	}
}

// Saving file in a host directory and returning its location
func (lfs *LocalFileSystem) SaveFile(fileName string, file multipart.File) (string, error) {
	// Create a new file on the host machine to store the uploaded image
	dst, err := os.Create(lfs.LocalDirectory + fileName)
	if err != nil {
		return "", fmt.Errorf("error creating the file in local - %w", err)
	}
	defer dst.Close()

	// Copy the contents of the uploaded file to the destination file
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("error copying the file - %w", err)
	}

	return dst.Name(), nil
}
