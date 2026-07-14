package workerpool

import "github.com/HellsKitchen99/camera-processor/internal/domain"

type WorkerPool struct {
	Workers int
	Jobs    <-chan domain.FrameJob
}

func NewWorkerPool(workers int, jobs <-chan domain.FrameJob) *WorkerPool {
	return &WorkerPool{
		Workers: workers,
		Jobs:    jobs,
	}
}

func (w *WorkerPool) StartWorkers() {
	for i := 0; i < w.Workers; i++ {
		go w.work()
	}
}

func (w *WorkerPool) work() {
	// брать кадры из Jobs
}
