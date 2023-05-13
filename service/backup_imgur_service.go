package service

import (
	"disfactory/imgur-backup/domain"
	log "github.com/sirupsen/logrus"
)

type BackupImgurService struct {
	factoryImageRepository domain.FactoryImageRepository
	backupRepository       domain.BackupRepository
	imageDownloader        domain.ImageDownloader

	WorkerNum int
}

func NewBackupImgurService(
	factoryImageRepository domain.FactoryImageRepository,
	backupRepository domain.BackupRepository,
	imageDownloader domain.ImageDownloader,
) *BackupImgurService {
	return &BackupImgurService{
		factoryImageRepository: factoryImageRepository,
		backupRepository:       backupRepository,
		imageDownloader:        imageDownloader,
	}
}

func (service *BackupImgurService) Backup() error {
	// Backup images if not exist in backup repository
	for {
		images, err := service.factoryImageRepository.GetImages(100, 0)
		if err != nil {
			return err
		}

		// If there is no image, stop the backup
		if len(images) == 0 {
			return nil
		}

		for _, image := range images {
			log.WithFields(log.Fields{
				"image_path": image.ImagePath,
			}).Info("Backup image")

			if service.backupRepository.IsExist(image.GetFileName()) {
				continue
			}

			// Download image
			image_data, err := service.imageDownloader.Download(image.ImagePath)
			if err != nil {
				log.WithFields(log.Fields{
					"image_path": image.ImagePath,
				}).Error("Failed to download image")
				continue
			}

			// Backup
			err = service.backupRepository.BackupImage(image.GetFileName(), image_data)
			if err != nil {
				log.WithFields(log.Fields{
					"image_path": image.ImagePath,
				}).Error("Failed to backup image")
				continue
			}
		}
	}
}
