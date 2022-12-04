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

func PreprocessImages(ctx context.Context, images chan *myimage.MyImage, datas chan *Data) error {
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

			preprocessed, err := PreprocessOne(ctx, img)
			if err != nil {
				logger.Error(ctx).Err(err).Msg("Got error while preprocessing")
				return errors.Wrap(err, "Preprocess one")
			}

			datas <- preprocessed

		}
	}
}

func PreprocessOne(ctx context.Context, img *myimage.MyImage) (*Data, error) {
	WhiteThreshold = threshold.OtsuThresholdValue(img.Img)
	temp, err := threshold.OtsuThreshold(img.Img, threshold.ThreshBinary)
	img.Img = temp
	bitset, err := toBitset(img.Img)
	if err != nil {
		logger.Error(ctx).Err(err).Msg("Got error while converting to bitset")
		return nil, errors.Wrap(err, "toBitset")
	}
	//_ = WriteBmp(bitset, "before-no-blik.bmp")

	services.Skeleton(bitset.Bin)
	//_ = WriteBmp(bitset, "after-no-blik.bmp")

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
