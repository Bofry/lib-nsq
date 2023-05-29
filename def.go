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
	logger *log.Logger = log.New(log.Writer(), LOGGER_PREFIX, log.LstdFlags|log.Lmsgprefix)
)

type (
	Config  = nsq.Config
	Message = nsq.Message

	MessageHandleProc func(ctx *ConsumeContext, message *Message) error
)
