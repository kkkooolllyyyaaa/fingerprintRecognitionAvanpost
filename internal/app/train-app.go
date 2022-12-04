package app

import (
	"context"
	"fingerprintRecognitionAvanpost/internal/file"
	"fingerprintRecognitionAvanpost/internal/myimage"
	"fingerprintRecognitionAvanpost/internal/preprocess"
	"fingerprintRecognitionAvanpost/pkg/logger"
	"golang.org/x/sync/errgroup"
	"sync"
)

const WorkersCnt = 32
const FileNamesChannelBufferSize = 1024
const ImagesChannelBufferSize = 128
const PreprocessedChannelBufferSize = 128

func RunTrain(ctx context.Context, erg *errgroup.Group) []*preprocess.Data {
	fileRoot := "files/train/SOCOFing/Real/"

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

	imagesChannel := make(chan *myimage.MyImage, ImagesChannelBufferSize)
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
		err := preprocess.PreprocessImages(ctx, imagesChannel, preprocessedChannel, false)
		close(preprocessedChannel)
		logger.Info(ctx).Msg("Done preprocessing images")
		return err
	})

	allData := make([]*preprocess.Data, 0)
	for oneData := range preprocessedChannel {
		allData = append(allData, oneData)
	}
	logger.Info(ctx).Msg("All data done")
	return allData
}
