package nsq

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Bofry/lib-nsq/tracing"
)

var _ MessageHandleProc = StopRecursiveForwardUnhandledMessageHandler

func StopRecursiveForwardUnhandledMessageHandler(ctx *ConsumeContext, msg *Message) error {
	ctx.logger.Fatal("invalid forward; it might be recursive forward message to unhandledMessageHandler")
	return nil
}

var (
	_ context.Context = new(ConsumeContext)
	_ tracing.Context = new(ConsumeContext)
)

type ConsumeContext struct {
	Topic string

	logger *log.Logger
	values map[interface{}]interface{}

	unhandledMessageHandler MessageHandleProc

	valuesOnce sync.Once
}

// Deadline implements context.Context.
func (*ConsumeContext) Deadline() (deadline time.Time, ok bool) {
	panic("unimplemented")
}

// Done implements context.Context.
func (*ConsumeContext) Done() <-chan struct{} {
	panic("unimplemented")
}

// Err implements context.Context.
func (*ConsumeContext) Err() error {
	panic("unimplemented")
}

// Value implements context.Context.
func (*ConsumeContext) Value(key interface{}) interface{} {
	panic("unimplemented")
}

func (c *ConsumeContext) SetValue(key, value interface{}) {
	if key == nil {
		return
	}
	if c.values == nil {
		c.valuesOnce.Do(func() {
			if c.values == nil {
				c.values = make(map[interface{}]interface{})
			}
		})
	}
	c.values[key] = value
}

func (c *ConsumeContext) ForwardUnhandledMessage(message *Message) {
	if c.unhandledMessageHandler != nil {
		ctx := &ConsumeContext{
			Topic:                   c.Topic,
			logger:                  c.logger,
			unhandledMessageHandler: StopRecursiveForwardUnhandledMessageHandler,
		}
		c.unhandledMessageHandler(ctx, message)
	}
}
