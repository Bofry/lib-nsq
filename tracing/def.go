package tracing

import (
	"context"
)

type (
	MessageState interface {
		Len() int
		Has(name string) bool
		Del(name string) []byte
		Set(name string, value []byte) (old []byte, err error)
		Value(name string) []byte
		Visit(visit func(name string, value []byte))
	}

	Context interface {
		context.Context

		SetValue(key, value interface{})
	}
)
