package tracing

import (
	"context"

	"github.com/Bofry/trace"
)

func injectSpan(ctx context.Context, span *trace.SeveritySpan) {
	if span != nil {
		if v, ok := ctx.(Context); ok {
			v.SetValue(_CONTEXT_KEY_SEVERITY_SPAN, span)
		}
	}
}
