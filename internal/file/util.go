package file

import (
	"image"
	"image/draw"
	"strconv"
	"strings"
)

func ExtractNumberFromFileName(filename string) int {
	numberString := ""
	i := 0
	for i < len(filename) && filename[i] >= '0' && filename[i] <= '9' {
		numberString += string(filename[i])
		i++
	}

	number, _ := strconv.Atoi(numberString)
	return number
}

func isBadFilename(filename string) bool {
	return len(strings.TrimSpace(filename)) == 0 || len(filename) == 0
}

func ToGray(img image.Image) *image.Gray {
	result := image.NewGray(img.Bounds())
	draw.Draw(result, result.Bounds(), img, img.Bounds().Min, draw.Src)
	return result
}
