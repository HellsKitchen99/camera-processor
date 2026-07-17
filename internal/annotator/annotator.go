package annotator

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/HellsKitchen99/camera-processor/internal/domain"
)

type Annotator struct {
}

func NewAnnotator() *Annotator {
	return &Annotator{}
}

func (a *Annotator) VisualizeDetections(img image.Image, detections []domain.Detection) *image.RGBA {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	draw.Draw(result, bounds, img, bounds.Min, draw.Src)

	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	for _, det := range detections {
		drawRect(result, int(det.X1), int(det.Y1), int(det.X2), int(det.Y2), red)
	}

	return result
}

func drawRect(img *image.RGBA, x1, y1, x2, y2 int, c color.Color) {
	// защита от выхода за границы
	b := img.Bounds()

	if x1 < b.Min.X {
		x1 = b.Min.X
	}
	if y1 < b.Min.Y {
		y1 = b.Min.Y
	}
	if x2 > b.Max.X-1 {
		x2 = b.Max.X - 1
	}
	if y2 > b.Max.Y-1 {
		y2 = b.Max.Y - 1
	}

	// толщина рамки
	thickness := 3

	// верхняя и нижняя линии
	for t := 0; t < thickness; t++ {
		for x := x1; x <= x2; x++ {
			img.Set(x, y1+t, c)
			img.Set(x, y2-t, c)
		}
	}

	// левая и правая линии
	for t := 0; t < thickness; t++ {
		for y := y1; y <= y2; y++ {
			img.Set(x1+t, y, c)
			img.Set(x2-t, y, c)
		}
	}
}
