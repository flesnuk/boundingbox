package boundingbox

import (
	"image"
	"runtime"
	"sync"
)

// Find the bounding box of the img using the config provided.
func Find(img image.Image, config Config) image.Rectangle {
	if config.Parallel == ParallelAuto {
		config.Parallel = runtime.NumCPU()
	}
	config.setDefaultPixHasData()
	if config.Parallel > 1 {
		return findParalel(img, config)
	}
	return find(img, config)
}

func find(img image.Image, config Config) image.Rectangle {
	if config.Bounds.Empty() {
		config.Bounds = img.Bounds()
	}
	bounds := config.Bounds
	bbox := image.Rectangle{image.Point{bounds.Max.X, -1}, image.Point{bounds.Min.X, bounds.Min.Y}}
	hasData := false

	width, height := bounds.Max.X, bounds.Max.Y
	rgbimg, _ := img.(image.RGBA64Image)

	for y := bounds.Min.Y; y < height; y++ {
		hasData = false
		for x := bounds.Min.X; x < width; x++ {
			if config.PixHasData(rgbimg.RGBA64At(x, y), config.Threshold) {
				hasData = true
				if x < bbox.Min.X {
					bbox.Min.X = x
				}
				if x >= bbox.Max.X {
					bbox.Max.X = x + 1
				}
			}
		}
		if hasData {
			if bbox.Min.Y < 0 {
				bbox.Min.Y = y
			}
			bbox.Max.Y = y + 1
		}

	}
	if bbox.Empty() {
		return image.Rectangle{}
	}
	return bbox
}

func findParalel(img image.Image, config Config) image.Rectangle {
	var wg sync.WaitGroup

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	nParts := config.Parallel
	r := make(chan image.Rectangle)

	// Split the work in nParts, but only dividing by the height (to maybe help with cache lines).
	for i := 0; i < 1; i++ {
		for j := 0; j < nParts; j++ {
			wg.Add(1)
			partConfig := config
			partConfig.Bounds = image.Rect(width*i, height/nParts*j, min(width*(i+1)+1, width), min(height/nParts*(j+1)+1, height))
			go func() {
				defer wg.Done()
				r <- find(img, partConfig)
			}()

		}
	}
	go func() {
		wg.Wait()
		close(r)
	}()
	result := image.Rectangle{}
	for v := range r {
		result = result.Union(v)
	}

	return result
}

type number interface {
	int | float64
}

func max[T number](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func min[T number](a, b T) T {
	if a < b {
		return a
	}
	return b
}
