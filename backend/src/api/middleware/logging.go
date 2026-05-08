package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/log"
)

func Logger(ctx context.Context, next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.FromContext(ctx).With().Str("request_id", uuid.New().String()).Logger()
		ctx = log.WithContext(ctx, &logger)
		r = r.WithContext(ctx)

		logger.Trace().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("user_agent", r.UserAgent()).
			Str("remote_addr", r.RemoteAddr).
			Msg("Received request")

		start := time.Now()

		rw := newResponseWriter(w)

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Truncate(1 * time.Microsecond)
		logger.Trace().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Int("status_code", rw.statusCode).
			Dur("duration_ms", duration).
			Msg("Processed request")
	})
}
