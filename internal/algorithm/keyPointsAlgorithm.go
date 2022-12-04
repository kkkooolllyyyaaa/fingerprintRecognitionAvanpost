package algorithm

import (
	"fingerprintRecognitionAvanpost/internal/preprocess"
)

type KeyPoints struct {
	data []*preprocess.Data
}

func NewKeyPointsAlgorithm(data []*preprocess.Data) *KeyPoints {
	return &KeyPoints{
		data: data,
	}
}

func (kp *KeyPoints) Predict(toPredict *preprocess.Data) (string, error) {
	var mi int64 = 9223372036854775807
	miIndex := 0

	for i, dt := range kp.data {
		got := dispersionForKeyPoints(dt, toPredict)
		if got < mi {
			mi = got
			miIndex = i
		}
	}
	return kp.data[miIndex].Filename, nil
}

func dispersionForKeyPoints(first *preprocess.Data, second *preprocess.Data) int64 {
	var sum int64 = 0
	for k, v := range first.KeyPoints {
		secV, ok := second.KeyPoints[k]
		if !ok {
			sum += dispersionBetween(v, 0)
		} else {
			sum += dispersionBetween(v, secV)
		}
	}

	for k, v := range second.KeyPoints {
		_, ok := first.KeyPoints[k]
		if !ok {
			sum += dispersionBetween(v, 0)
		}
	}
	return sum
}

func dispersionBetween(first, second int) int64 {
	return int64((first - second) * (first - second))
}

//func dispersionBetween(first, second int) int64 {
//	return int64(math.Abs(float64(first - second)))
//}
