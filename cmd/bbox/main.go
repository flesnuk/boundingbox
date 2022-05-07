package main

import (
	"image/color"

	"github.com/flesnuk/boundingbox"
)

func main() {
	path := ""
	var img, err = LoadImage(path)
	if err != nil {
		panic(err)
	}
	var gray = ImageToGray(img)

	config := boundingbox.NewConfig(color.Gray{200})
	config.Parallel = -1

	gg := gray.SubImage(boundingbox.Find(img, config))

	err = SaveImagePNG(gg, "bbox.png")
	if err != nil {
		panic(err)
	}
}
