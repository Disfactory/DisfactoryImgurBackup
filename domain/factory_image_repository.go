package domain

import (
	"github.com/google/uuid"
	"path/filepath"
)

type FactoryImage struct {
	Id        uuid.UUID `db:"uuid"`
	ImagePath string    `db:"image_path"`
}

func (i FactoryImage) GetFileName() string {
	// Get the file name from the link
	// e.g. https://i.imgur.com/3JQ3Z0Y.jpg
	// 3JQ3Z0Y.jpg
	return filepath.Base(i.ImagePath)
}

type FactoryImageRepository interface {
	GetImages(size int, offset int) ([]FactoryImage, error)
}
