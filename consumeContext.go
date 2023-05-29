package nsq

import "log"

var _ MessageHandleProc = StopRecursiveForwardUnhandledMessageHandler

func StopRecursiveForwardUnhandledMessageHandler(ctx *ConsumeContext, msg *Message) error {
	ctx.logger.Fatal("invalid forward; it might be recursive forward message to unhandledMessageHandler")
	return nil
}

type ConsumeContext struct {
	Topic string

	logger *log.Logger

	unhandledMessageHandler MessageHandleProc
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
