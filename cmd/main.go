package main

import (
	"context"
	"fingerprintRecognitionAvanpost/pkg/logger"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	erg, ctx := errgroup.WithContext(ctxWithCancel)
	logger.Info(ctx).Msg("Starts main")

	// graceful shutdown, listen for os signals
	erg.Go(func() error {
		signalsListenChan := make(chan os.Signal, 1)
		signal.Notify(signalsListenChan, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-signalsListenChan:
			logger.Warn(ctx).Msgf("Received signal: %s, context will be cancelled\n", sig)
			cancel()
		case <-ctx.Done():
			logger.Debug(ctx).Msg("cmd.bot.main context done")
			return ctx.Err()
		}

		return nil
	})

	ticker := time.NewTicker(time.Duration(10) * time.Second)
	erg.Go(func() error {
		select {
		case <-ctx.Done():
			logger.Debug(ctx).Msg("context done")
			return ctx.Err()
		case <-ticker.C:
			select {
			case <-ctx.Done():
				logger.Debug(ctx).Msg("context done")
				return ctx.Err()
			default:
				logger.Debug(ctx).Msg("Regular tick...")
				fmt.Println("Hello, World!")
			}
		}
		return ctx.Err()
	})

	// handle errors
	err := erg.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			logger.Warn(ctx).Err(err).Msg("Context was cancelled")
		} else {
			logger.Error(ctx).Err(err).Msg("Received error while application runtime")
		}
	} else {
		logger.Debug(ctx).Msg("Application finished gracefully")
	}
}
