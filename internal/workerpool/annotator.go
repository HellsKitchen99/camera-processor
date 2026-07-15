package workerpool

import (
	"image"

	"github.com/HellsKitchen99/camera-processor/internal/domain"
)

type Annotator interface {
	VisualizeDetections(img image.Image, detections []domain.Detection) *image.RGBA
}
