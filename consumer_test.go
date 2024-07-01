package nsq_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	nsq "github.com/Bofry/lib-nsq"
	gonsq "github.com/nsqio/go-nsq"
)

var _ gonsq.ConnDelegate = new(DummyConnDelegate)

type DummyConnDelegate struct{}

func (d *DummyConnDelegate) OnBackoff(*gonsq.Conn)                         {}
func (d *DummyConnDelegate) OnClose(*gonsq.Conn)                           {}
func (d *DummyConnDelegate) OnContinue(*gonsq.Conn)                        {}
func (d *DummyConnDelegate) OnError(*gonsq.Conn, []byte)                   {}
func (d *DummyConnDelegate) OnHeartbeat(*gonsq.Conn)                       {}
func (d *DummyConnDelegate) OnIOError(*gonsq.Conn, error)                  {}
func (d *DummyConnDelegate) OnMessage(*gonsq.Conn, *gonsq.Message)         {}
func (d *DummyConnDelegate) OnMessageFinished(*gonsq.Conn, *gonsq.Message) {}
func (d *DummyConnDelegate) OnMessageRequeued(*gonsq.Conn, *gonsq.Message) {}
func (d *DummyConnDelegate) OnResponse(*gonsq.Conn, []byte)                {}
func (d *DummyConnDelegate) OnResume(*gonsq.Conn)                          {}

func TestConsumer_Subscribe(t *testing.T) {
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
			topic := "gotestTopic1"
			for _, word := range []string{"Welcome", "to", "the", "Nsq", "Golang", "client", "library"} {
				p.Write(topic, []byte(word))
			}
		}

		{
			topic := "gotestTopic2"
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

	err := c.Subscribe([]string{"gotestTopic1", "gotestTopic2"})
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
			topic := "gotestTopic3"
			for _, word := range []string{"Welcome", "to", "the", "Nsq", "Golang", "client", "library"} {
				p.DeferredWrite(topic, 2*time.Second, []byte(word))
			}
		}

		{
			topic := "gotestTopic4"
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

	err := c.Subscribe([]string{"gotestTopic3", "gotestTopic4"})
	if err != nil {
		panic(err)
	}
	c.Pause("gotestTopic3")

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

func TestConsumer_PauseAndResume(t *testing.T) {
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
			topic := "gotestTopic5"
			for _, word := range []string{"Welcome", "to", "the", "Nsq", "Golang", "client", "library"} {
				p.DeferredWrite(topic, 2*time.Second, []byte(word))
			}
		}

		{
			topic := "gotestTopic6"
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

	err := c.Subscribe([]string{"gotestTopic5", "gotestTopic6"})
	if err != nil {
		panic(err)
	}
	c.Pause("gotestTopic5")
	c.Resume("gotestTopic5")

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
