package app

import (
	"context"
	"fingerprintRecognitionAvanpost/internal/file"
	"fingerprintRecognitionAvanpost/internal/preprocess"
	"fingerprintRecognitionAvanpost/pkg/logger"
	"golang.org/x/sync/errgroup"
	"image"
	"sync"
)

const WorkersCnt = 4
const FileNamesChannelBufferSize = 20
const ImagesChannelBufferSize = 20
const PreprocessedChannelBufferSize = 20

func RunTrain(ctx context.Context, erg *errgroup.Group) error {
	//fileRoot := "files/train/SOCOFing/OneImage/"
	fileRoot := "files/train/SOCOFing/OneImageTwo/"

	//fileRoot := "files/train/SOCOFing/TenPeople/"
	//fileRoot := "files/train/SOCOFing/Real/"
	//fileRoot := "files/train/SOCOFing/Altered/Altered-Hard/"

	fileNamesChannels := make([]chan string, WorkersCnt)
	for i := range fileNamesChannels {
		fileNamesChannels[i] = make(chan string, FileNamesChannelBufferSize)
	}

	fileWorker := file.NewBmpWorker(fileRoot)
	erg.Go(func() error {
		err := fileWorker.ExtractFilePaths(ctx, fileNamesChannels, WorkersCnt)
		logger.Info(ctx).Msg("Done extracting file paths")
		for i := range fileNamesChannels {
			logger.Info(ctx).Msgf("Closing channel fileNamesChannels[%d]", i)
			close(fileNamesChannels[i])
		}
		return err
	})

	imagesChannel := make(chan *image.Gray, ImagesChannelBufferSize)
	var wg sync.WaitGroup
	wg.Add(WorkersCnt)

	for j := 0; j < WorkersCnt; j++ {
		neededIndex := j
		neededChan := fileNamesChannels[neededIndex]
		logger.Info(ctx).Msgf("Launch read images for %d", neededIndex)
		erg.Go(func() error {
			err := fileWorker.ReadImages(ctx, neededChan, imagesChannel)
			logger.Info(ctx).Msgf("Done reading images for %d", neededIndex)
			wg.Done()
			return err
		})
	}

	go func() {
		wg.Wait()
		logger.Info(ctx).Msg("Wait done, read all images")
		close(imagesChannel)
	}()

	preprocessedChannel := make(chan *preprocess.Data, PreprocessedChannelBufferSize)
	erg.Go(func() error {
		err := preprocess.PreprocessImages(ctx, imagesChannel, preprocessedChannel)
		close(preprocessedChannel)
		logger.Info(ctx).Msg("Done preprocessing images")
		return err
	})
	return nil
}
