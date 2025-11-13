package compression

import (
	"image"
	"image/color"
	"math"
	"runtime"
	"sync"
)

const blockSize = 8

type DCTCompressor struct{}

func NewDCTCompressor() *DCTCompressor {
	return &DCTCompressor{}
}

func (d *DCTCompressor) Compress(img image.Image, numComponents int) (*CompressionResult, error) {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()

	padded := PadImage(ImageToGray(img),
		((width+blockSize-1)/blockSize)*blockSize,
		((height+blockSize-1)/blockSize)*blockSize)

	numBlocksY := padded.Bounds().Dy() / blockSize
	numBlocksX := padded.Bounds().Dx() / blockSize
	dctBlocks := make([][][]float64, numBlocksY)
	for i := range dctBlocks {
		dctBlocks[i] = make([][]float64, numBlocksX)
	}

	numWorkers := runtime.NumCPU()
	var wg sync.WaitGroup
	workChan := make(chan struct{ i, j int }, numBlocksY*numBlocksX)

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for work := range workChan {
				block := extractBlock(padded, work.i*blockSize, work.j*blockSize)
				dctBlocks[work.i][work.j] = dct2D(block)
			}
		}()
	}

	for i := 0; i < numBlocksY; i++ {
		for j := 0; j < numBlocksX; j++ {
			workChan <- struct{ i, j int }{i, j}
		}
	}
	close(workChan)
	wg.Wait()

	if numComponents > blockSize*blockSize {
		numComponents = blockSize * blockSize
	}

	samplePoints := GenerateSamplePoints(numComponents)
	componentLevels := make([]ComponentLevel, len(samplePoints))

	var wg2 sync.WaitGroup
	for idx, k := range samplePoints {
		wg2.Add(1)
		go func(idx, k int) {
			defer wg2.Done()
			reconstructed := dctReconstructWithComponents(dctBlocks, k, width, height)
			componentLevels[idx] = ComponentLevel{
				NumComponents: k,
				DataSize:      k * len(dctBlocks) * len(dctBlocks[0]) * 8,
				Image:         reconstructed,
			}
		}(idx, k)
	}
	wg2.Wait()

	return &CompressionResult{
		Method:          "DCT",
		OriginalSize:    width * height * 3,
		ComponentLevels: componentLevels,
	}, nil
}

func dctReconstructWithComponents(dctBlocks [][][]float64, numComponents, origWidth, origHeight int) image.Image {
	reconstructed := image.NewGray(image.Rect(0, 0, len(dctBlocks[0])*blockSize, len(dctBlocks)*blockSize))

	numWorkers := runtime.NumCPU()
	var wg sync.WaitGroup
	var mu sync.Mutex
	workChan := make(chan struct{ i, j int }, len(dctBlocks)*len(dctBlocks[0]))

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for work := range workChan {
				truncated := truncateComponents(dctBlocks[work.i][work.j], numComponents)
				block := idct2D(truncated)

				mu.Lock()
				for y := 0; y < blockSize; y++ {
					for x := 0; x < blockSize; x++ {
						py := work.i*blockSize + y
						px := work.j*blockSize + x
						if py < origHeight && px < origWidth {
							reconstructed.SetGray(px, py, color.Gray{
								Y: uint8(Clamp(block[y*blockSize+x], 0, 255)),
							})
						}
					}
				}
				mu.Unlock()
			}
		}()
	}

	for i := 0; i < len(dctBlocks); i++ {
		for j := 0; j < len(dctBlocks[0]); j++ {
			workChan <- struct{ i, j int }{i, j}
		}
	}
	close(workChan)
	wg.Wait()

	cropped := image.NewGray(image.Rect(0, 0, origWidth, origHeight))
	for y := 0; y < origHeight; y++ {
		for x := 0; x < origWidth; x++ {
			cropped.SetGray(x, y, reconstructed.GrayAt(x, y))
		}
	}

	return cropped
}

func dct2D(block []float64) []float64 {
	size := int(math.Sqrt(float64(len(block))))
	result := make([]float64, len(block))

	for u := 0; u < size; u++ {
		for v := 0; v < size; v++ {
			sum := 0.0
			for x := 0; x < size; x++ {
				for y := 0; y < size; y++ {
					sum += block[x*size+y] *
						math.Cos((2*float64(x)+1)*float64(u)*math.Pi/(2*float64(size))) *
						math.Cos((2*float64(y)+1)*float64(v)*math.Pi/(2*float64(size)))
				}
			}

			cu, cv := 1.0, 1.0
			if u == 0 {
				cu = 1.0 / math.Sqrt(2)
			}
			if v == 0 {
				cv = 1.0 / math.Sqrt(2)
			}

			result[u*size+v] = 0.25 * cu * cv * sum
		}
	}

	return result
}

func idct2D(dctBlock []float64) []float64 {
	size := int(math.Sqrt(float64(len(dctBlock))))
	result := make([]float64, len(dctBlock))

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			sum := 0.0
			for u := 0; u < size; u++ {
				for v := 0; v < size; v++ {
					cu, cv := 1.0, 1.0
					if u == 0 {
						cu = 1.0 / math.Sqrt(2)
					}
					if v == 0 {
						cv = 1.0 / math.Sqrt(2)
					}

					sum += cu * cv * dctBlock[u*size+v] *
						math.Cos((2*float64(x)+1)*float64(u)*math.Pi/(2*float64(size))) *
						math.Cos((2*float64(y)+1)*float64(v)*math.Pi/(2*float64(size)))
				}
			}

			result[x*size+y] = 0.25 * sum
		}
	}

	return result
}

func truncateComponents(dctBlock []float64, n int) []float64 {
	result := make([]float64, len(dctBlock))
	zigzag := generateZigzagOrder()

	for i := 0; i < n && i < len(zigzag); i++ {
		result[zigzag[i]] = dctBlock[zigzag[i]]
	}

	return result
}

func generateZigzagOrder() []int {
	zigzag := make([]int, blockSize*blockSize)
	idx := 0

	for sum := 0; sum < 2*blockSize-1; sum++ {
		if sum%2 == 0 {
			for i := sum; i >= 0; i-- {
				j := sum - i
				if i < blockSize && j < blockSize {
					zigzag[idx] = i*blockSize + j
					idx++
				}
			}
		} else {
			for j := sum; j >= 0; j-- {
				i := sum - j
				if i < blockSize && j < blockSize {
					zigzag[idx] = i*blockSize + j
					idx++
				}
			}
		}
	}

	return zigzag
}

func extractBlock(img *image.Gray, startY, startX int) []float64 {
	block := make([]float64, blockSize*blockSize)
	for y := 0; y < blockSize; y++ {
		for x := 0; x < blockSize; x++ {
			block[y*blockSize+x] = float64(img.GrayAt(startX+x, startY+y).Y)
		}
	}
	return block
}
