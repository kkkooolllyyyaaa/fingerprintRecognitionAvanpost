package preprocess

import (
	"image"
	"image/color"
)

const TopLeftOffset = 4

func toBitset(img image.Image) (*Bitset, error) {
	bounds := img.Bounds()
	minX := bounds.Min.X
	maxX := bounds.Max.X - TopLeftOffset
	minY := bounds.Min.Y
	maxY := bounds.Max.Y - TopLeftOffset

	bitset := NewZeroes(maxY-minY, maxX-minX)
	for i := minY; i < maxY; i++ {
		for j := minX; j < maxX; j++ {
			r, g, b := normalizeColor(img.At(j, i))
			bitset.Bin[i][j] = binarizeRgb(r, g, b)
		}
	}

	return bitset, nil
}

const NormalizeK = 256

func normalizeColor(c color.Color) (r, g, b uint32) {
	r, g, b, _ = c.RGBA()
	return r / NormalizeK, g / NormalizeK, b / NormalizeK
}

const WhiteThreshold = 30

func binarizeRgb(r, g, b uint32) bool {
	return r <= WhiteThreshold && g <= WhiteThreshold && b <= WhiteThreshold
}
