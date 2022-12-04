package main

import (
	"context"
	"fingerprintRecognitionAvanpost/internal/app"
	"fingerprintRecognitionAvanpost/pkg/logger"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	erg, ctx := errgroup.WithContext(ctxWithCancel)
	logger.Info(ctx).Msg("Starting train main...")
	var server = &http.Server{Addr: ":8080", Handler: nil}

	// graceful shutdown, listen for os signals
	erg.Go(func() error {
		signalsListenChan := make(chan os.Signal, 1)
		signal.Notify(signalsListenChan, syscall.SIGINT, syscall.SIGTERM)

		logger.Info(ctx).Msg("Listening for system signals...")
		select {
		case sig := <-signalsListenChan:
			_ = server.Shutdown(ctx)
			logger.Warn(ctx).Msgf("Received signal: %s, context will be cancelled\n", sig)
			cancel()
		case <-ctx.Done():
			_ = server.Shutdown(ctx)
			logger.Debug(ctx).Msg("context done")
			return ctx.Err()
		}

		return nil
	})

	//// run train application
	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//erg.Go(func() error {
	//	logger.Info(ctx).Msg("Running train application...")
	//	trainErr := app.RunTrain(ctx, erg, allData)
	//	logger.Info(ctx).Msg("Train done")
	//	wg.Done()
	//	return trainErr
	//})
	//wg.Wait()

	// run train application
	logger.Info(ctx).Msg("Running train application...")
	allData := app.RunTrain(ctx, erg)
	logger.Info(ctx).Msg("Train done")

	// run test application
	erg.Go(func() error {
		logger.Info(ctx).Msg("Running test application...")
		return app.RunTest(ctx, erg, server, allData)
	})

	// handle errors
	err := erg.Wait()
	if err != nil {
		_ = server.Shutdown(ctx)
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
