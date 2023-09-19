package logger

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

const (
	TraceIDAttr = "trace_id"
	SpanIDAttr  = "span_id"
)

type ContextHandler struct {
	slog.Handler
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	nh := h.Handler

	if a := h.attrs(ctx); len(a) != 0 {
		nh = h.Handler.WithAttrs(a)
	}

	return nh.Handle(ctx, r)
}

func (h ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return ContextHandler{
		Handler: h.Handler.WithAttrs(attrs),
	}
}

func (h ContextHandler) WithGroup(name string) slog.Handler {
	return ContextHandler{
		Handler: h.Handler.WithGroup(name),
	}
}

func (h ContextHandler) attrs(ctx context.Context) []slog.Attr {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return nil
	}

	attrs := make([]slog.Attr, 0, 2)

	if span.SpanContext().HasTraceID() {
		attrs = append(attrs, slog.String(TraceIDAttr, span.SpanContext().TraceID().String()))
	}
	if span.SpanContext().HasSpanID() {
		attrs = append(attrs, slog.String(SpanIDAttr, span.SpanContext().SpanID().String()))
	}

	return attrs
}
