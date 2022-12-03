package main

import (
	"context"
	"fingerprintRecognitionAvanpost/internal/app"
	"fingerprintRecognitionAvanpost/pkg/logger"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	erg, ctx := errgroup.WithContext(ctxWithCancel)
	logger.Info(ctx).Msg("Starting train main...")

	// graceful shutdown, listen for os signals
	erg.Go(func() error {
		signalsListenChan := make(chan os.Signal, 1)
		signal.Notify(signalsListenChan, syscall.SIGINT, syscall.SIGTERM)

		logger.Info(ctx).Msg("Listening for system signals...")
		select {
		case sig := <-signalsListenChan:
			logger.Warn(ctx).Msgf("Received signal: %s, context will be cancelled\n", sig)
			cancel()
		case <-ctx.Done():
			logger.Debug(ctx).Msg("context done")
			return ctx.Err()
		}

		return nil
	})

	// run train application
	erg.Go(func() error {
		logger.Info(ctx).Msg("Running application...")
		return app.Run(ctx, erg)
	})

	// handle errors
	err := erg.Wait()
	if err != nil {
		logger.Error(ctx).Err(err).Msg("Handling errors...")
		if errors.Is(err, context.Canceled) {
			logger.Warn(ctx).Err(err).Msg("Context was cancelled")
		} else {
			logger.Error(ctx).Err(err).Msg("Received error while application runtime")
		}
	} else {
		logger.Debug(ctx).Msg("Application finished gracefully")
	}
}
