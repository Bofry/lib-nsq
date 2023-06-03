package nsq_test

import (
	"context"
	"fmt"
	"time"

	nsq "github.com/Bofry/lib-nsq"
)

func Example() {
	// publish
	{
		p, err := nsq.NewProducer(&nsq.ProducerConf{
			Address:           []string{"127.0.0.1:4150"},
			Config:            nsq.NewConfig(),
			ReplicationFactor: 1,
		})
		if err != nil {
			if p != nil {
				p.Close()
			}
			panic(err)
		}
		defer p.Close()

		topic := "myTopic"
		for _, word := range []string{"Welcome", "to", "the", "Nsq", "Golang", "client", "library"} {
			p.Write(topic, []byte(word))
		}
	}

	// subscribe
	{
		// the config only for test use !!
		config := nsq.NewConfig()
		{
			config.LookupdPollInterval = time.Second * 3
			config.DefaultRequeueDelay = 0
			config.MaxBackoffDuration = time.Millisecond * 50
			config.LowRdyIdleTimeout = time.Second * 1
			config.RDYRedistributeInterval = time.Millisecond * 20
		}

		c := &nsq.Consumer{
			NsqAddress:         "nsqlookupd://127.0.0.1:4160",
			Channel:            "gotest",
			HandlerConcurrency: 3,
			Config:             config,
			MessageHandler: nsq.MessageHandleProc(func(ctx *nsq.ConsumeContext, message *nsq.Message) error {
				fmt.Printf("[%s] %+v\n", ctx.Topic, string(message.Body))
				message.Finish()
				return nil
			}),
			UnhandledMessageHandler: nil,
		}

		ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

		err := c.Subscribe([]string{"myTopic"})
		if err != nil {
			panic(err)
		}

		select {
		case <-ctx.Done():
			c.Close()
			return
		}
	}
	// Output:
	// [myTopic] Welcome
	// [myTopic] to
	// [myTopic] the
	// [myTopic] Nsq
	// [myTopic] Golang
	// [myTopic] client
	// [myTopic] library
}
