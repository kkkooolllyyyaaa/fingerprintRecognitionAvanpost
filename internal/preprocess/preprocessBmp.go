package preprocess

import (
	"context"
	"fingerprintRecognitionAvanpost/internal/myimage"
	"fingerprintRecognitionAvanpost/internal/services"
	"fingerprintRecognitionAvanpost/internal/threshold"
	"fingerprintRecognitionAvanpost/pkg/logger"
	"github.com/pkg/errors"
	"golang.org/x/image/bmp"
	"image"
	"os"
)

func PreprocessImages(ctx context.Context, images chan *myimage.MyImage, datas chan *Data, write bool) error {
	for {
		select {
		case <-ctx.Done():
			logger.Info(ctx).Msg("preprocess context done")
			return ctx.Err()
		case img, ok := <-images:
			if !ok {
				logger.Warn(ctx).Msg("Stop listening images channel")
				return nil
			}

			preprocessed, err := PreprocessOne(ctx, img, write)
			if err != nil {
				logger.Error(ctx).Err(err).Msg("Got error while preprocessing")
				return errors.Wrap(err, "Preprocess one")
			}

			datas <- preprocessed

		}
	}
}

func PreprocessOne(ctx context.Context, img *myimage.MyImage, write bool) (*Data, error) {
	WhiteThreshold = threshold.OtsuThresholdValue(img.Img) / 2
	//temp, err := threshold.OtsuThreshold(img.Img, threshold.ThreshBinary)
	//img.Img = temp
	bitset, err := toBitset(img.Img)
	if err != nil {
		logger.Error(ctx).Err(err).Msg("Got error while converting to bitset")
		return nil, errors.Wrap(err, "toBitset")
	}
	if write {
		_ = WriteBmp(bitset, img.Filename)
	}

	services.Skeleton(bitset.Bin)
	if write {
		_ = WriteBmp(bitset, "skeleton_"+img.Filename)
	}

	counted, _ := services.DefineDots(bitset.Bin)

	return NewData(counted, img.Filename), nil
}

func WriteBmp(img image.Image, filename string) error {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		return errors.Wrap(err, "Creating new file")
	}

	err = bmp.Encode(f, img)
	if err != nil {
		return errors.Wrap(err, "Encoding bmp")
	}

	return nil
}
