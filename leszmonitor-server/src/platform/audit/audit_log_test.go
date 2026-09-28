package audit

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/log"
)

func TestNewAuditLogEntry_TraceIDFromContext(t *testing.T) {
	traceID := uuid.New()
	ctx := log.WithTraceID(context.Background(), traceID)

	entry, err := NewAuditLogEntry(ctx, AuditLogParams{Action: ActionCreateTag, IsSuccess: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.TraceID == nil || *entry.TraceID != traceID.String() {
		t.Fatalf("expected trace ID %q, got %v", traceID, entry.TraceID)
	}
}

func TestNewAuditLogEntry_NoTraceIDInContext(t *testing.T) {
	entry, err := NewAuditLogEntry(context.Background(), AuditLogParams{Action: ActionCreateTag, IsSuccess: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.TraceID != nil {
		t.Fatalf("expected nil trace ID, got %q", *entry.TraceID)
	}
}
