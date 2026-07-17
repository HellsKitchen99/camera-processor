package workerpool

import "image"

type ImageWriter interface {
	SaveImage(img image.Image, cameraId int) error
}
