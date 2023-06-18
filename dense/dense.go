package dense

import (
	"github.com/r2dtools/mnist/loader"
	"gonum.org/v1/gonum/mat"
)

// rows - rows count in image matrix
// columns - columns count in image matrix
func NewImageDenses(images []loader.Image, rows, columns int) []*mat.Dense {
	return NewNormalizedImageDenses(images, rows, columns, 1)
}

// rows - rows count in image matrix
// columns - columns count in image matrix
// n - devide all image bytes by n
func NewNormalizedImageDenses(images []loader.Image, rows, columns, n int) []*mat.Dense {
	denses := make([]*mat.Dense, len(images))
	var denseData []float64

	for i, image := range images {
		denseData = convertImageToFloat64(image, n)
		denses[i] = mat.NewDense(rows, columns, denseData)
	}

	return denses
}

func NewLabelVecDense(labels []loader.Label) *mat.VecDense {
	data := make([]float64, len(labels))
	var (
		label loader.Label
		i     int
	)

	for i, label = range labels {
		data[i] = float64(label)
	}

	return mat.NewVecDense(len(labels), data)
}

func convertImageToFloat64(image loader.Image, n int) []float64 {
	data := make([]float64, len(image))
	c := float64(n)
	var (
		b byte
		i int
	)

	for i, b = range image {
		data[i] = float64(b) / c
	}

	return data
}
