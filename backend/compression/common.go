package compression

import (
	"errors"
	"image"
	"image/color"
)

var (
	ErrInvalidImage = errors.New("invalid image")
	ErrSVDFailed    = errors.New("SVD factorization failed")
)

type CompressionResult struct {
	Method          string           `json:"method"`
	OriginalSize    int              `json:"originalSize"`
	ComponentLevels []ComponentLevel `json:"componentLevels"`
}

type ComponentLevel struct {
	NumComponents int         `json:"numComponents"`
	DataSize      int         `json:"dataSize"`
	Image         image.Image `json:"-"`
}

type Compressor interface {
	Compress(img image.Image, numComponents int) (*CompressionResult, error)
}

func ImageToGray(img image.Image) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray.Set(x, y, img.At(x, y))
		}
	}
	return gray
}

func ImageToRGBA(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}
	return rgba
}

func PadImage(img *image.Gray, newWidth, newHeight int) *image.Gray {
	bounds := img.Bounds()
	padded := image.NewGray(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			padded.Set(x, y, img.At(x+bounds.Min.X, y+bounds.Min.Y))
		}
	}

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			if x >= bounds.Dx() || y >= bounds.Dy() {
				srcX, srcY := x, y
				if srcX >= bounds.Dx() {
					srcX = bounds.Dx() - 1
				}
				if srcY >= bounds.Dy() {
					srcY = bounds.Dy() - 1
				}
				padded.Set(x, y, padded.At(srcX, srcY))
			}
		}
	}

	return padded
}

func Clamp(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func GenerateSamplePoints(maxComponents int) []int {
	if maxComponents <= 20 {
		points := make([]int, maxComponents)
		for i := range points {
			points[i] = i + 1
		}
		return points
	}

	points := []int{1}
	current := 1
	step := 1

	for current < maxComponents {
		current += step
		if current > maxComponents {
			current = maxComponents
		}
		points = append(points, current)

		if len(points)%3 == 0 {
			step = int(float64(step) * 1.5)
			if step < 1 {
				step = 1
			}
		}
	}

	if points[len(points)-1] != maxComponents {
		points = append(points, maxComponents)
	}

	return points
}

func ImageToGrayscaleMatrix(img image.Image) [][]float64 {
	bounds := img.Bounds()
	matrix := make([][]float64, bounds.Dy())
	for y := 0; y < bounds.Dy(); y++ {
		matrix[y] = make([]float64, bounds.Dx())
		for x := 0; x < bounds.Dx(); x++ {
			r, g, b, _ := img.At(x+bounds.Min.X, y+bounds.Min.Y).RGBA()
			matrix[y][x] = 0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)
		}
	}
	return matrix
}

func MatrixToGrayImage(matrix [][]float64) *image.Gray {
	img := image.NewGray(image.Rect(0, 0, len(matrix[0]), len(matrix)))
	for y := 0; y < len(matrix); y++ {
		for x := 0; x < len(matrix[0]); x++ {
			img.SetGray(x, y, color.Gray{Y: uint8(Clamp(matrix[y][x], 0, 255))})
		}
	}
	return img
}
