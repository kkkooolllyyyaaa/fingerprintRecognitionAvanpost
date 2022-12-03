package file

import (
	"context"
	"fingerprintRecognitionAvanpost/pkg/logger"
	"github.com/pkg/errors"
	"golang.org/x/image/bmp"
	"image"
	"os"
	"path/filepath"
	"sync/atomic"
)

type Worker struct {
	fileRoot string
	FilesCnt int32
}

func New(fileRoot string) *Worker {
	return &Worker{
		fileRoot: fileRoot,
		FilesCnt: 0,
	}
}

func (w *Worker) ExtractFilePaths(ctx context.Context, fileNames []chan string, workersCnt int) error {
	filePathErr := filepath.Walk(w.fileRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "Walk incoming error")
		}

		if info.IsDir() {
			return nil
		}

		atomic.AddInt32(&w.FilesCnt, 1)

		fileName := info.Name()
		personIndex := ExtractNumberFromFileName(fileName)
		calculatedChan := fileNames[personIndex%workersCnt]
		calculatedChan <- fileName
		logger.Info(ctx).Msgf("Wrote %s to %d chan, curSize=%d", fileName, personIndex%workersCnt, len(calculatedChan))
		return nil
	})

	return filePathErr
}

func (w *Worker) ReadImages(ctx context.Context, fileNamesChan <-chan string, imagesChan chan<- image.Image) error {
	for {
		select {

		case fileName := <-fileNamesChan:
			if isBadFilename(fileName) {
				continue
			}

			img, err := w.decodeToBmp(fileName)
			if err != nil {
				return errors.Wrap(err, "Binarization")
			}

			logger.Info(ctx).Msgf("Wrote img %s", fileName)
			imagesChan <- img

		case <-ctx.Done():
			logger.Warn(ctx).Msg("ReadImages context done")
			return ctx.Err()

		}
	}
}

func (w *Worker) decodeToBmp(path string) (image.Image, error) {
	filePath := w.fileRoot + path

	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "opening file")
	}

	img, err := bmp.Decode(file)
	if err != nil {
		return nil, errors.Wrap(err, "decoding bmp")
	}

	if err = file.Close(); err != nil {
		return nil, errors.Wrap(err, "File close error")
	}

	return img, nil
}
