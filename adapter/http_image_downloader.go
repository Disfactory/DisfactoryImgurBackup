package adapter

import (
	"io/ioutil"
	"net/http"
)

type HttpImageDownloader struct {
}

func NewHttpImageDownloader() HttpImageDownloader {
	return HttpImageDownloader{}
}

func (d HttpImageDownloader) Download(image_url string) ([]byte, error) {
	// Download image from image_url
	// e.g. https://i.imgur.com/3JQ3Z0Y.jpg
	// return image_data, err

	http_client := &http.Client{}
	req, err := http.NewRequest("GET", image_url, nil)
	if err != nil {
		return nil, err
	}

	// Set Chrome User-Agent
	req.Header.Set("User-Agent",
		" Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36 Edg/113.0.1774.35")

	// Download image
	resp, err := http_client.Do(req)
	if err != nil {
		return nil, err
	}

	// Read image data
	defer resp.Body.Close()
	image_data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return image_data, nil
}
