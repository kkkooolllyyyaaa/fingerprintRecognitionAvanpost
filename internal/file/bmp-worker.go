package file

import (
	"github.com/pkg/errors"
	"golang.org/x/image/bmp"
	"image/color"
	"os"
)

type Worker struct {
	fileRoot string
}

func New(fileRoot string) *Worker {
	return &Worker{
		fileRoot: fileRoot,
	}
}

const TopLeftOffset = 4

func (w *Worker) ReadToBitset(path string) (*Bitset, error) {
	filePath := w.fileRoot + path

	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "opening file")
	}

	image, err := bmp.Decode(file)
	if err != nil {
		return nil, errors.Wrap(err, "decoding bmp")
	}

	bounds := image.Bounds()
	minX := bounds.Min.X
	maxX := bounds.Max.X - TopLeftOffset
	minY := bounds.Min.Y
	maxY := bounds.Max.Y - TopLeftOffset

	bitset := NewZeroes(maxY-minY, maxX-minX)
	for i := minY; i < maxY; i++ {
		for j := minX; j < maxX; j++ {
			r, g, b := normalizeColor(image.At(j, i))
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
