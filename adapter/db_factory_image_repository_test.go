package adapter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestCase: Connect to 127.0.0.1:5432 and get 10 images from offset 0
// Expected: 10 images are returned
func TestGetImages(t *testing.T) {
	pg_params := PGParameters{
		Account:  "postgres",
		Password: "postgres",
		Host:     "127.0.0.1",
		Port:     "5432",
		DBName:   "disfactory_data",
	}

	repo, err := NewDBFactoryImageRepository(pg_params)
	assert.Nil(t, err)

	images_0_9, err := repo.GetImages(10, 0)
	assert.Nil(t, err)
	assert.Equal(t, 10, len(images_0_9))

	// The order of images_0_9 should be latest to oldest
	for i := 0; i < len(images_0_9)-1; i++ {
		println(images_0_9[i].CreatedAt.String())
		assert.True(t, images_0_9[i].CreatedAt.After(images_0_9[i+1].CreatedAt),
			"images_0_9[%d].CreatedAt %s should be after images_0_9[%d].CreatedAt %s",
			i, images_0_9[i].CreatedAt.String(),
			i+1, images_0_9[i+1].CreatedAt.String())
	}

	images_10_19, err := repo.GetImages(10, 1)
	assert.Nil(t, err)
	assert.Equal(t, 10, len(images_10_19))

	// images_10_19 should not contain images_0_9
	for _, image_0_9 := range images_0_9 {
		for _, image_10_19 := range images_10_19 {
			assert.NotEqual(t, image_0_9.Id, image_10_19.Id)
		}
	}

	assert.Nil(t, repo.Close())
}
