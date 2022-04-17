package boundingbox

import (
	"image"
	"image/color"
)

const ParallelAuto = -1

// Config contains the parameters used.
type Config struct {
	// When a pixel is considered as the edge of the image (has meaningful data).
	Threshold color.RGBA64
	// Check if pixel surpasses the Threshold. If nil it will use a basic pix > threshold function (gray-like comparison).
	PixHasData func(pix color.RGBA64, threshold color.RGBA64) bool
	// Region to check, if empty, image bounds are used instead.
	Bounds image.Rectangle
	// Number of gorutines used if n > 1. Use ParallelAuto (-1) for using NumCPU
	Parallel int
}

// NewConfig returns a default config using the supplied threshold color.
// The PixHasData function is nil and use a basic pix > threshold function (gray-like comparison).
func NewConfig(threshold color.Color) Config {
	ct := color.RGBA64Model.Convert(threshold).(color.RGBA64)
	return Config{Threshold: ct}
}

// NewConfigInverse returns a default config using the supplied threshold color.
// The PixHasData function use a basic pix < threshold function. Useful when background color is white.
func NewConfigInverse(threshold color.Color) Config {
	ct := color.RGBA64Model.Convert(threshold).(color.RGBA64)
	return Config{Threshold: ct, PixHasData: PixHasDataWhiteBG}
}

// WithThreshold sets the threshold limit.
func (config *Config) WithThreshold(threshold color.RGBA64) *Config {
	config.Threshold = threshold
	return config
}

// WithPixHasData sets the pixHasData function.
func (config *Config) WithPixHasData(pixHasData func(pix color.RGBA64, threshold color.RGBA64) bool) *Config {
	config.PixHasData = pixHasData
	return config
}

// WithBounds sets the bounds region.
func (config *Config) WithBounds(bounds image.Rectangle) *Config {
	config.Bounds = bounds
	return config
}

// WithParallel sets the parallel number of goroutines.
func (config *Config) WithParallel(parallel int) *Config {
	config.Parallel = parallel
	return config
}

func (conf *Config) setDefaultPixHasData() {
	if conf.PixHasData == nil {
		conf.PixHasData = PixHasDataBlackBG
	}
}

// PixHasDataBlackBG is a simple function that returns true if pix has greater RGB values than threshold.
func PixHasDataBlackBG(pix color.RGBA64, threshold color.RGBA64) bool {
	r1, g1, b1, _ := pix.RGBA()
	r2, g2, b2, _ := threshold.RGBA()
	return r1 > r2 && g1 > g2 && b1 > b2
}

// PixHasDataWhiteBG is a simple function that returns true if pix has lesser RGB values than threshold.
func PixHasDataWhiteBG(pix color.RGBA64, threshold color.RGBA64) bool {
	r1, g1, b1, _ := pix.RGBA()
	r2, g2, b2, _ := threshold.RGBA()
	return r1 < r2 && g1 < g2 && b1 < b2
}
