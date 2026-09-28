package log

import (
	"context"

	"github.com/google/uuid"
)

type traceIDKey struct{}

func WithTraceID(ctx context.Context, traceID uuid.UUID) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

func TraceIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	traceID, ok := ctx.Value(traceIDKey{}).(uuid.UUID)
	if !ok {
		return uuid.Nil, false
	}
	return traceID, true
}
