package workerpool

import (
	"image"

	"github.com/HellsKitchen99/camera-processor/internal/domain"
)

type YoloDetector interface {
	DetectCow(img image.Image) ([]domain.Detection, error)
}
