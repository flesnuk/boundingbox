# Bounding box

Library for finding the [minimum bounding box](https://en.wikipedia.org/wiki/Minimum_bounding_box) of an image. It does so by simply iterating the image and checking which pixels surpass certain threshold, so they are considered as data points in the image which will delimit the bounding box. 

## Usage

```go
package main

import (
	"image"
	_ "image/png"
	"os"
    "github.com/flesnuk/boundingbox"
)

func main() {
    infile, _ := os.Open("0004-001.png")
    defer infile.Close()
    img, _ := image.Decode(infile)
    bbox := boundingbox.Find(img, boundingbox.NewConfig(color.Gray{30}))
}
```

The NewConfig helper is useful for creating a configuration with the specified threshold color. By default the function used to compare the pixels is defined and exported as `PixHasDataBlackBG`. This does a simple comparison of RGB values, and if all are greater than the color threshold, it returns true.

There is a few options that can be customized:

- PixHasData: set custom function to compare when a pixel is considered as data or not.
- Bounds: by default uses the image bounds, but can be changed with this parameter.
- Parallel: sets the number of goroutines to process the image in parallel. Set to -1 for auto setting (use the runtime.NumCPU())


