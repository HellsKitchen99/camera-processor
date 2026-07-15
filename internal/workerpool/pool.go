package workerpool

import (
	"fmt"

	"github.com/HellsKitchen99/camera-processor/internal/domain"
	"github.com/sirupsen/logrus"
)

type WorkerPool struct {
	yoloDetector YoloDetector
	annotator    Annotator

	Workers int
	Jobs    <-chan domain.FrameJob
}

func NewWorkerPool(yoloDetector YoloDetector, annoattor Annotator, workers int, jobs <-chan domain.FrameJob) *WorkerPool {
	return &WorkerPool{
		yoloDetector: yoloDetector,
		annotator:    annoattor,

		Workers: workers,
		Jobs:    jobs,
	}
}

func (w *WorkerPool) StartWorkers() {
	for i := 0; i < w.Workers; i++ {
		go func(workerId int) {
			w.worker(workerId)
		}(i)
	}
}

func (w *WorkerPool) worker(workerId int) {
	for job := range w.Jobs {
		image := job.Image
		detections, err := w.yoloDetector.DetectCow(image)
		if err != nil {
			logrus.Errorf("CAMERA %v WORKER %v: %v", job.CameraID, workerId, err)
			continue
		}
		if len(detections) == 0 {
			continue
		}
		for _, detection := range detections {
			fmt.Printf(
				"worker=%d camera=%d cow confidence=%.2f box=(%.0f %.0f %.0f %.0f)\n",
				workerId,
				job.CameraID,
				detection.Confidence,
				detection.X1,
				detection.Y1,
				detection.X2,
				detection.Y2,
			)
		}
		finalImage := w.annotator.VisualizeDetections(job.Image, detections)
		fmt.Println(finalImage)
	}
}
