package adapter

import (
	"os"
	"path/filepath"
)

type LocalFileBackupRepository struct {
	RootDir string
}

// Create a new LocalFileBackupRepository
func NewLocalFileBackupRepository(root_dir string) LocalFileBackupRepository {
	return LocalFileBackupRepository{RootDir: root_dir}
}

// Generate the image file path
func (r LocalFileBackupRepository) generateImageFilePath(image_name string) string {
	// Generate the image file path
	// e.g. /home/backup/3JQ3Z0Y.jpg
	return filepath.Join(r.RootDir, image_name)
}

// Check if the image file exists
func (r LocalFileBackupRepository) IsExist(image_name string) bool {
	image_path := r.generateImageFilePath(image_name)
	if _, err := os.Stat(image_path); err == nil {
		return true
	}
	return false
}

// Backup the image file
func (r LocalFileBackupRepository) BackupImage(image_name string, image_data []byte) error {
	image_path := r.generateImageFilePath(image_name)

	// Write image_data to file if file not exists create it
	// e.g. /home/backup/3JQ3Z0Y.jpg
	fi, err := os.OpenFile(image_path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer fi.Close()
	fi.Write(image_data)

	return nil
}
