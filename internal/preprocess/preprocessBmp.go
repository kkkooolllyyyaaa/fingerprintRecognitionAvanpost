package preprocess

import (
	"context"
	"fingerprintRecognitionAvanpost/pkg/logger"
	"github.com/pkg/errors"
	"image"
)

func PreprocessImages(ctx context.Context, images chan image.Image, bitsets chan *Bitset) error {
	for {
		select {
		case <-ctx.Done():
			logger.Info(ctx).Msg("preprocess context done")
			return ctx.Err()
		case img := <-images:
			select {
			case <-ctx.Done():
				logger.Debug(ctx).Msg("preprocess context done")
				return ctx.Err()
			default:
				bitset, err := toBitset(img)
				if err != nil {
					logger.Error(ctx).Err(err).Msg("Got error while converting to bitset")
					return errors.Wrap(err, "toBitset")
				}

				logger.Info(ctx).Msg("Converted img to bitset")
				bitsets <- bitset
			}
		}
	}
}
