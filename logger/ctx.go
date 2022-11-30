package logger

import (
	"context"

	"go.uber.org/zap"
)

type loggerCtxKey struct{}

func AddToCtx(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, logger)
}
