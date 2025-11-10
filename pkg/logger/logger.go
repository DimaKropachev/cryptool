package logger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type Key string

const key Key = "logger"

type Logger struct {
	l *zap.Logger
}

func New(ctx context.Context) (context.Context, error) {
	logger, err := zap.NewProduction()
	
	if err != nil {
		return nil, fmt.Errorf("error creating logger: %w", err)
	}

	ctx = context.WithValue(ctx, key, &Logger{logger})
	return ctx, nil
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	return ctx.Value(key).(*Logger)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.l.Info(msg, fields...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.l.Fatal(msg, fields...)
}
