package camera

import "github.com/HellsKitchen99/camera-processor/internal/domain"

type CameraReader struct {
	CamerasUrl []domain.Camera
	Jobs       chan<- domain.FrameJob
}

func NewCameraReader(camerasUrl []domain.Camera, jobs chan<- domain.FrameJob) *CameraReader {
	return &CameraReader{
		CamerasUrl: camerasUrl,
		Jobs:       jobs,
	}
}

func (c *CameraReader) Run() {
	for i := 0; i < len(c.CamerasUrl); i++ {
		go func(cameraId int, cameraUrl string) {
			c.reader(cameraId, cameraUrl)
		}(c.CamerasUrl[i].ID, c.CamerasUrl[i].URL)
	}
}

func (c *CameraReader) reader(cameraId int, cameraUrl string) {
	// считывать камеру и кидать все в Jobs
}
