package filesystem

import (
	"mime/multipart"

	"github.com/ksindhwani/imagegram/pkg/config"
	"github.com/ksindhwani/imagegram/pkg/filesystem/local"
)

const (
	LOCAL = "local"
	S3    = "s3"
)

type FileSystem interface {
	SaveFile(fileName string, file multipart.File) (string, error)
}

func New(fileSystemType string, config *config.Config) (FileSystem, error) {
	switch fileSystemType {
	case LOCAL:
		return &local.LocalFileSystem{
			HostDirectory:  config.HostImageDirectory,
			LocalDirectory: config.LocalImageDirectory,
		}, nil
	default:
		return &local.LocalFileSystem{
			HostDirectory:  config.HostImageDirectory,
			LocalDirectory: config.LocalImageDirectory,
		}, nil
	}
}
