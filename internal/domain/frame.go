package domain

import (
	"image"
	"time"
)

type Frame struct {
}

type FrameJob struct {
	CameraId   int
	Image      image.Image
	RecievedAt time.Time
}
