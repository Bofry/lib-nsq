package nsq

import (
	"fmt"
	"reflect"
	"sync"
)

const (
	MESSAGE_STATE_NAME_MAX_LENGTH = 255
	MESSAGE_STATE_VALUE_MAX_SIZE  = 0x0fff
)

type MessageState struct {
	values map[string][]byte
	size   int

	valuesOnce sync.Once
}

func (s *MessageState) Len() int {
	if s.values == nil {
		return 0
	}
	return len(s.values)
}

func (s *MessageState) Has(name string) bool {
	if s.values == nil {
		return false
	}

	_, ok := s.values[name]
	return ok
}

func (s *MessageState) Del(name string) []byte {
	if s.values == nil {
		return nil
	}

	if v, ok := s.values[name]; ok {
		// deduct size from total size
		s.size -= (len(v) + len(name))

		delete(s.values, name)
		return v
	}
	return nil
}

func (s *MessageState) SetString(name, value string) (old []byte, err error) {
	return s.Set(name, []byte(value))
}

func (s *MessageState) Set(name string, value []byte) (old []byte, err error) {
	if err := s.validateName(name); err != nil {
		return nil, err
	}
	if len(value) > MESSAGE_STATE_VALUE_MAX_SIZE {
		return nil, fmt.Errorf("specified value exceed (max size: %d)", MESSAGE_STATE_VALUE_MAX_SIZE)
	}

	if s.values == nil {
		if len(value) == 0 {
			return nil, nil
		}

		s.valuesOnce.Do(func() {
			s.values = make(map[string][]byte)
		})
	}
	// name existed?
	if v, ok := s.values[name]; ok {
		old = v

		// ignored when new value equals old value
		if reflect.DeepEqual(v, value) {
			return old, nil
		}

		// deduct size from total size
		s.size -= (len(v) + len(name))

		// delete the key and exit when the value is empty
		if len(value) == 0 {
			delete(s.values, name)
			return old, nil
		}
	}

	if len(value) > 0 {
		s.values[name] = value
		// increase added size
		s.size += (len(value) + len(name))
	}
	return old, nil
}

func (s *MessageState) Value(name string) []byte {
	if s.values == nil {
		return nil
	}

	return s.values[name]
}

func (s *MessageState) Visit(visit func(name string, value []byte)) {
	if s.values != nil {
		for k, v := range s.values {
			visit(k, v)
		}
	}
}

func (s *MessageState) validateName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("specified name is invalid")
	}
	if len(name) > MESSAGE_STATE_NAME_MAX_LENGTH {
		return fmt.Errorf("specified name is too long (max lenght: %d)", MESSAGE_STATE_NAME_MAX_LENGTH)
	}
	for i := 0; i < len(name); i++ {
		ch := name[i]
		if !(s.isValidNameChar(ch)) {
			return fmt.Errorf("specified name contains invalid '%c' at %d", ch, i)
		}
	}
	return nil
}

func (s *MessageState) byteSize() int {
	return s.size
}

func (s *MessageState) isValidNameChar(ch byte) bool {
	if ch == '_' || ch == '-' ||
		(ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') {
		return true
	}
	return false
}
