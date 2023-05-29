package nsq

import (
	"reflect"
	"testing"
)

func TestMessageState_Init(t *testing.T) {
	state := MessageState{}

	{
		size := state.Len()
		var expectedSize int = 0
		if expectedSize != size {
			t.Errorf("MessageState.Len() expected: %v, got: %v", expectedSize, size)
		}

		bytes := state.byteSize()
		var expectedByteSize int = 0
		if expectedByteSize != bytes {
			t.Errorf("MessageState.byteSize() expected: %v, got: %v", expectedByteSize, bytes)
		}

		ok := state.Has("foo")
		var expectedOk bool = false
		if expectedOk != ok {
			t.Errorf("MessageState.Has(foo) expected: %v, got: %v", expectedOk, ok)
		}
	}
}

func TestMessageState_Set(t *testing.T) {
	state := MessageState{}

	// add foo, <empty>
	{
		old, err := state.Set("foo", []byte(""))
		var expectedErr error = nil
		if expectedErr != err {
			t.Errorf("MessageState.Set().Err expected: %v, got: %v", expectedErr, err)
		}
		var expectedOld []byte
		if !reflect.DeepEqual(expectedOld, old) {
			t.Errorf("MessageState.Set().Old expected: %v, got: %v", expectedOld, old)
		}

		size := state.Len()
		var expectedSize int = 0
		if expectedSize != size {
			t.Errorf("MessageState.Len() expected: %v, got: %v", expectedSize, size)
		}

		bytes := state.byteSize()
		var expectedByteSize int = 0
		if expectedByteSize != bytes {
			t.Errorf("MessageState.byteSize() expected: %v, got: %v", expectedByteSize, bytes)
		}

		ok := state.Has("foo")
		var expectedOk bool = false
		if expectedOk != ok {
			t.Errorf("MessageState.Has(foo) expected: %v, got: %v", expectedOk, ok)
		}
	}

	// add foo, bar
	{
		old, err := state.Set("foo", []byte("bar"))
		var expectedErr error = nil
		if expectedErr != err {
			t.Errorf("MessageState.Set() expected: %v, got: %v", expectedErr, err)
		}
		var expectedOld []byte
		if !reflect.DeepEqual(expectedOld, old) {
			t.Errorf("MessageState.Set().Old expected: %v, got: %v", expectedOld, old)
		}

		size := state.Len()
		var expectedSize int = 1
		if expectedSize != size {
			t.Errorf("MessageState.Len() expected: %v, got: %v", expectedSize, size)
		}

		bytes := state.byteSize()
		var expectedByteSize int = 6
		if expectedByteSize != bytes {
			t.Errorf("MessageState.byteSize() expected: %v, got: %v", expectedByteSize, bytes)
		}

		ok := state.Has("foo")
		var expectedOk bool = true
		if expectedOk != ok {
			t.Errorf("MessageState.Has(foo) expected: %v, got: %v", expectedOk, ok)
		}
	}

	// add foo2, bar2
	{
		old, err := state.Set("foo2", []byte("bar2"))
		var expectedErr error = nil
		if expectedErr != err {
			t.Errorf("MessageState.Set() expected: %v, got: %v", expectedErr, err)
		}
		var expectedOld []byte
		if !reflect.DeepEqual(expectedOld, old) {
			t.Errorf("MessageState.Set().Old expected: %v, got: %v", expectedOld, old)
		}

		size := state.Len()
		var expectedSize int = 2
		if expectedSize != size {
			t.Errorf("MessageState.Len() expected: %v, got: %v", expectedSize, size)
		}

		bytes := state.byteSize()
		var expectedByteSize int = 14
		if expectedByteSize != bytes {
			t.Errorf("MessageState.byteSize() expected: %v, got: %v", expectedByteSize, bytes)
		}

		ok := state.Has("foo2")
		var expectedOk bool = true
		if expectedOk != ok {
			t.Errorf("MessageState.Has(foo) expected: %v, got: %v", expectedOk, ok)
		}
	}

	// add baz, <empty>
	{
		old, err := state.Set("baz", []byte(""))
		var expectedErr error = nil
		if expectedErr != err {
			t.Errorf("MessageState.Set() expected: %v, got: %v", expectedErr, err)
		}
		var expectedOld []byte
		if !reflect.DeepEqual(expectedOld, old) {
			t.Errorf("MessageState.Set().Old expected: %v, got: %v", expectedOld, old)
		}

		size := state.Len()
		var expectedSize int = 2
		if expectedSize != size {
			t.Errorf("MessageState.Len() expected: %v, got: %v", expectedSize, size)
		}

		bytes := state.byteSize()
		var expectedByteSize int = 14
		if expectedByteSize != bytes {
			t.Errorf("MessageState.byteSize() expected: %v, got: %v", expectedByteSize, bytes)
		}

		ok := state.Has("baz")
		var expectedOk bool = false
		if expectedOk != ok {
			t.Errorf("MessageState.Has(foo) expected: %v, got: %v", expectedOk, ok)
		}
	}

	// replace foo2, <empty>. It same as Del().
	{
		old, err := state.Set("foo2", []byte(""))
		var expectedErr error = nil
		if expectedErr != err {
			t.Errorf("MessageState.Set() expected: %v, got: %v", expectedErr, err)
		}
		var expectedOld []byte = []byte("bar2")
		if !reflect.DeepEqual(expectedOld, old) {
			t.Errorf("MessageState.Set().Old expected: %v, got: %v", expectedOld, old)
		}

		size := state.Len()
		var expectedSize int = 1
		if expectedSize != size {
			t.Errorf("MessageState.Len() expected: %v, got: %v", expectedSize, size)
		}

		bytes := state.byteSize()
		var expectedByteSize int = 6
		if expectedByteSize != bytes {
			t.Errorf("MessageState.byteSize() expected: %v, got: %v", expectedByteSize, bytes)
		}

		ok := state.Has("foo2")
		var expectedOk bool = false
		if expectedOk != ok {
			t.Errorf("MessageState.Has(foo) expected: %v, got: %v", expectedOk, ok)
		}
	}

	// re-add foo2, bar2
	{
		old, err := state.Set("foo2", []byte("bar2"))
		var expectedErr error = nil
		if expectedErr != err {
			t.Errorf("MessageState.Set() expected: %v, got: %v", expectedErr, err)
		}
		var expectedOld []byte
		if !reflect.DeepEqual(expectedOld, old) {
			t.Errorf("MessageState.Set().Old expected: %v, got: %v", expectedOld, old)
		}

		size := state.Len()
		var expectedSize int = 2
		if expectedSize != size {
			t.Errorf("MessageState.Len() expected: %v, got: %v", expectedSize, size)
		}

		bytes := state.byteSize()
		var expectedByteSize int = 14
		if expectedByteSize != bytes {
			t.Errorf("MessageState.byteSize() expected: %v, got: %v", expectedByteSize, bytes)
		}

		ok := state.Has("foo2")
		var expectedOk bool = true
		if expectedOk != ok {
			t.Errorf("MessageState.Has(foo) expected: %v, got: %v", expectedOk, ok)
		}
	}
}

func TestMessageState_Del(t *testing.T) {
	state := MessageState{}

	// delete missing key foo
	{
		size := state.Len()
		var expectedSize int = 0
		if expectedSize != size {
			t.Errorf("MessageState.Len() expected: %v, got: %v", expectedSize, size)
		}

		bytes := state.byteSize()
		var expectedByteSize int = 0
		if expectedByteSize != bytes {
			t.Errorf("MessageState.byteSize() expected: %v, got: %v", expectedByteSize, bytes)
		}

		ok := state.Has("foo")
		var expectedOk bool = false
		if expectedOk != ok {
			t.Errorf("MessageState.Has(foo) expected: %v, got: %v", expectedOk, ok)
		}

		old := state.Del("foo")
		var expectedOld []byte
		if !reflect.DeepEqual(expectedOld, old) {
			t.Errorf("MessageState.Del(foo).Old expected: %v, got: %v", expectedOld, old)
		}
	}

	// delete key foo
	{
		state.Set("foo", []byte("bar"))
		size := state.Len()
		var expectedSize int = 1
		if expectedSize != size {
			t.Errorf("MessageState.Len() expected: %v, got: %v", expectedSize, size)
		}

		old := state.Del("foo")
		var expectedOld []byte = []byte("bar")
		if !reflect.DeepEqual(expectedOld, old) {
			t.Errorf("MessageState.Del().Old expected: %v, got: %v", expectedOld, old)
		}

		size = state.Len()
		expectedSize = 0
		if expectedSize != size {
			t.Errorf("MessageState.Len() expected: %v, got: %v", expectedSize, size)
		}

		bytes := state.byteSize()
		var expectedByteSize int = 0
		if expectedByteSize != bytes {
			t.Errorf("MessageState.byteSize() expected: %v, got: %v", expectedByteSize, bytes)
		}

		ok := state.Has("foo")
		var expectedOk bool = false
		if expectedOk != ok {
			t.Errorf("MessageState.Has(foo) expected: %v, got: %v", expectedOk, ok)
		}
	}
}

func TestMessageState_Value(t *testing.T) {
	state := MessageState{}

	// get missing key foo
	{
		size := state.Len()
		var expectedSize int = 0
		if expectedSize != size {
			t.Errorf("MessageState.Len() expected: %v, got: %v", expectedSize, size)
		}

		bytes := state.byteSize()
		var expectedByteSize int = 0
		if expectedByteSize != bytes {
			t.Errorf("MessageState.byteSize() expected: %v, got: %v", expectedByteSize, bytes)
		}

		ok := state.Has("foo")
		var expectedOk bool = false
		if expectedOk != ok {
			t.Errorf("MessageState.Has(foo) expected: %v, got: %v", expectedOk, ok)
		}

		value := state.Value("foo")
		var expectedValue []byte
		if !reflect.DeepEqual(expectedValue, value) {
			t.Errorf("MessageState.Value(foo) expected: %v, got: %v", expectedValue, value)
		}
	}

	// get foo
	{
		state.Set("foo", []byte("bar"))
		size := state.Len()
		var expectedSize int = 1
		if expectedSize != size {
			t.Errorf("MessageState.Len() expected: %v, got: %v", expectedSize, size)
		}

		value := state.Value("foo")
		var expectedValue []byte = []byte("bar")
		if !reflect.DeepEqual(expectedValue, value) {
			t.Errorf("MessageState.Value(foo) expected: %v, got: %v", expectedValue, value)
		}
	}

	// get foo after value has been updated
	{
		state.Set("foo", []byte("bar2"))
		size := state.Len()
		var expectedSize int = 1
		if expectedSize != size {
			t.Errorf("MessageState.Len() expected: %v, got: %v", expectedSize, size)
		}

		value := state.Value("foo")
		var expectedValue []byte = []byte("bar2")
		if !reflect.DeepEqual(expectedValue, value) {
			t.Errorf("MessageState.Value(foo) expected: %v, got: %v", expectedValue, value)
		}
	}
}

func TestMessageState_Visit(t *testing.T) {
	state := MessageState{}

	// visit empty MessageState
	{
		size := state.Len()
		var expectedSize int = 0
		if expectedSize != size {
			t.Errorf("MessageState.Len() expected: %v, got: %v", expectedSize, size)
		}

		bytes := state.byteSize()
		var expectedByteSize int = 0
		if expectedByteSize != bytes {
			t.Errorf("MessageState.byteSize() expected: %v, got: %v", expectedByteSize, bytes)
		}

		ok := state.Has("foo")
		var expectedOk bool = false
		if expectedOk != ok {
			t.Errorf("MessageState.Has(foo) expected: %v, got: %v", expectedOk, ok)
		}

		var values map[string][]byte = make(map[string][]byte)
		state.Visit(func(name string, value []byte) {
			values[name] = value
		})
		var expectedValues map[string][]byte = map[string][]byte{}
		if !reflect.DeepEqual(expectedValues, values) {
			t.Errorf("MessageState.Visit() expected: %v, got: %v", expectedValues, values)
		}
	}

	{
		state.Set("foo", []byte("bar"))
		size := state.Len()
		var expectedSize int = 1
		if expectedSize != size {
			t.Errorf("MessageState.Len() expected: %v, got: %v", expectedSize, size)
		}

		var values map[string][]byte = make(map[string][]byte)
		state.Visit(func(name string, value []byte) {
			values[name] = value
		})
		var expectedValues map[string][]byte = map[string][]byte{
			"foo": []byte("bar"),
		}
		if !reflect.DeepEqual(expectedValues, values) {
			t.Errorf("MessageState.Visit() expected: %v, got: %v", expectedValues, values)
		}
	}

	{
		state.Set("foo", []byte("bar2"))
		size := state.Len()
		var expectedSize int = 1
		if expectedSize != size {
			t.Errorf("MessageState.Len() expected: %v, got: %v", expectedSize, size)
		}

		var values map[string][]byte = make(map[string][]byte)
		state.Visit(func(name string, value []byte) {
			values[name] = value
		})
		var expectedValues map[string][]byte = map[string][]byte{
			"foo": []byte("bar2"),
		}
		if !reflect.DeepEqual(expectedValues, values) {
			t.Errorf("MessageState.Visit() expected: %v, got: %v", expectedValues, values)
		}
	}

	{
		state.Set("foo", []byte(""))
		size := state.Len()
		var expectedSize int = 0
		if expectedSize != size {
			t.Errorf("MessageState.Len() expected: %v, got: %v", expectedSize, size)
		}

		var values map[string][]byte = make(map[string][]byte)
		state.Visit(func(name string, value []byte) {
			values[name] = value
		})
		var expectedValues map[string][]byte = map[string][]byte{}
		if !reflect.DeepEqual(expectedValues, values) {
			t.Errorf("MessageState.Visit() expected: %v, got: %v", expectedValues, values)
		}
	}
}
