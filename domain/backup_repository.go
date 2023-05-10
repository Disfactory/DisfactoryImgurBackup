package domain

type BackupRepository interface {
	IsExist(image_name string) bool
	BackupImage(image_name string, image_data []byte) error
}
