package app

import (
	"context"
	"encoding/json"
	"fingerprintRecognitionAvanpost/internal/algorithm"
	"fingerprintRecognitionAvanpost/internal/file"
	"fingerprintRecognitionAvanpost/internal/myimage"
	"fingerprintRecognitionAvanpost/internal/preprocess"
	"fingerprintRecognitionAvanpost/pkg/logger"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

const WorkersCntCheck = 10
const FileNamesChannelBufferSizeCheck = 1024
const ImagesChannelBufferSizeCheck = 128
const PreprocessedChannelBufferSizeCheck = 128
const PredictWorkersCount = 32

type Answers struct {
	kv map[string]string
}

func RunCheck(ctx context.Context, erg *errgroup.Group, allData []*preprocess.Data) error {
	answersFilePath := "files/test/labels.json"
	f, err := os.ReadFile(answersFilePath)
	if err != nil {
		return errors.Wrap(err, "Couldn't read answers json")
	}
	answrs := &Answers{
		kv: make(map[string]string),
	}
	err = answrs.unmarshalAnswers(f)
	if err != nil {
		return errors.Wrap(err, "Couldn't unmarshal json")
	}

	//fileRoot := "files/train/SOCOFing/Altered/Altered-Easy/"
	fileRoot := "files/test/images/"
	//fileRoot := "files/test/two/"

	fileNamesChannels := make([]chan string, WorkersCntCheck)
	for i := range fileNamesChannels {
		fileNamesChannels[i] = make(chan string, FileNamesChannelBufferSizeCheck)
	}

	fileWorker := file.NewBmpWorker(fileRoot)
	erg.Go(func() error {
		err := fileWorker.ExtractFilePaths(ctx, fileNamesChannels, WorkersCntCheck)
		logger.Info(ctx).Msg("[CHECK] Done extracting file paths")
		for i := range fileNamesChannels {
			logger.Info(ctx).Msgf("[CHECK] Closing channel fileNamesChannels[%d]", i)
			close(fileNamesChannels[i])
		}
		return err
	})

	imagesChannel := make(chan *myimage.MyImage, ImagesChannelBufferSizeCheck)
	var wg sync.WaitGroup
	wg.Add(WorkersCntCheck)

	for j := 0; j < WorkersCntCheck; j++ {
		neededIndex := j
		neededChan := fileNamesChannels[neededIndex]
		logger.Info(ctx).Msgf("[CHECK] Launch read images for %d", neededIndex)
		erg.Go(func() error {
			err := fileWorker.ReadImages(ctx, neededChan, imagesChannel)
			logger.Info(ctx).Msgf("[CHECK] Done reading images for %d", neededIndex)
			wg.Done()
			return err
		})
	}

	go func() {
		wg.Wait()
		logger.Info(ctx).Msg("[CHECK] Wait done, read all images")
		close(imagesChannel)
	}()

	preprocessedChannel := make(chan *preprocess.Data, PreprocessedChannelBufferSizeCheck)
	erg.Go(func() error {
		err := preprocess.PreprocessImages(ctx, imagesChannel, preprocessedChannel, false)
		close(preprocessedChannel)
		logger.Info(ctx).Msg("[CHECK] Done preprocessing images")
		return err
	})

	alg := algorithm.NewKeyPointsAlgorithm(allData)

	var eqCnt int32 = 0
	var allCnt int32 = 0

	var oblErrCnt int32 = 0
	var crErrCnt int32 = 0
	var zcutErrCnt int32 = 0

	logger.Info(ctx).Msg("[CHECK] Starting check files")
	wgPredict := sync.WaitGroup{}
	wgPredict.Add(PredictWorkersCount)
	for i := 0; i < PredictWorkersCount; i++ {
		index := i
		erg.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					wgPredict.Done()
					logger.Info(ctx).Msgf("predict %d context done", index)
					return ctx.Err()
				case oneData, ok := <-preprocessedChannel:
					if !ok {
						wgPredict.Done()
						logger.Warn(ctx).Msgf("Stop predicting in %d", index)
						return nil
					}

					if allCnt%100 == 0 {
						logger.Info(ctx).Msgf("Processed %d", allCnt)
					}

					foundFilename, err := alg.Predict(oneData)
					if err != nil {
						logger.Warn(ctx).Err(err).Msg("[CHECK] Got error while predicting")
						continue
					}

					gotFilename := oneData.Filename
					isEquals := checkEquals(foundFilename, gotFilename, answrs)
					//isEquals := checkEquals(foundFilename, gotFilename)
					if isEquals {
						atomic.AddInt32(&eqCnt, 1)
					} else {
						if strings.Contains(gotFilename, "CR") {
							atomic.AddInt32(&crErrCnt, 1)
						}
						if strings.Contains(gotFilename, "Obl") {
							atomic.AddInt32(&oblErrCnt, 1)
						}
						if strings.Contains(gotFilename, "Zcut") {
							atomic.AddInt32(&zcutErrCnt, 1)
						}
					}
					atomic.AddInt32(&allCnt, 1)
				}
			}
		})
	}

	go func() {
		wgPredict.Wait()
		fmt.Printf("Zcut=%d Obl=%d CR=%d", zcutErrCnt, oblErrCnt, crErrCnt)
		logger.Info(ctx).Msgf("accuracy = %d / %d", eqCnt, allCnt)
	}()
	return nil
}

//func checkEquals(f1, f2 string) bool {
//	return strings.HasPrefix(f2[:len(f2)-4], f1[:len(f1)-4])
//}

func checkEquals(gotAnswer, key string, answers *Answers) bool {
	actualAnswer, ok := answers.kv[key]
	if !ok {
		return false
	}
	actualAnswer += ".BMP"
	return actualAnswer == gotAnswer
}

func (a *Answers) unmarshalAnswers(data []byte) error {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	a.kv = make(map[string]string)
	for k, v := range m {
		var p string
		if err := json.Unmarshal(v, &p); err != nil {
			return err
		}
		a.kv[k] = p
	}
	return nil
}
