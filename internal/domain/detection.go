package domain

type Detection struct {
	ClassID    int
	ClassName  string
	Confidence float32

	X1 float32
	Y1 float32
	X2 float32
	Y2 float32
}
