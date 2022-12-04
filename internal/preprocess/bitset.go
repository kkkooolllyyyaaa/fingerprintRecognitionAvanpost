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

func (b *Bitset) Print() {
	for i := 0; i < b.H; i++ {
		for j := 0; j < b.W; j++ {
			if b.Bin[i][j] {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func (b *Bitset) PrintWithDots(dots [][]int) {
	for i := 0; i < b.H; i++ {
		for j := 0; j < b.W; j++ {
			if b.Bin[i][j] {

				isDot := false
				for k := 0; k < len(dots[0]); k++ {
					if dots[0][k] == i && dots[1][k] == j {
						isDot = true
						break
					}
				}
				if isDot {
					fmt.Print("*")
				} else {
					fmt.Print(".")
				}
			}
		}
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

func (bitset *Bitset) ColorModel() color.Model {
	return &BlackWhiteColorModel{}
}

func (bitset *Bitset) Bounds() image.Rectangle {
	return image.Rect(0, 0, bitset.W, bitset.H)
}

func (bitset *Bitset) At(x, y int) color.Color {
	return &BlackWhiteColor{
		isBlack: bitset.Bin[y][x],
	}
}
