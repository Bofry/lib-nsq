package tracing

type MessageStateCarrier struct {
	value MessageState
}

func NewMessageStateCarrier(state MessageState) *MessageStateCarrier {
	return &MessageStateCarrier{
		value: state,
	}
}

// Get returns the value associated with the passed key.
func (hc MessageStateCarrier) Get(key string) string {
	return string(hc.value.Value(key))
}

// Set stores the key-value pair.
func (hc MessageStateCarrier) Set(key string, value string) {
	hc.value.Set(key, []byte(value))
}

// Keys lists the keys stored in this carrier.
func (hc MessageStateCarrier) Keys() []string {
	state := hc.value
	keys := make([]string, 0, state.Len())
	state.Visit(func(key string, value []byte) {
		keys = append(keys, key)
	})
	return keys
}
