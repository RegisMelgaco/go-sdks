package httpresp

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

type loggerCtxKey struct{}

func Logger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), loggerCtxKey{}, logger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
