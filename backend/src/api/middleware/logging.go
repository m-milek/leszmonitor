package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/log"
)

func Logger(ctx context.Context, next http.Handler) http.Handler {
	baseLogger := log.FromContext(ctx)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()

		reqCtx := context.WithValue(r.Context(), "request_id", reqID)

		logger := baseLogger.With().Str("request_id", reqID).Logger()
		reqCtx = log.WithContext(reqCtx, &logger)

		r = r.WithContext(reqCtx)

		shouldLog := strings.HasPrefix(r.URL.Path, "/api/")

		if shouldLog {
			logger.Trace().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("user_agent", r.UserAgent()).
				Str("remote_addr", r.RemoteAddr).
				Msg("Received request")
		}

		start := time.Now()

		rw := newResponseWriter(w)

		next.ServeHTTP(rw, r)

		if shouldLog {
			duration := time.Since(start).Truncate(1 * time.Microsecond)
			logger.Trace().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Int("status_code", rw.statusCode).
				Dur("duration_ms", duration).
				Msg("Processed request")
		}
	})
}

func Recoverer(ctx context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger := log.FromContext(r.Context())
				logger.Error().Any("panic", err).Msg("Panic recovered")
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
