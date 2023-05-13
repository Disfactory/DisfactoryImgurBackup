package adapter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Try to download a image from imgur url
func TestDownloadImage(t *testing.T) {
	url := "https://i.imgur.com/Y3eWZxt.jpg"
	downloader := NewHttpImageDownloader()
	_, err := downloader.Download(url)
	assert.Nil(t, err)
}
