package nsq

import "github.com/nsqio/go-nsq"

var _ MessageHandleProc = StopRecursiveForwardUnhandledMessageHandler

func StopRecursiveForwardUnhandledMessageHandler(message *Message) error {
	logger.Fatal("invalid forward; it might be recursive forward message to unhandledMessageHandler")
	return nil
}

type Message struct {
	*nsq.Message

	Topic string

	unhandledMessageHandler MessageHandleProc
}

func (m *Message) ForwardUnhandledMessage(message *Message) {
	if message.unhandledMessageHandler != nil {
		handler := message.unhandledMessageHandler

		message.unhandledMessageHandler = StopRecursiveForwardUnhandledMessageHandler

		handler(message)
	}
}
