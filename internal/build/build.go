package build

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/HellsKitchen99/camera-processor/internal/annotator"
	"github.com/HellsKitchen99/camera-processor/internal/camera"
	"github.com/HellsKitchen99/camera-processor/internal/config"
	"github.com/HellsKitchen99/camera-processor/internal/domain"
	"github.com/HellsKitchen99/camera-processor/internal/imagewriter"
	"github.com/HellsKitchen99/camera-processor/internal/workerpool"
	ort "github.com/yalue/onnxruntime_go"
)

func Build() error {
	cameras, err := config.LoadCamerasUrl()
	if err != nil {
		return err
	}
	modelPath := config.LoadModelPath()
	imagePath := config.LoadImagePath()
	workers := config.LoadWorkersAmount()
	jobsQueueSize := config.LoadJobsQueueSize()
	libPath := config.LoadLibPath()

	ort.SetSharedLibraryPath(libPath)

	if err := ort.InitializeEnvironment(); err != nil {
		return err
	}
	defer ort.DestroyEnvironment()

	var wgWriters sync.WaitGroup
	var wgReaders sync.WaitGroup

	// jobs channel
	jobs := make(chan domain.FrameJob, jobsQueueSize)

	// common context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		cancel()
		wgWriters.Wait()
		close(jobs)
	}()

	// annotator
	annotator := annotator.NewAnnotator()

	// camera reader
	cameraReader := camera.NewCameraReader(cameras, jobs, ctx, &wgWriters)

	// image writer
	imageWriter := imagewriter.NewImageWriter(imagePath)

	// worker pool
	workerPool := workerpool.NewWorkerPool(modelPath, annotator, imageWriter, workers, jobs, &wgReaders)

	cameraReader.Run()
	workerPool.StartWorkers()
	wgReaders.Wait()
	return nil
}
