package algorithm

import "fingerprintRecognitionAvanpost/internal/preprocess"

type Algorithm interface {
	Predict(toPredict *preprocess.Data) (string, error)
}
