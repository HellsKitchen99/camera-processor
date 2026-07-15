package detector

import (
	"image"
	"sync"

	"github.com/HellsKitchen99/camera-processor/internal/domain"
	ort "github.com/yalue/onnxruntime_go"
)

const (
	inputWidth         = 640
	inputHeight        = 640
	cowClassID         = 19
	yoloOutputChannels = 84
	yoloOutputBoxes    = 8400
)

type YoloDetector struct {
	mu sync.Mutex

	modelPath string

	session *ort.AdvancedSession

	inputTensor  *ort.Tensor[float32]
	outputTensor *ort.Tensor[float32]

	inputData  []float32
	outputData []float32
}

func NewYoloDetector(modelPath string) (*YoloDetector, error) {
	inputData := make([]float32, 1*3*inputHeight*inputWidth)
	inputTensor, err := ort.NewTensor(ort.NewShape(1, 3, inputHeight, inputWidth), inputData)
	if err != nil {
		return nil, err
	}
	outputTensor, err := ort.NewEmptyTensor[float32](ort.NewShape(1, yoloOutputChannels, yoloOutputBoxes))
	if err != nil {
		return nil, err
	}
	session, err := ort.NewAdvancedSession(modelPath, []string{"images"}, []string{"output0"}, []ort.Value{inputTensor}, []ort.Value{outputTensor}, nil)
	if err != nil {
		_ = inputTensor.Destroy()
		_ = outputTensor.Destroy()
		return nil, err
	}

	return &YoloDetector{
		mu:           sync.Mutex{},
		modelPath:    modelPath,
		session:      session,
		inputTensor:  inputTensor,
		outputTensor: outputTensor,
		inputData:    inputData,
		outputData:   outputTensor.GetData(),
	}, nil
}

func (d *YoloDetector) DetectCow(img image.Image) ([]domain.Detection, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := PreprocessImage(img, d.inputData); err != nil {
		return []domain.Detection{}, err
	}
	if err := d.session.Run(); err != nil {
		return []domain.Detection{}, err
	}
	detections := PostprocessImage(d.outputData, img.Bounds().Dx(), img.Bounds().Dy())
	return detections, nil
}

func (d *YoloDetector) Close() {
	if d.session != nil {
		_ = d.session.Destroy()
		d.session = nil
	}

	if d.inputTensor != nil {
		_ = d.inputTensor.Destroy()
		d.inputTensor = nil
	}

	if d.outputTensor != nil {
		_ = d.outputTensor.Destroy()
		d.outputTensor = nil
	}
}
