package detector

import (
	"fmt"
	"image"
)

func PreprocessImage(img image.Image, result []float32) error {
	needBuff := 3 * inputWidth * inputHeight
	if len(result) < needBuff {
		return fmt.Errorf("preprocess buffer too small: got %v, need %v", len(result), needBuff)
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	for y := 0; y < inputHeight; y++ {
		srcY := bounds.Min.Y + y*height/inputHeight

		for x := 0; x < inputWidth; x++ {
			srcX := bounds.Min.X + x*width/inputWidth

			r, g, b, _ := img.At(srcX, srcY).RGBA()

			i := y*inputWidth + x

			result[i] = float32(r>>8) / 255.0
			result[inputWidth*inputHeight+i] = float32(g>>8) / 255.0
			result[2*inputWidth*inputHeight+i] = float32(b>>8) / 255.0
		}
	}
	return nil
}
