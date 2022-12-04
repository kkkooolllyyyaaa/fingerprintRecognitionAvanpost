package preprocess

import "fmt"

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
				fmt.Print(".")
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
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
