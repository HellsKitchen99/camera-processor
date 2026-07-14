package config

import (
	"fmt"
	"os"
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
