package logger

import (
	"context"

	"go.uber.org/zap"
)

type loggerCtxKey struct{}

func AddToCtx(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, logger)
}

func FromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(loggerCtxKey{}).(*zap.Logger)
	if !ok {
		panic("failed to retrieve logger from context")
	}

	return logger
}
