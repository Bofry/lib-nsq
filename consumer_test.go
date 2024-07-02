package nsq_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	nsq "github.com/Bofry/lib-nsq"
)

func TestConsumer_Subscribe(t *testing.T) {
	t.Parallel()

	// publish
	{
		p, err := nsq.NewProducer(&nsq.ProducerConfig{
			Address:           __TEST_NSQD_SERVERS,
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

		{
			topic := "gotestConsumer_Subscribe1"
			for _, word := range []string{"Welcome", "to", "the", "Nsq", "Golang", "client", "library"} {
				p.Write(topic, []byte(word))
			}
		}

		{
			topic := "gotestConsumer_Subscribe2"
			for _, word := range []string{"foo", "bar", "baz", "qux", "quux", "corge", "grault", "garply", "waldo", "fred", "plugh", "xyzzy", "thud"} {
				p.Write(topic, []byte(word))
			}
		}
	}

	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	var msgCnt int = 0

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
		NsqAddress:         __TEST_NSQLOOKUPD_ADDRESS,
		Channel:            "gotest",
		HandlerConcurrency: 3,
		Config:             config,
		MessageHandler: nsq.MessageHandleProc(func(message *nsq.Message) error {
			fmt.Printf("[%s] %+v\n", message.Topic, string(message.Body))
			message.Finish()
			msgCnt++
			return nil
		}),
	}

	err := c.Subscribe([]string{"gotestConsumer_Subscribe1", "gotestConsumer_Subscribe2"})
	if err != nil {
		panic(err)
	}

	select {
	case <-ctx.Done():
		c.Close()

		// assert
		{
			var expectedMsgCnt int = 20
			if msgCnt != expectedMsgCnt {
				t.Errorf("expect %d messages, but got %d messages", expectedMsgCnt, msgCnt)
			}
		}
		return
	}
}

func TestConsumer_Pause(t *testing.T) {
	t.Parallel()

	// publish
	{
		p, err := nsq.NewProducer(&nsq.ProducerConfig{
			Address:           __TEST_NSQD_SERVERS,
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

		{
			topic := "gotestConsumer_Pause1"
			for _, word := range []string{"Welcome", "to", "the", "Nsq", "Golang", "client", "library"} {
				p.DeferredWrite(topic, 2*time.Second, []byte(word))
			}
		}

		{
			topic := "gotestConsumer_Pause2"
			for _, word := range []string{"foo", "bar", "baz", "qux", "quux", "corge", "grault", "garply", "waldo", "fred", "plugh", "xyzzy", "thud"} {
				p.DeferredWrite(topic, 2*time.Second, []byte(word))
			}
		}
	}

	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	var msgCnt int = 0

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
		NsqAddress:         __TEST_NSQLOOKUPD_ADDRESS,
		Channel:            "gotest",
		HandlerConcurrency: 3,
		Config:             config,
		MessageHandler: nsq.MessageHandleProc(func(message *nsq.Message) error {
			fmt.Printf("[%s] %+v\n", message.Topic, string(message.Body))
			message.Finish()
			msgCnt++
			return nil
		}),
	}

	err := c.Subscribe([]string{"gotestConsumer_Pause1", "gotestConsumer_Pause2"})
	if err != nil {
		panic(err)
	}
	c.Pause("gotestConsumer_Pause1")

	select {
	case <-ctx.Done():
		c.Close()

		// assert
		{
			var expectedMsgCnt int = 13
			if msgCnt != expectedMsgCnt {
				t.Errorf("expect %d messages, but got %d messages", expectedMsgCnt, msgCnt)
			}
		}
		return
	}
}

func TestConsumer_Resume(t *testing.T) {
	t.Parallel()

	// publish
	{
		p, err := nsq.NewProducer(&nsq.ProducerConfig{
			Address:           __TEST_NSQD_SERVERS,
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

		{
			topic := "gotestConsumer_Resume1"
			for _, word := range []string{"Welcome", "to", "the", "Nsq", "Golang", "client", "library"} {
				p.DeferredWrite(topic, 2*time.Second, []byte(word))
			}
		}

		{
			topic := "gotestConsumer_Resume2"
			for _, word := range []string{"foo", "bar", "baz", "qux", "quux", "corge", "grault", "garply", "waldo", "fred", "plugh", "xyzzy", "thud"} {
				p.DeferredWrite(topic, 2*time.Second, []byte(word))
			}
		}
	}

	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	var msgCnt int = 0

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
		NsqAddress:         __TEST_NSQLOOKUPD_ADDRESS,
		Channel:            "gotest",
		HandlerConcurrency: 3,
		Config:             config,
		MessageHandler: nsq.MessageHandleProc(func(message *nsq.Message) error {
			fmt.Printf("[%s] %+v\n", message.Topic, string(message.Body))
			message.Finish()
			msgCnt++
			return nil
		}),
	}

	err := c.Subscribe([]string{"gotestConsumer_Resume1", "gotestConsumer_Resume2"})
	if err != nil {
		panic(err)
	}
	c.Pause("gotestConsumer_Resume1")
	c.Resume("gotestConsumer_Resume1")

	select {
	case <-ctx.Done():
		c.Close()

		// assert
		{
			var expectedMsgCnt int = 20
			if msgCnt != expectedMsgCnt {
				t.Errorf("expect %d messages, but got %d messages", expectedMsgCnt, msgCnt)
			}
		}
		return
	}
}
