package compression

import (
	"image"
	"image/color"
	"sync"

	"gonum.org/v1/gonum/mat"
)

type SVDCompressor struct{}

func NewSVDCompressor() *SVDCompressor {
	return &SVDCompressor{}
}

func (s *SVDCompressor) Compress(img image.Image, numComponents int) (*CompressionResult, error) {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()

	grayMatrix := imageToMatrix(img)

	var svd mat.SVD
	if ok := svd.Factorize(grayMatrix, mat.SVDFull); !ok {
		return nil, ErrSVDFailed
	}

	var u, v mat.Dense
	values := svd.Values(nil)
	svd.UTo(&u)
	svd.VTo(&v)

	if numComponents > len(values) {
		numComponents = len(values)
	}

	samplePoints := GenerateSamplePoints(numComponents)
	componentLevels := make([]ComponentLevel, len(samplePoints))

	var wg sync.WaitGroup
	for idx, k := range samplePoints {
		wg.Add(1)
		go func(idx, k int) {
			defer wg.Done()
			reconstructed := svdReconstructWithComponents(&u, &v, values, k, width, height)
			componentLevels[idx] = ComponentLevel{
				NumComponents: k,
				DataSize:      k * (height + width + 1) * 8,
				Image:         reconstructed,
			}
		}(idx, k)
	}
	wg.Wait()

	return &CompressionResult{
		Method:          "SVD",
		OriginalSize:    width * height * 3,
		ComponentLevels: componentLevels,
	}, nil
}

func svdReconstructWithComponents(u, v *mat.Dense, values []float64, k, width, height int) image.Image {
	ur, _ := u.Dims()
	vr, _ := v.Dims()

	uTrunc, vTrunc := mat.NewDense(ur, k, nil), mat.NewDense(vr, k, nil)

	for i := 0; i < ur; i++ {
		for j := 0; j < k; j++ {
			uTrunc.Set(i, j, u.At(i, j))
		}
	}

	for i := 0; i < vr; i++ {
		for j := 0; j < k; j++ {
			vTrunc.Set(i, j, v.At(i, j))
		}
	}

	var temp mat.Dense
	temp.Mul(uTrunc, mat.NewDiagDense(k, values[:k]))

	var result mat.Dense
	result.Mul(&temp, vTrunc.T())

	img := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetGray(x, y, color.Gray{Y: uint8(Clamp(result.At(y, x), 0, 255))})
		}
	}

	return img
}

func imageToMatrix(img image.Image) *mat.Dense {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	data := make([]float64, height*width)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x+bounds.Min.X, y+bounds.Min.Y).RGBA()
			data[y*width+x] = 0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)
		}
	}

	return mat.NewDense(height, width, data)
}
