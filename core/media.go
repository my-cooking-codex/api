package core

import (
	"github.com/h2non/bimg"
)

// Return the minimum of two integers
func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Optimises an image to JPEG format,
// ensuring that the image is no larger than maxSize
func OptimiseImageToJPEG(image []byte, maxSize int) ([]byte, error) {
	loadedImage := bimg.NewImage(image)
	size, err := loadedImage.Size()
	if err != nil {
		return nil, err
	}
	optimisedImage, err := loadedImage.Process(bimg.Options{
		Width:         intMin(size.Width, maxSize),
		Type:          bimg.JPEG,
		StripMetadata: true,
		Quality:       80,
	})

	return optimisedImage, err
}
