package nsq

import (
	"time"

	nsq "github.com/nsqio/go-nsq"
)

var _ nsq.MessageDelegate = new(clientMessageDelegate)

type clientMessageDelegate struct {
	message  *Message
	original nsq.MessageDelegate
}

// OnFinish implements MessageDelegate.
func (d *clientMessageDelegate) OnFinish(m *nsq.Message) {
	d.original.OnFinish(m)
	if d.message.Delegate != nil {
		d.message.Delegate.OnFinish(d.message)
	}
}

// OnRequeue implements MessageDelegate.
func (d *clientMessageDelegate) OnRequeue(m *nsq.Message, delay time.Duration, backoff bool) {
	if d.message.Delegate != nil {
		d.message.Delegate.OnRequeue(d.message, delay, backoff)
	}
	d.original.OnRequeue(m, delay, backoff)
}

// OnTouch implements MessageDelegate.
func (d *clientMessageDelegate) OnTouch(m *nsq.Message) {
	d.original.OnTouch(m)
	if d.message.Delegate != nil {
		d.message.Delegate.OnTouch(d.message)
	}
}
