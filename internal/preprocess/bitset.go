package preprocess

import (
	"fmt"
	"image"
	"image/color"
)

type Bitset struct {
	H   int // y
	W   int // x
	Bin [][]bool
}

func NewZeroes(height, width int) *Bitset {
	bin := make([][]bool, height)
	for i := range bin {
		bin[i] = make([]bool, width)
	}

	return &Bitset{
		H:   height,
		W:   width,
		Bin: bin,
	}
}

type BlackWhiteColor struct {
	isBlack bool
}

func (bwColor *BlackWhiteColor) RGBA() (r, g, b, a uint32) {
	if bwColor.isBlack {
		return 0 * NormalizeK, 0 * NormalizeK, 0 * NormalizeK, 255 * NormalizeK
	}
	return 255 * NormalizeK, 255 * NormalizeK, 255 * NormalizeK, 255 * NormalizeK
}

type BlackWhiteColorModel struct{}

func (bwColorModel *BlackWhiteColorModel) Convert(c color.Color) color.Color {
	return c
}

// ColorModel returns the Image's color model.
func (bitset *Bitset) ColorModel() color.Model {
	return &BlackWhiteColorModel{}
}

// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (bitset *Bitset) Bounds() image.Rectangle {
	return image.Rect(0, 0, bitset.W, bitset.H)
}

// At returns the color of the pixel at (x, y).
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (bitset *Bitset) At(x, y int) color.Color {
	return &BlackWhiteColor{
		isBlack: bitset.Bin[y][x],
	}
}

func (bitset *Bitset) Print() {
	for i := 0; i < bitset.H; i++ {
		for j := 0; j < bitset.W; j++ {
			if bitset.Bin[i][j] {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
