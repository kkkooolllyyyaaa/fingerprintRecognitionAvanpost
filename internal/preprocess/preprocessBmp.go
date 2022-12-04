package preprocess

import (
	"context"
	"fingerprintRecognitionAvanpost/internal/services"
	"fingerprintRecognitionAvanpost/pkg/logger"
	"github.com/pkg/errors"
	"golang.org/x/image/bmp"
	"image"
	"os"
)

func PreprocessImages(ctx context.Context, images chan image.Image, datas chan *Data) error {
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

func PreprocessOne(ctx context.Context, image image.Image) (*Data, error) {
	bitset, err := toBitset(image)
	if err != nil {
		logger.Error(ctx).Err(err).Msg("Got error while converting to bitset")
		return nil, errors.Wrap(err, "toBitset")
	}
	_ = WriteBmp(bitset, "before.bmp")

	services.Skeleton(bitset.Bin)
	_ = WriteBmp(bitset, "after.bmp")

	counted := services.KeyPointsCount(bitset.Bin)

	return NewData(counted), nil
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
