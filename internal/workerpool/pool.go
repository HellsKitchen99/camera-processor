package workerpool

import (
	"fmt"
	"sync"

	"github.com/HellsKitchen99/camera-processor/internal/detector"
	"github.com/HellsKitchen99/camera-processor/internal/domain"
	"github.com/sirupsen/logrus"
)

type WorkerPool struct {
	annotator   Annotator
	imageWriter ImageWriter
	wg          *sync.WaitGroup
	modelPath   string

	Workers int
	Jobs    <-chan domain.FrameJob
}

func NewWorkerPool(modelPath string, annotator Annotator, imageWriter ImageWriter, workers int, jobs <-chan domain.FrameJob, wg *sync.WaitGroup) *WorkerPool {
	return &WorkerPool{
		annotator:   annotator,
		imageWriter: imageWriter,
		wg:          wg,
		modelPath:   modelPath,

		Workers: workers,
		Jobs:    jobs,
	}
}

func (w *WorkerPool) StartWorkers() {
	for i := 0; i < w.Workers; i++ {
		w.wg.Add(1)
		go func(workerId int) {
			defer w.wg.Done()
			w.worker(workerId)
		}(i)
	}
}

func (w *WorkerPool) worker(workerId int) {
	yoloDetector, err := detector.NewYoloDetector(w.modelPath)
	if err != nil {
		logrus.Errorf("WORKER %v: %v\n", workerId, err)
		return
	}
	defer yoloDetector.Close()
	for job := range w.Jobs {
		image := job.Image
		detections, err := yoloDetector.DetectCow(image)
		if err != nil {
			logrus.Errorf("CAMERA %v WORKER %v: %v\n", job.CameraID, workerId, err)
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
		if err := w.imageWriter.SaveImage(finalImage, job.CameraID); err != nil {
			logrus.Errorf("CAMERA %v WORKER %v: %v\n", job.CameraID, workerId, err)
		}
	}
}
