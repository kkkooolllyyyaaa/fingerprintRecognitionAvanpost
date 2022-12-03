package algorithm

import "fingerprintRecognitionAvanpost/internal/preprocess"

type KeyPoints struct {
	data []preprocess.Data
}

func NewKeyPointsAlgorithm(data []preprocess.Data) *KeyPoints {
	return &KeyPoints{
		data: data,
	}
}

func (kp *KeyPoints) Predict(toPredict preprocess.Data) (int, error) {
	return 0, nil
}
