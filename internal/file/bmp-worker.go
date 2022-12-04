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

type BmpWorker struct {
	fileRoot      string
	WalkFilesCnt  int32
	ReadImagesCnt int32
}

func NewBmpWorker(fileRoot string) *BmpWorker {
	return &BmpWorker{
		fileRoot:      fileRoot,
		WalkFilesCnt:  0,
		ReadImagesCnt: 0,
	}
}

func (bw *BmpWorker) ExtractFilePaths(ctx context.Context, fileNames []chan string, workersCnt int) error {
	filePathErr := filepath.Walk(bw.fileRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "Walk incoming error")
		}

		if info.IsDir() {
			return nil
		}

		atomic.AddInt32(&bw.WalkFilesCnt, 1)

		fileName := info.Name()
		personIndex := ExtractNumberFromFileName(fileName)
		calculatedChan := fileNames[personIndex%workersCnt]
		calculatedChan <- fileName
		//logger.Info(ctx).Msgf("Wrote %s to %d chan, curSize=%d", fileName, personIndex%workersCnt, len(calculatedChan))
		return nil
	})

	return filePathErr
}

func (bw *BmpWorker) ReadImages(ctx context.Context, fileNamesChan <-chan string, imagesChan chan<- image.Image) error {
	for {
		select {

		case fileName, ok := <-fileNamesChan:
			if !ok {
				logger.Warn(ctx).Msg("Stop listening filenames channel")
				return nil
			}

			if isBadFilename(fileName) {
				continue
			}

			img, err := bw.ExtractImage(fileName)
			if err != nil {
				return errors.Wrap(err, "Binarization error")
			}

			//logger.Info(ctx).Msgf("Wrote img %s", fileName)

			atomic.AddInt32(&bw.ReadImagesCnt, 1)

			imagesChan <- img

		case <-ctx.Done():
			logger.Warn(ctx).Msg("ReadImages context done")
			return ctx.Err()

		}
	}
}

func (bw *BmpWorker) ExtractImage(path string) (image.Image, error) {
	filePath := bw.fileRoot + path

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

//func (bw *BmpWorker) WriteToBmp(bitset *preprocess.Bitset) error {
//}
