package mangabox

import (
	"image"
	"image/color"

	"github.com/flesnuk/boundingbox"
)

// ClampRectangle returns the rectangle moved by p while not exceeding the limits
// https://stackoverflow.com/a/64250050
// TODO

// AdjustSimpleBbox
func AdjustSimpleBbox(img image.Image, minMarginPct float64, maxMarginPct float64) image.Rectangle {
	bounds := img.Bounds()

	minMargin := image.Point{int(minMarginPct*float64(bounds.Max.X) + 0.5), int(minMarginPct*float64(bounds.Max.Y) + 0.5)}
	maxMargin := image.Point{int(maxMarginPct*float64(bounds.Max.X) + 0.5), int(maxMarginPct*float64(bounds.Max.Y) + 0.5)}
	bbox := boundingbox.Find(img, boundingbox.NewConfigInverse(color.Gray{200}))

	newbbox := image.Rectangle{
		image.Point{
			X: Max(0, Min(maxMargin.X, bbox.Min.X-minMargin.X)),
			Y: Max(0, Min(maxMargin.Y, bbox.Min.Y-minMargin.Y)),
		},
		image.Point{
			X: Min(bounds.Max.X, Max(bounds.Max.X-maxMargin.X, bbox.Max.X+minMargin.X)),
			Y: Min(bounds.Max.Y, Max(bounds.Max.Y-maxMargin.Y, bbox.Max.Y+minMargin.Y)),
		},
	}

	return newbbox

}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
