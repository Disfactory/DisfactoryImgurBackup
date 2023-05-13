package service

import (
	"disfactory/imgur-backup/domain"
	"disfactory/imgur-backup/utils"
	"time"

	"github.com/go-co-op/gocron"
)

type BackupImgurService struct {
	factoryImageRepository domain.FactoryImageRepository
	backupRepository       domain.BackupRepository
	imageDownloader        domain.ImageDownloader
	schedule               string
	scheduler              *gocron.Scheduler
}

func NewBackupImgurService(
	factoryImageRepository domain.FactoryImageRepository,
	backupRepository domain.BackupRepository,
	imageDownloader domain.ImageDownloader,
	schedule string,
) *BackupImgurService {
	return &BackupImgurService{
		factoryImageRepository: factoryImageRepository,
		backupRepository:       backupRepository,
		imageDownloader:        imageDownloader,
		schedule:               schedule,
	}
}

func (service *BackupImgurService) Start() {
	service.scheduler = gocron.NewScheduler(time.UTC)
	utils.Logger().Infof("Start backup cronjob: %s", service.schedule)
	_, err := service.scheduler.Cron(service.schedule).Do(service.Backup)
	if err != nil {
		utils.Logger().Panic("Failed to set cronjob", err)
	}
	service.scheduler.StartAsync()
}

func (service *BackupImgurService) Backup() error {
	logger := utils.Logger()
	logger.Info("Start backup")
	defer logger.Info("Finish backup")
	// Backup images if not exist in backup repository
	for i := 0; i < 65536; i++ {
		images, err := service.factoryImageRepository.GetImages(100, i)
		if err != nil {
			logger.Errorf("Failed to get images", err)
			return err
		}

		// If there is no image, stop the backup
		if len(images) == 0 {
			return nil
		}

		for _, image := range images {
			if service.backupRepository.IsExist(image.GetFileName()) {
				continue
			}

			// Download image
			logger.WithFields(utils.LogFields{
				"image_path": image.ImagePath,
				"created_at": image.CreatedAt,
			}).Info("Download image")
			image_data, err := service.imageDownloader.Download(image.ImagePath)
			if err != nil {
				utils.Logger().WithFields(utils.LogFields{
					"image_path": image.ImagePath,
				}).Error("Failed to download image", err)
				continue
			}

			// Backup
			err = service.backupRepository.BackupImage(image.GetFileName(), image_data)
			if err != nil {
				utils.Logger().WithFields(utils.LogFields{
					"image_path": image.ImagePath,
				}).Error("Failed to backup image", err)
				continue
			}
		}
	}

	return nil
}
