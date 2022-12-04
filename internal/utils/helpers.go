package utils

import (
	"image"
	"image/color"
)

func ForEachPixel(size image.Point, f func(x int, y int)) {
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			f(x, y)
		}
	}
}

func ForEachGrayPixel(img *image.Gray, f func(pixel color.Gray)) {
	ForEachPixel(img.Bounds().Size(), func(x, y int) {
		pixel := img.GrayAt(x, y)
		f(pixel)
	})
}
