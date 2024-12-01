package nsq

import (
	"net/http"

	nsq "github.com/nsqio/go-nsq"
)

var _ ConsumerOption = ConsumerOptionFunc(nil)

type ConsumerOptionFunc func(consumer *nsq.Consumer)

// apply implements ConsumerOption.
func (fn ConsumerOptionFunc) apply(consumer *nsq.Consumer) {
	fn(consumer)
}

// //////////////////////////////////////
func WithLoggerLevel(lv nsq.LogLevel) ConsumerOption {
	return ConsumerOptionFunc(func(consumer *nsq.Consumer) {
		if consumer == nil {
			return
		}
		consumer.SetLoggerLevel(lv)
	})
}

// //////////////////////////////////////
func WithLookupdHttpClient(client *http.Client) ConsumerOption {
	return ConsumerOptionFunc(func(consumer *nsq.Consumer) {
		if consumer == nil {
			return
		}
		consumer.SetLookupdHttpClient(client)
	})
}
