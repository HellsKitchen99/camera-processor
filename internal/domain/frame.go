package domain

import (
	"image"
	"time"
)

type Frame struct {
}

type FrameJob struct {
	CameraID   int
	Image      image.Image
	ReceivedAt time.Time
}
