package httpx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/constants"
	"github.com/m-milek/leszmonitor/platform/log"
	"github.com/rs/zerolog"
)

func serveWithTraceHeader(t *testing.T, header string) (uuid.UUID, string) {
	t.Helper()
	logger := zerolog.Nop()
	ctx := log.WithContext(context.Background(), &logger)

	var ctxTraceID uuid.UUID
	handler := Logger(ctx, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := log.TraceIDFromContext(r.Context())
		if !ok {
			t.Fatal("trace ID not found in request context")
		}
		ctxTraceID = id
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	if header != "" {
		req.Header.Set(constants.HTTPHeaderXTraceID, header)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	return ctxTraceID, rec.Header().Get(constants.HTTPHeaderXTraceID)
}

func TestLogger_GeneratesTraceID(t *testing.T) {
	ctxTraceID, respHeader := serveWithTraceHeader(t, "")

	if ctxTraceID == uuid.Nil {
		t.Fatal("expected generated trace ID")
	}
	if respHeader != ctxTraceID.String() {
		t.Fatalf("response header %q does not match context trace ID %q", respHeader, ctxTraceID)
	}
}

func TestLogger_IgnoresIncomingTraceID(t *testing.T) {
	incoming := uuid.New()
	ctxTraceID, respHeader := serveWithTraceHeader(t, incoming.String())

	if ctxTraceID == incoming {
		t.Fatal("expected incoming trace ID to be ignored")
	}
	if respHeader != ctxTraceID.String() {
		t.Fatalf("response header %q does not match context trace ID %q", respHeader, ctxTraceID)
	}
}
