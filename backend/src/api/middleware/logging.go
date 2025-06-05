package middleware

import (
	"github.com/m-milek/leszmonitor/logger"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Api.Trace().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Msg("Received request")

		start := time.Now()

		rw := newResponseWriter(w)

		// Call the next handler in the chain
		next.ServeHTTP(rw, r)

		duration := time.Since(start).Truncate(1 * time.Microsecond)
		logger.Api.Trace().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Int("status_code", rw.statusCode).
			Dur("duration_ms", duration).
			Msg("Processed request")
	})
}
