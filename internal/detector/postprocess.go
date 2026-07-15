package detector

import (
	"math"
	"sort"

	"github.com/HellsKitchen99/camera-processor/internal/domain"
)

const (
	confidenceThreshold = 0.5
	nmsThreshold        = 0.45
)

func PostprocessImage(output []float32, originalW, originalH int) []domain.Detection {
	detections := make([]domain.Detection, 0)

	scaleX := float32(originalW) / float32(inputWidth)
	scaleY := float32(originalH) / float32(inputHeight)

	for i := 0; i < yoloOutputBoxes; i++ {
		x := output[0*yoloOutputBoxes+i]
		y := output[1*yoloOutputBoxes+i]
		w := output[2*yoloOutputBoxes+i]
		h := output[3*yoloOutputBoxes+i]

		confidence := output[(4+cowClassID)*yoloOutputBoxes+i]
		if confidence < confidenceThreshold {
			continue
		}

		x1 := (x - w/2) * scaleX
		y1 := (y - h/2) * scaleY
		x2 := (x + w/2) * scaleX
		y2 := (y + h/2) * scaleY

		x1 = clamp(x1, 0, float32(originalW))
		y1 = clamp(y1, 0, float32(originalH))
		x2 = clamp(x2, 0, float32(originalW))
		y2 = clamp(y2, 0, float32(originalH))

		detections = append(detections, domain.Detection{
			ClassID:    cowClassID,
			ClassName:  "cow",
			Confidence: confidence,
			X1:         x1,
			Y1:         y1,
			X2:         x2,
			Y2:         y2,
		})
	}

	return nms(detections, nmsThreshold)
}

// help funcs
func nms(detections []domain.Detection, threshold float32) []domain.Detection {
	if len(detections) == 0 {
		return nil
	}

	sort.Slice(detections, func(i, j int) bool {
		return detections[i].Confidence > detections[j].Confidence
	})

	result := make([]domain.Detection, 0, len(detections))
	removed := make([]bool, len(detections))

	for i := 0; i < len(detections); i++ {
		if removed[i] {
			continue
		}

		current := detections[i]
		result = append(result, current)

		for j := i + 1; j < len(detections); j++ {
			if removed[j] {
				continue
			}

			if iou(current, detections[j]) > threshold {
				removed[j] = true
			}
		}
	}

	return result
}

func iou(a, b domain.Detection) float32 {
	x1 := max(a.X1, b.X1)
	y1 := max(a.Y1, b.Y1)
	x2 := min(a.X2, b.X2)
	y2 := min(a.Y2, b.Y2)

	intersectionW := max(0, x2-x1)
	intersectionH := max(0, y2-y1)
	intersectionArea := intersectionW * intersectionH

	areaA := max(0, a.X2-a.X1) * max(0, a.Y2-a.Y1)
	areaB := max(0, b.X2-b.X1) * max(0, b.Y2-b.Y1)

	unionArea := areaA + areaB - intersectionArea
	if unionArea <= 0 {
		return 0
	}

	return intersectionArea / unionArea
}

func clamp(value, minValue, maxValue float32) float32 {
	return float32(math.Max(float64(minValue), math.Min(float64(value), float64(maxValue))))
}

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
