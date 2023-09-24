package nsq

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

var _ MessageDelegate = new(mockMessageDelegate)

type mockMessageDelegate struct {
	OnFinishCalledCount  int
	OnRequeueCalledCount int
	OnTouchCalledCount   int
}

// OnFinish implements MessageDelegate.
func (d *mockMessageDelegate) OnFinish(*Message) {
	d.OnFinishCalledCount++
}

// OnRequeue implements MessageDelegate.
func (d *mockMessageDelegate) OnRequeue(m *Message, delay time.Duration, backoff bool) {
	d.OnRequeueCalledCount++
}

// OnTouch implements MessageDelegate.
func (d *mockMessageDelegate) OnTouch(*Message) {
	d.OnTouchCalledCount++
}

var _ nsq.MessageDelegate = new(mockNsqMessageDelegate)

type mockNsqMessageDelegate struct {
	OnFinishCalledCount  int
	OnRequeueCalledCount int
	OnTouchCalledCount   int
}

// OnFinish implements nsq.MessageDelegate.
func (d *mockNsqMessageDelegate) OnFinish(*nsq.Message) {
	d.OnFinishCalledCount++
}

// OnRequeue implements nsq.MessageDelegate.
func (d *mockNsqMessageDelegate) OnRequeue(m *nsq.Message, delay time.Duration, backoff bool) {
	d.OnRequeueCalledCount++
}

// OnTouch implements nsq.MessageDelegate.
func (d *mockNsqMessageDelegate) OnTouch(*nsq.Message) {
	d.OnTouchCalledCount++
}

func TestMessage(t *testing.T) {
	var nsqMessage *nsq.Message
	{
		nsqMessage = nsq.NewMessage(
			nsq.MessageID{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'a', 's', 'd', 'f', 'g', 'h'},
			[]byte("foo"),
		)
		nsqMessage.Delegate = new(mockNsqMessageDelegate)
	}

	originalDelegate := nsqMessage.Delegate
	targetDelegate := new(mockMessageDelegate)

	message := &Message{
		Message:  nsqMessage,
		Topic:    "gotest",
		Delegate: targetDelegate,
	}

	nsqMessage.Delegate = &clientMessageDelegate{
		message:  message,
		original: originalDelegate,
	}

	message.Finish()
	{
		var expectedOnFinishCalledCount int = 1
		if expectedOnFinishCalledCount != targetDelegate.OnFinishCalledCount {
			t.Errorf("mockMessageDelegate.OnFinishCalledCount expect:: %v, got:: %v\n", expectedOnFinishCalledCount, targetDelegate.OnFinishCalledCount)
		}
		var expectedOnRequeueCalledCount int = 0
		if expectedOnRequeueCalledCount != targetDelegate.OnRequeueCalledCount {
			t.Errorf("mockMessageDelegate.OnRequeueCalledCount expect:: %v, got:: %v\n", expectedOnRequeueCalledCount, targetDelegate.OnRequeueCalledCount)
		}
		var expectedOnTouchCalledCount int = 0
		if expectedOnTouchCalledCount != targetDelegate.OnTouchCalledCount {
			t.Errorf("mockMessageDelegate.OnTouchCalledCount expect:: %v, got:: %v\n", expectedOnTouchCalledCount, targetDelegate.OnTouchCalledCount)
		}
	}

}

func TestMessage_Content_Well(t *testing.T) {
	message := Message{
		Message: nsq.NewMessage(
			nsq.MessageID{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'a', 's', 'd', 'f', 'g', 'h'},
			bytes.Join([][]byte{
				/* signature           */ {0x1b, 0x4e, 0x53, 0x51},
				/* version             */ {0x01},
				/* reserved            */ {0x00, 0x00},
				/* state total size    */ {0x00, 0x1a},
				/* state tag           */

				/*   tag size          */ {0x00, 0x09},
				/*   tag name          */ {'f', 'o', 'o', '1'},
				/*   tag seperator     */ {':'},
				/*   tag value         */ {'b', 'a', 'r', '1'},

				/*   tag size          */ {0x00, 0x0d},
				/*   tag name          */ {'f', 'o', 'o', '2'},
				/*   tag seperator     */ {':'},
				/*   tag value         */ {'b', 'a', 'r', '2', '-', 'b', 'a', 'z'},

				/* state end delimiter */ {'\r', '\n'},
				/* body                */ []byte("the quick brown fox jumps over the lazy dog"),
			}, nil),
		),
		Topic: "gotest",
	}
	content := message.Content()
	{
		var state map[string][]byte = make(map[string][]byte)
		content.State.Visit(func(name string, value []byte) {
			state[name] = value
		})
		expectedState := map[string][]byte{
			"foo1": []byte("bar1"),
			"foo2": []byte("bar2-baz"),
		}
		if !reflect.DeepEqual(expectedState, state) {
			t.Errorf("MessageContent.State() expected: %v, got: %v", expectedState, state)
		}
	}
	{
		expectedBody := []byte("the quick brown fox jumps over the lazy dog")
		if !reflect.DeepEqual(expectedBody, content.Body) {
			t.Errorf("MessageContent.Body expected: %v, got: %v", expectedBody, content.Body)
		}
	}
}

func TestMessage_Content_Bad(t *testing.T) {
	message := Message{
		Message: nsq.NewMessage(
			nsq.MessageID{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'a', 's', 'd', 'f', 'g', 'h'},
			[]byte("the quick brown fox jumps over the lazy dog"),
		),
		Topic: "gotest",
	}
	content := message.Content()
	{
		var state map[string][]byte = make(map[string][]byte)
		content.State.Visit(func(name string, value []byte) {
			state[name] = value
		})
		expectedState := map[string][]byte{}
		if !reflect.DeepEqual(expectedState, state) {
			t.Errorf("MessageContent.State() expected: %v, got: %v", expectedState, state)
		}
	}
	{
		expectedBody := []byte("the quick brown fox jumps over the lazy dog")
		if !reflect.DeepEqual(expectedBody, content.Body) {
			t.Errorf("MessageContent.Body expected: %v, got: %v", expectedBody, content.Body)
		}
	}
}
