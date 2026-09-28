package httpx

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/constants"
	"github.com/m-milek/leszmonitor/platform/log"
)

func Logger(ctx context.Context, next http.Handler) http.Handler {
	baseLogger := log.FromContext(ctx)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New()
		w.Header().Set(constants.HTTPHeaderXTraceID, traceID.String())

		reqCtx := log.WithTraceID(r.Context(), traceID)

		logger := baseLogger.With().Str("trace_id", traceID.String()).Logger()
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
