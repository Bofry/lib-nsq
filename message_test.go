package nsq

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/nsqio/go-nsq"
)

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
