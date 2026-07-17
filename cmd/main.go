package main

import (
	"github.com/HellsKitchen99/camera-processor/internal/build"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := build.Build(); err != nil {
		logrus.Error(err)
		return
	}
}
