package threshold

import (
	"fingerprintRecognitionAvanpost/internal/histogram"
	"fingerprintRecognitionAvanpost/internal/utils"
	"github.com/pkg/errors"
	"image"
	"image/color"
)

type Method int

const (
	// ThreshBinary
	// if (x, y) > threshold
	//     maxValue
	// else
	//     0
	ThreshBinary Method = iota

	// ThreshBinaryInv
	// if (x, y) > threshold
	//     0
	// else
	//     maxValue
	ThreshBinaryInv
)

func Threshold(img *image.Gray, t uint8, method Method) (*image.Gray, error) {
	var setPixel func(*image.Gray, int, int)
	switch method {
	case ThreshBinary:
		setPixel = func(gray *image.Gray, x int, y int) {
			pixel := img.GrayAt(x, y).Y
			if pixel < t {
				gray.SetGray(x, y, color.Gray{Y: utils.MinUint8})
			} else {
				gray.SetGray(x, y, color.Gray{Y: utils.MaxUint8})
			}
		}
	case ThreshBinaryInv:
		setPixel = func(gray *image.Gray, x int, y int) {
			pixel := img.GrayAt(x, y).Y
			if pixel < t {
				gray.SetGray(x, y, color.Gray{Y: utils.MaxUint8})
			} else {
				gray.SetGray(x, y, color.Gray{Y: utils.MinUint8})
			}
		}
	default:
		return nil, errors.New("invalid threshold method")
	}
	return threshold(img, setPixel), nil
}

// https://en.wikipedia.org/wiki/Otsu%27s_method
func OtsuThreshold(img *image.Gray, method Method) (*image.Gray, error) {
	return Threshold(img, OtsuThresholdValue(img), method)
}

func threshold(img *image.Gray, setPixel func(*image.Gray, int, int)) *image.Gray {
	size := img.Bounds().Size()
	gray := image.NewGray(img.Bounds())
	utils.ParallelForEachPixel(size, func(x, y int) {
		setPixel(gray, x, y)
	})
	return gray
}

func OtsuThresholdValue(img *image.Gray) uint8 {
	hist := histogram.HistogramGray(img)
	size := img.Bounds().Size()
	totalNumberOfPixels := size.X * size.Y

	var sumHist float64
	for i, bin := range hist {
		sumHist += float64(uint64(i) * bin)
	}

	var sumBackground float64
	var weightBackground int
	var weightForeground int

	maxVariance := 0.0
	var thresh uint8
	for i, bin := range hist {
		weightBackground += int(bin)
		if weightBackground == 0 {
			continue
		}
		weightForeground = totalNumberOfPixels - weightBackground
		if weightForeground == 0 {
			break
		}

		sumBackground += float64(uint64(i) * bin)

		meanBackground := float64(sumBackground) / float64(weightBackground)
		meanForeground := (sumHist - sumBackground) / float64(weightForeground)

		variance := float64(weightBackground) * float64(weightForeground) * (meanBackground - meanForeground) * (meanBackground - meanForeground)

		if variance > maxVariance {
			maxVariance = variance
			thresh = uint8(i)
		}
	}
	return thresh
}
