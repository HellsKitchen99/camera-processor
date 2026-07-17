package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/HellsKitchen99/camera-processor/internal/domain"
)

func LoadCamerasUrl() ([]domain.Camera, error) {
	camerasUrl := os.Getenv("CAMERAS_URL")
	if camerasUrl == "" {
		return []domain.Camera{}, fmt.Errorf("there is no any camera url")
	}
	camerasUrlSep := strings.Split(camerasUrl, ",")
	cameras := make([]domain.Camera, len(camerasUrlSep))
	for id, url := range camerasUrlSep {
		url = strings.TrimSpace(url)
		if url == "" {
			continue
		}
		cameras[id] = domain.Camera{
			ID:  id,
			URL: url,
		}
	}
	return cameras, nil
}

func LoadModelPath() string {
	modelPath := os.Getenv("MODEL_PATH")
	return modelPath
}

func LoadImagePath() string {
	imagePath := os.Getenv("IMAGE_PATH")
	return imagePath
}

func LoadWorkersAmount() int {
	workersAmount, err := strconv.Atoi(os.Getenv("WORKERS_AMOUNT"))
	if err != nil {
		workersAmount = 1
	}
	return workersAmount
}

func LoadJobsQueueSize() int {
	jobsQueueSize, err := strconv.Atoi(os.Getenv("JOBS_QUEUE_SIZE"))
	if err != nil {
		jobsQueueSize = 100
	}
	return jobsQueueSize
}

func LoadLibPath() string {
	libPath := os.Getenv("ONNXRUNTIME_LIB_PATH")
	return libPath
}
