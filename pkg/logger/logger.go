package logger

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"os"
)

type contextKey int

const (
	loggerContextKey contextKey = iota
)

var logger zerolog.Logger

func init() {
	logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.Level(zerolog.InfoLevel)

	zerolog.TimestampFieldName = "ts"
	zerolog.MessageFieldName = "msg"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func FromContext(ctx context.Context) *zerolog.Logger {
	if log, ok := ctx.Value(loggerContextKey).(*zerolog.Logger); ok {
		return log
	}
	return &logger
}

func ToContext(ctx context.Context, l *zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, l)
}

func Level() zerolog.Level {
	return logger.GetLevel()
}

func SetLevel(l zerolog.Level) {
	logger.Level(l)
}

func Debug(ctx context.Context) *zerolog.Event {
	return FromContext(ctx).Debug()
}

func Info(ctx context.Context) *zerolog.Event {
	return FromContext(ctx).Info()
}

func Warn(ctx context.Context) *zerolog.Event {
	return FromContext(ctx).Warn()
}

func Error(ctx context.Context) *zerolog.Event {
	return FromContext(ctx).Error()
}

func Fatal(ctx context.Context) *zerolog.Event {
	return FromContext(ctx).Fatal()
}

func Panic(ctx context.Context) *zerolog.Event {
	return FromContext(ctx).Panic()
}
