package file

import "fingerprintRecognitionAvanpost/internal/preprocess"

type DataWorker struct {
	fileRoot string
	FilesCnt int32
}

func NewDataWorker(fileRoot string) *DataWorker {
	return &DataWorker{
		fileRoot: fileRoot,
		FilesCnt: 0,
	}
}

func (pw *DataWorker) InitByReadingAll() ([]preprocess.Data, error) {

	return nil, nil
}
