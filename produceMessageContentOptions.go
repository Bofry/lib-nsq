package nsq

import (
	"context"

	"github.com/Bofry/lib-nsq/tracing"
	"github.com/Bofry/trace"
	"go.opentelemetry.io/otel/propagation"
)

var _ ProduceMessageContentOption = new(ProduceMessageContentOptionProc)

type ProduceMessageContentOptionProc func(topic string, msg *MessageContent) error

func (proc ProduceMessageContentOptionProc) apply(topic string, msg *MessageContent) error {
	return proc(topic, msg)
}

func WithTracePropagation(ctx context.Context, propagator propagation.TextMapPropagator) ProduceMessageContentOptionProc {
	return func(topic string, msg *MessageContent) error {
		_ = topic

		carrier := tracing.NewMessageStateCarrier(&msg.State)
		if propagator == nil {
			propagator = trace.GetTextMapPropagator()
		}
		propagator.Inject(ctx, carrier)
		return nil
	}
}
