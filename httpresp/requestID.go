package httpresp

import (
	"net/http"

	"github.com/google/uuid"
	loggerPkg "github.com/regismelgaco/go-sdks/logger"
	"go.uber.org/zap"
)

const requestIDKey = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(requestIDKey)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		w.Header().Set(requestIDKey, requestID)

		logger := loggerPkg.FromContext(r.Context())
		logger = logger.With(zap.String("request_id", requestID))
		ctx := loggerPkg.AddToCtx(r.Context(), logger)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
