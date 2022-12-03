package app

import (
	"context"
	"fingerprintRecognitionAvanpost/internal/file"
	"fingerprintRecognitionAvanpost/internal/preprocess"
	"golang.org/x/sync/errgroup"
	"image"
	"sync"
)

const WorkersCnt = 4
const FileNamesChannelBufferSize = 10
const ImagesChannelBufferSize = 100
const PreprocessedChannelBufferSize = 6000

func RunTrain(ctx context.Context, erg *errgroup.Group) error {
	//fileRoot := "files/train/SOCOFing/OneImage/"
	//fileRoot := "files/train/SOCOFing/TenPeople/"
	fileRoot := "files/train/SOCOFing/Real/"
	//fileRoot := "files/train/SOCOFing/Altered/Altered-Hard/"

	fileNamesChannels := make([]chan string, WorkersCnt)
	for i := range fileNamesChannels {
		fileNamesChannels[i] = make(chan string, FileNamesChannelBufferSize)
	}

	fileWorker := file.New(fileRoot)
	erg.Go(func() error {
		err := fileWorker.ExtractFilePaths(ctx, fileNamesChannels, WorkersCnt)
		for i := range fileNamesChannels {
			close(fileNamesChannels[i])
		}
		return err
	})

	imagesChannel := make(chan image.Image, ImagesChannelBufferSize)
	var wg sync.WaitGroup
	wg.Add(WorkersCnt)

	for j := 0; j < WorkersCnt; j++ {
		neededChan := fileNamesChannels[j]
		erg.Go(func() error {
			err := fileWorker.ReadImages(ctx, neededChan, imagesChannel)
			wg.Done()
			return err
		})
	}

	go func() {
		wg.Wait()
		close(imagesChannel)
	}()

	preprocessedChannel := make(chan *preprocess.Bitset, PreprocessedChannelBufferSize)
	erg.Go(func() error {
		defer close(preprocessedChannel)
		return preprocess.PreprocessImages(ctx, imagesChannel, preprocessedChannel)
	})

	return nil
}

//func SetupRouter(service currency.Service) *gin.Engine {
//	/**
//	@description Init Router
//	*/
//	router := gin.Default()
//	/**
//	@description Setup Mode Application
//	*/
//	gin.SetMode(gin.ReleaseMode)
//
//	/**
//	@description Setup Middleware
//	*/
//	router.Use(cors.New(cors.Config{
//		AllowOrigins:  []string{"*"},
//		AllowMethods:  []string{"*"},
//		AllowHeaders:  []string{"*"},
//		AllowWildcard: true,
//	}))
//	router.Use(helmet.Default())
//	router.Use(gzip.Gzip(gzip.BestCompression))
//	/**
//	@description Init All Route
//	*/
//	route.InitCurrencyRoutes(service, router)
//
//	return route
//}
