package imagewriter

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"time"
)

type ImageWriter struct {
	imagePath string
}

func NewImageWriter(imagePath string) *ImageWriter {
	return &ImageWriter{
		imagePath: imagePath,
	}
}

func (i *ImageWriter) SaveImage(img image.Image, cameraId int) error {
	cameraDir := filepath.Join(i.imagePath, fmt.Sprintf("camera_%v", cameraId))
	if err := os.MkdirAll(cameraDir, 0755); err != nil {
		return err
	}
	fileName := fmt.Sprintf("detection_%v.jpg", time.Now().Format("20060102_150405_000"))
	fullPath := filepath.Join(cameraDir, fileName)
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := jpeg.Encode(file, img, &jpeg.Options{Quality: 90}); err != nil {
		return err
	}
	return nil
}
