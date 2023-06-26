package tables

import "time"

type ImageTable struct {
	ImageId                int64
	PostId                 int64
	ImageFileName          string
	Location               string
	ConvertedImageName     string
	ConvertedImageLocation string
	UploadedAt             time.Time
}
