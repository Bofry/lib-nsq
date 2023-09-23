package nsq

import (
	"log"

	"github.com/nsqio/go-nsq"
)

const (
	SERVICE_NSQD       = "nsqd"
	SERVICE_NSQLOOKUPD = "nsqlookupd"

	LOGGER_PREFIX string = "[lib-nsq] "
)

var (
	defaultLogger *log.Logger = log.New(log.Writer(), LOGGER_PREFIX, log.LstdFlags|log.Lmsgprefix)
)

type (
	Config          = nsq.Config
	MessageDelegate = nsq.MessageDelegate

	MessageHandleProc func(message *Message) error

	ProduceMessageContentOption interface {
		apply(msg *MessageContent) error
	}
)

func DefaultLogger() *log.Logger {
	return defaultLogger
}
