package httpx

import (
	"bufio"
	"net"
	"net/http"
)

// responseWriter is a custom ResponseWriter that captures the status code.
type responseWriter struct {
	http.ResponseWriter

	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

// WrapResponseWriter wraps w so the status code is captured and [http.Hijacker] stays available.
func WrapResponseWriter(w http.ResponseWriter) http.ResponseWriter {
	return newResponseWriter(w)
}

// WriteHeader captures the status code.
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Hijack implements the [http.Hijacker] interface for websocket compatibility.
func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}
