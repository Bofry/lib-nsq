package tracing

import (
	"context"

	"github.com/Bofry/trace"
)

type MessageStateCtxSpanExtractor int

// Extract implements trace.SpanExtractor
func (MessageStateCtxSpanExtractor) Extract(ctx context.Context) *trace.SeveritySpan {
	v := ctx.Value(_CONTEXT_KEY_SEVERITY_SPAN)
	if sp, ok := v.(*trace.SeveritySpan); ok {
		return sp
	}
	return nil
}
