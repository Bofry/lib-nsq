package nsq

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func TestMessageContent_WriteTo_Case01(t *testing.T) {
	content := MessageContent{}

	content.State.Set("foo", []byte("bar"))
	content.Body = []byte("the quick brown fox jumps over the lazy dog")

	var (
		payload bytes.Buffer
		writer  = bufio.NewWriter(&payload)
	)

	payloadSize, err := content.WriteTo(writer)
	if err != nil {
		t.Fatalf("should not error, but got: %v", err)
	}
	err = writer.Flush()
	if err != nil {
		t.Fatalf("should not error, but got: %v", err)
	}

	expectedPayloadSize := 63
	if expectedPayloadSize != int(payloadSize) {
		t.Errorf("PayloadSize expected: %v, got: %v", expectedPayloadSize, payloadSize)
	}
	expectedPayload := bytes.Join([][]byte{
		/* signature           */ {0x1b, 0x4e, 0x53, 0x51},
		/* version             */ {0x01},
		/* reserved            */ {0x00, 0x00},
		/* state total size    */ {0x00, 0x09},
		/* state tag           */
		/*   tag size          */ {0x00, 0x07},
		/*   tag name          */ {'f', 'o', 'o'},
		/*   tag seperator     */ {':'},
		/*   tag value         */ {'b', 'a', 'r'},
		/* state end delimiter */ {'\r', '\n'},
		/* body                */ []byte("the quick brown fox jumps over the lazy dog"),
	}, nil)
	if !reflect.DeepEqual(expectedPayload, payload.Bytes()) {
		t.Errorf("Payload expected: %v, got: %v", expectedPayload, payload.Bytes())
	}
}

func TestMessageContent_WriteTo_Case02(t *testing.T) {
	content := MessageContent{}

	content.State.Set("foo1", []byte("bar1"))
	content.State.Set("foo2", []byte("bar2-baz"))
	content.Body = []byte("the quick brown fox jumps over the lazy dog")

	var (
		payload bytes.Buffer
		writer  = bufio.NewWriter(&payload)
	)

	payloadSize, err := content.WriteTo(writer)
	if err != nil {
		t.Fatalf("should not error, but got: %v", err)
	}
	err = writer.Flush()
	if err != nil {
		t.Fatalf("should not error, but got: %v", err)
	}

	expectedPayloadSize := 80
	if expectedPayloadSize != int(payloadSize) {
		t.Errorf("PayloadSize expected: %v, got: %v", expectedPayloadSize, payloadSize)
	}
	expectedPayload := bytes.Join([][]byte{
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
	}, nil)
	if !reflect.DeepEqual(expectedPayload, payload.Bytes()) {
		t.Errorf("Payload expected: %v, got: %v", expectedPayload, payload.Bytes())
	}
}

func TestDecodeMessageContent(t *testing.T) {
	b := bytes.Join([][]byte{
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
	}, nil)

	content, err := DecodeMessageContent(b)
	if err != nil {
		t.Fatalf("should not error, but got: %v", err)
	}

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

func TestMessageContent_Sanity(t *testing.T) {
	content := MessageContent{}

	content.State.Set("foo1", []byte("bar1"))
	content.State.Set("foo2", []byte("bar2-baz"))
	content.Body = []byte("the quick brown fox jumps over the lazy dog")

	var (
		payload bytes.Buffer
		writer  = bufio.NewWriter(&payload)
	)

	payloadSize, err := content.WriteTo(writer)
	if err != nil {
		t.Fatalf("should not error, but got: %v", err)
	}
	err = writer.Flush()
	if err != nil {
		t.Fatalf("should not error, but got: %v", err)
	}

	{
		expectedPayloadSize := 80
		if expectedPayloadSize != int(payloadSize) {
			t.Errorf("PayloadSize expected: %v, got: %v", expectedPayloadSize, payloadSize)
		}
		expectedPayload := bytes.Join([][]byte{
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
		}, nil)
		if !reflect.DeepEqual(expectedPayload, payload.Bytes()) {
			t.Errorf("Payload expected: %v, got: %v", expectedPayload, payload.Bytes())
		}
	}

	out, err := DecodeMessageContent(payload.Bytes())
	if err != nil {
		t.Fatalf("should not error, but got: %v", err)
	}

	{
		var state map[string][]byte = make(map[string][]byte)
		out.State.Visit(func(name string, value []byte) {
			state[name] = value
		})
		expectedState := content.State.values
		if !reflect.DeepEqual(expectedState, state) {
			t.Errorf("MessageContent.State() expected: %v, got: %v", expectedState, state)
		}
	}
	{
		expectedBody := content.Body
		if !reflect.DeepEqual(expectedBody, out.Body) {
			t.Errorf("MessageContent.Body expected: %v, got: %v", expectedBody, out.Body)
		}
	}
}
