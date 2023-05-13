package domain

type ImageDownloader interface {
	Download(imagePath string) ([]byte, error)
}
