package histogram

import (
	"fingerprintRecognitionAvanpost/internal/utils"
	"image"
	"image/color"
)

const BufferSize = 256

func HistogramGray(img *image.Gray) [BufferSize]uint64 {
	var res [BufferSize]uint64
	utils.ForEachGrayPixel(img, func(pixel color.Gray) {
		res[pixel.Y]++
	})
	return res
}
