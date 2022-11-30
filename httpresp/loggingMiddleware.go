package httpresp

import (
	"net/http"

	log "github.com/regismelgaco/go-sdks/logger"
	"go.uber.org/zap"
)

func Logger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := log.AddToCtx(r.Context(), logger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
