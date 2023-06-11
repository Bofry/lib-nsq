package nsq_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	nsq "github.com/Bofry/lib-nsq"
	"github.com/Bofry/lib-nsq/tracing"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slices"
)

var (
	traceIDStr = "4bf92f3577b34da6a3ce929d0e0e4736"
	spanIDStr  = "00f067aa0ba902b7"

	__TEST_NSQD_SERVERS       []string
	__TEST_NSQD_ADDRESS       string
	__TEST_NSQLOOKUPD_ADDRESS string

	__ENV_FILE        = "nsq_test.env"
	__ENV_FILE_SAMPLE = "nsq_test.env.sample"

	__TEST_TRACE_ID = mustTraceIDFromHex(traceIDStr)
	__TEST_SPAN_ID  = mustSpanIDFromHex(spanIDStr)

	__TEST_PROPAGATOR = propagation.TraceContext{}
	__TEST_CONTEXT    = mustSpanContext()
)

func mustTraceIDFromHex(s string) (t trace.TraceID) {
	var err error
	t, err = trace.TraceIDFromHex(s)
	if err != nil {
		panic(err)
	}
	return
}

func mustSpanIDFromHex(s string) (t trace.SpanID) {
	var err error
	t, err = trace.SpanIDFromHex(s)
	if err != nil {
		panic(err)
	}
	return
}

func mustSpanContext() context.Context {
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    __TEST_TRACE_ID,
		SpanID:     __TEST_SPAN_ID,
		TraceFlags: 0,
	})
	return trace.ContextWithSpanContext(context.Background(), sc)
}

func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func TestMain(m *testing.M) {
	_, err := os.Stat(__ENV_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			err = copyFile(__ENV_FILE_SAMPLE, __ENV_FILE)
			if err != nil {
				panic(err)
			}
		}
	}

	{
		f, err := os.Open(__ENV_FILE)
		if err != nil {
			panic(err)
		}
		env, err := godotenv.Parse(f)
		if err != nil {
			panic(err)
		}
		__TEST_NSQD_SERVERS = strings.Split(env["TEST_NSQD_SERVERS"], ",")
		__TEST_NSQD_ADDRESS = env["TEST_NSQD_ADDRESS"]
		__TEST_NSQLOOKUPD_ADDRESS = env["TEST_NSQLOOKUPD_ADDRESS"]
	}
	m.Run()
}

func TestProducer_Write(t *testing.T) {
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
			t.Fatal(err)
		}

		topic := "goTestProducer_Write"
		for _, word := range []string{"Welcome", "to", "the", "Nsq", "Golang", "client", "library"} {
			p.Write(topic, []byte(word))
		}

		p.Close()
	}

	var receivedMessages []string
	defer func() {
		var expectMessages = []string{"Welcome", "to", "the", "Nsq", "Golang", "client", "library"}
		sort.Strings(receivedMessages)
		sort.Strings(expectMessages)

		if !slices.Equal(expectMessages, receivedMessages) {
			t.Errorf("received messages expected: %v, got: %v", expectMessages, receivedMessages)
		}
	}()

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
			NsqAddress:         __TEST_NSQLOOKUPD_ADDRESS,
			Channel:            "gotest",
			HandlerConcurrency: 3,
			Config:             config,
			MessageHandler: nsq.MessageHandleProc(func(message *nsq.Message) error {
				// t.Logf("[%s] %+v\n", message.Topic, string(message.Body))
				receivedMessages = append(receivedMessages, string(message.Body))
				message.Finish()
				return nil
			}),
		}

		ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

		err := c.Subscribe([]string{"goTestProducer_Write"})
		if err != nil {
			t.Fatal(err)
		}

		select {
		case <-ctx.Done():
			c.Close()
			return
		}
	}
}

func TestProducer_WriteContent(t *testing.T) {
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
			t.Fatal(err)
		}

		topic := "goTestProducer_WriteContent"
		for _, word := range []string{"Welcome", "to", "the", "Nsq", "Golang", "client", "library"} {
			msg := nsq.MessageContent{
				Body: []byte(word),
			}
			msg.State.Set("foo", []byte("bar"))

			p.WriteContent(topic, &msg)
		}
		p.Close()
	}

	var receivedMessages []string
	defer func() {
		var expectMessages = []string{"Welcome", "to", "the", "Nsq", "Golang", "client", "library"}
		sort.Strings(receivedMessages)
		sort.Strings(expectMessages)

		if !slices.Equal(expectMessages, receivedMessages) {
			t.Errorf("received messages expected: %v, got: %v", expectMessages, receivedMessages)
		}
	}()

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
			NsqAddress:         __TEST_NSQLOOKUPD_ADDRESS,
			Channel:            "gotest",
			HandlerConcurrency: 3,
			Config:             config,
			MessageHandler: nsq.MessageHandleProc(func(message *nsq.Message) error {
				content := message.Content()
				if content == nil {
					t.Errorf("missing MessageContent")
				}

				if content.State.Len() == 0 {
					t.Errorf("missing MessageState")
				}
				actualStateFooValue := content.State.Value("foo")
				expectStateFooValue := []byte("bar")
				if !reflect.DeepEqual(expectStateFooValue, actualStateFooValue) {
					t.Errorf("MessageState[foo] expected: %v, got: %v", expectStateFooValue, actualStateFooValue)
				}

				// t.Logf("[%s] %+v\n", message.Topic, string(content.Body))
				receivedMessages = append(receivedMessages, string(content.Body))
				message.Finish()
				return nil
			}),
		}

		ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

		err := c.Subscribe([]string{"goTestProducer_WriteContent"})
		if err != nil {
			t.Fatal(err)
		}

		select {
		case <-ctx.Done():
			c.Close()
			return
		}
	}
}

func TestProducer_WriteContent_WithTracePropagation(t *testing.T) {
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
			t.Fatal(err)
		}

		topic := "goTestProducer_WriteContent_WithTracePropagation"
		for _, word := range []string{"Welcome", "to", "the", "Nsq", "Golang", "client", "library"} {

			p.WriteContent(topic,
				&nsq.MessageContent{
					State: nsq.MessageState{},
					Body:  []byte(word),
				},
				nsq.WithTracePropagation(__TEST_CONTEXT, __TEST_PROPAGATOR),
			)
		}
		p.Close()
	}

	var receivedMessages []string
	defer func() {
		var expectMessages = []string{"Welcome", "to", "the", "Nsq", "Golang", "client", "library"}
		sort.Strings(receivedMessages)
		sort.Strings(expectMessages)

		if !slices.Equal(expectMessages, receivedMessages) {
			t.Errorf("received messages expected: %v, got: %v", expectMessages, receivedMessages)
		}
	}()

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
			NsqAddress:         __TEST_NSQLOOKUPD_ADDRESS,
			Channel:            "gotest",
			HandlerConcurrency: 3,
			Config:             config,
			MessageHandler: nsq.MessageHandleProc(func(message *nsq.Message) error {
				content := message.Content()
				if content == nil {
					t.Errorf("missing MessageContent")
				}

				if content.State.Len() == 0 {
					t.Errorf("missing MessageState")
				}
				if len(content.State.Value("traceparent")) == 0 {
					t.Errorf("missing MessageState[traceparent]")
				}

				// test carrier
				carrier := tracing.NewMessageStateCarrier(&content.State)
				if len(carrier.Get("traceparent")) == 0 {
					t.Errorf("missing MessageStateCarrier.Get(traceparent)")
				}
				propagator := __TEST_PROPAGATOR
				ctx := propagator.Extract(context.Background(), carrier)
				spx := trace.SpanContextFromContext(ctx)
				expectedTraceID := __TEST_TRACE_ID
				if !reflect.DeepEqual(expectedTraceID, spx.TraceID()) {
					t.Errorf("received trace id expected: %v, got: %v", expectedTraceID, spx.TraceID())
				}
				expectedSpanID := __TEST_SPAN_ID
				if !reflect.DeepEqual(expectedSpanID, spx.SpanID()) {
					t.Errorf("received span id expected: %v, got: %v", expectedSpanID, spx.SpanID())
				}

				// t.Logf("[%s] %+v\n", message.Topic, string(content.Body))
				receivedMessages = append(receivedMessages, string(content.Body))
				message.Finish()
				return nil
			}),
		}

		ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

		err := c.Subscribe([]string{"goTestProducer_WriteContent_WithTracePropagation"})
		if err != nil {
			t.Fatal(err)
		}

		select {
		case <-ctx.Done():
			c.Close()
			return
		}
	}
}
