package domain

import (
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type FactoryImage struct {
	Id        uuid.UUID `db:"id"`
	ImagePath string    `db:"image_path"`
	CreatedAt time.Time `db:"created_at"`
}

func (i FactoryImage) GetFileName() string {
	// Get the file name from the link
	// e.g. https://i.imgur.com/3JQ3Z0Y.jpg
	// 3JQ3Z0Y.jpg
	return filepath.Base(i.ImagePath)
}

type FactoryImageRepository interface {
	GetImages(size int, page int) ([]FactoryImage, error)
}
