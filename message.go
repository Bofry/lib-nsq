package nsq

import (
	"sync/atomic"

	"github.com/nsqio/go-nsq"
)

var _ MessageHandleProc = StopRecursiveForwardUnhandledMessageHandler

func StopRecursiveForwardUnhandledMessageHandler(message *Message) error {
	logger.Fatal("invalid forward; it might be recursive forward message to unhandledMessageHandler")
	return nil
}

type Message struct {
	*nsq.Message

	Topic string

	unhandledMessageHandler     MessageHandleProc
	unhandledMessageHandlerFlag *int32
}

func (m *Message) ForwardUnhandledMessageHandler() error {
	if m.unhandledMessageHandler != nil {
		if atomic.CompareAndSwapInt32(m.unhandledMessageHandlerFlag, 0, 1) {
			return m.unhandledMessageHandler(m)
		}
		return StopRecursiveForwardUnhandledMessageHandler(m)
	}
	return nil
}
