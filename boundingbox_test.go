package boundingbox

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"math"
	"os"
	"testing"
)

func TestFind(t *testing.T) {
	testCases := []struct {
		name string
		src  image.Image
		want image.Rectangle
	}{
		{
			"empty",
			&image.Gray{
				Rect:   image.Rect(0, 0, 3, 3),
				Stride: 3,
				Pix: []uint8{
					255, 255, 255,
					255, 255, 255,
					255, 255, 255,
				},
			},
			image.Rect(0, 0, 0, 0),
		},
		{
			"full",
			&image.Gray{
				Rect:   image.Rect(0, 0, 3, 3),
				Stride: 3,
				Pix: []uint8{
					0x0, 0x0, 0x0,
					0x0, 0x0, 0x0,
					0x0, 0x0, 0x0,
				},
			},
			image.Rect(0, 0, 3, 3),
		},
		{
			"1",
			&image.Gray{
				Rect:   image.Rect(0, 0, 1, 1),
				Stride: 1,
				Pix: []uint8{
					255,
				},
			},
			image.Rect(0, 0, 0, 0),
		},
		{
			"0",
			&image.Gray{
				Rect:   image.Rect(0, 0, 1, 1),
				Stride: 1,
				Pix: []uint8{
					0x0,
				},
			},
			image.Rect(0, 0, 1, 1),
		},
		{
			"00",
			&image.Gray{
				Rect:   image.Rect(0, 0, 2, 1),
				Stride: 1,
				Pix: []uint8{
					0x0, 0x0,
				},
			},
			image.Rect(0, 0, 2, 1),
		},
		{
			"11",
			&image.Gray{
				Rect:   image.Rect(0, 0, 2, 1),
				Stride: 1,
				Pix: []uint8{
					255, 255,
				},
			},
			image.Rect(0, 0, 0, 0),
		},
		{
			"0 0",
			&image.Gray{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 1,
				Pix: []uint8{
					0x0,
					0x0,
				},
			},
			image.Rect(0, 0, 1, 2),
		},
		{
			"1 1",
			&image.Gray{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 1,
				Pix: []uint8{
					255,
					255,
				},
			},
			image.Rect(0, 0, 0, 0),
		},
		{
			"simple",
			&image.Gray{
				Rect:   image.Rect(0, 0, 3, 3),
				Stride: 3,
				Pix: []uint8{
					255, 255, 255,
					255, 0x0, 255,
					255, 255, 255,
				},
			},
			image.Rect(1, 1, 2, 2),
		},
		{
			"two points",
			&image.Gray{
				Rect:   image.Rect(0, 0, 5, 5),
				Stride: 5,
				Pix: []uint8{
					255, 255, 255, 255, 255,
					255, 255, 255, 255, 255,
					255, 0x0, 255, 0x0, 255,
					255, 255, 255, 255, 255,
					255, 255, 255, 255, 255,
				},
			},
			image.Rect(1, 2, 4, 3),
		},
		{
			"diamond",
			&image.Gray{
				Rect:   image.Rect(0, 0, 5, 5),
				Stride: 5,
				Pix: []uint8{
					255, 255, 255, 255, 255,
					255, 255, 0x0, 255, 255,
					255, 0x0, 0x0, 0x0, 255,
					255, 255, 0x0, 255, 255,
					255, 255, 255, 255, 255,
				},
			},
			image.Rect(1, 1, 4, 4),
		},
		{
			"x",
			&image.Gray{
				Rect:   image.Rect(0, 0, 5, 5),
				Stride: 5,
				Pix: []uint8{
					255, 255, 255, 255, 255,
					255, 0x0, 255, 0x0, 255,
					255, 255, 0x0, 255, 255,
					255, 0x0, 255, 0x0, 255,
					255, 255, 255, 255, 255,
				},
			},
			image.Rect(1, 1, 4, 4),
		},
		{
			"triangle-lower-left",
			&image.Gray{
				Rect:   image.Rect(0, 0, 5, 5),
				Stride: 5,
				Pix: []uint8{
					0x0, 255, 255, 255, 255,
					0x0, 0x0, 255, 255, 255,
					0x0, 0x0, 0x0, 255, 255,
					0x0, 0x0, 0x0, 0x0, 255,
					0x0, 0x0, 0x0, 0x0, 0x0,
				},
			},
			image.Rect(0, 0, 5, 5),
		},
		{
			"small-triangle-lower-left",
			&image.Gray{
				Rect:   image.Rect(0, 0, 5, 5),
				Stride: 5,
				Pix: []uint8{
					255, 255, 255, 255, 255,
					255, 0x0, 255, 255, 255,
					255, 0x0, 0x0, 255, 255,
					255, 0x0, 0x0, 0x0, 255,
					255, 255, 255, 255, 255,
				},
			},
			image.Rect(1, 1, 4, 4),
		},
		{
			"hline",
			&image.Gray{
				Rect:   image.Rect(0, 0, 5, 5),
				Stride: 5,
				Pix: []uint8{
					255, 255, 255, 255, 255,
					255, 255, 255, 255, 255,
					0x0, 0x0, 0x0, 0x0, 0x0,
					255, 255, 255, 255, 255,
					255, 255, 255, 255, 255,
				},
			},
			image.Rect(0, 2, 5, 3),
		},
		{
			"vline",
			&image.Gray{
				Rect:   image.Rect(0, 0, 5, 5),
				Stride: 5,
				Pix: []uint8{
					255, 255, 0x0, 255, 255,
					255, 255, 0x0, 255, 255,
					255, 255, 0x0, 255, 255,
					255, 255, 0x0, 255, 255,
					255, 255, 0x0, 255, 255,
				},
			},
			image.Rect(2, 0, 3, 5),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Find(tc.src, NewConfigInverse(color.Gray{200}))
			if !got.Eq(tc.want) {
				t.Fatalf("got result %#v want %#v", got, tc.want)
			}
		})
	}

	for _, tc := range testCases {
		t.Run(tc.name+"-parallel", func(t *testing.T) {
			config := NewConfigInverse(color.Gray{200})
			config.Parallel = 2
			got := Find(tc.src, config)
			if !got.Eq(tc.want) {
				t.Fatalf("got result %#v want %#v", got, tc.want)
			}
		})
	}

}

func BenchmarkFind(b *testing.B) {
	b.StopTimer()
	var img = GenerateImage()
	b.StartTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Find(img, NewConfigInverse(color.Gray{200}))
	}
}

func BenchmarkFindParalel(b *testing.B) {
	b.StopTimer()
	var img = GenerateImage()
	config := NewConfigInverse(color.Gray{200})
	config.Parallel = ParallelAuto
	b.StartTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Find(img, config)
	}
}

// ImageToGray converts an image.Image into an image.Gray.
func ImageToGray(img image.Image) *image.Gray {
	var (
		bounds = img.Bounds()
		gray   = image.NewGray(bounds)
	)
	rgbimg, _ := img.(image.RGBA64Image)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			gray.SetRGBA64(x, y, rgbimg.RGBA64At(x, y))
		}
	}
	return gray
}

// LoadImage reads and loads an image from a file path.
func LoadImage(path string) (image.Image, error) {
	infile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer infile.Close()
	img, _, err := image.Decode(infile)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// SaveImagePNG save an image to a PNG file.
func SaveImagePNG(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

type Circle struct {
	X, Y, R float64
}

func (c *Circle) Brightness(x, y float64) uint8 {
	var dx, dy float64 = c.X - x, c.Y - y
	d := math.Sqrt(dx*dx+dy*dy) / c.R
	if d > 1 {
		return 0
	} else {
		return 255
	}
}

// GenerateImage, from http://tech.nitoyon.com/en/blog/2015/12/31/go-image-gen/
func GenerateImage() image.Image {
	var w, h int = 1280, 720
	var hw, hh float64 = float64(w / 2), float64(h / 2)
	r := 200.0
	θ := 2 * math.Pi / 3
	cr := &Circle{hw - r*math.Sin(0), hh - r*math.Cos(0), 100}
	cg := &Circle{hw - r*math.Sin(θ), hh - r*math.Cos(θ), 100}
	cb := &Circle{hw - r*math.Sin(-θ), hh - r*math.Cos(-θ), 100}

	m := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			c := color.RGBA{
				cr.Brightness(float64(x), float64(y)),
				cg.Brightness(float64(x), float64(y)),
				cb.Brightness(float64(x), float64(y)),
				255,
			}
			m.Set(x, y, c)
		}
	}
	return m
}
