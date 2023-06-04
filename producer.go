package nsq

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"sync"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

type Producer struct {
	pool *ProducerPool

	logger *log.Logger

	wg          sync.WaitGroup
	mutex       sync.Mutex
	disposed    bool
	initialized bool
}

func NewProducer(config *ProducerConfig) (*Producer, error) {
	instance := &Producer{}

	var err error
	err = instance.init(config)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (p *Producer) AllHandles() []*nsq.Producer {
	return p.pool.handles
}

func (p *Producer) Handle() *nsq.Producer {
	return p.pool.handles[p.pool.current]
}

func (p *Producer) WriteContent(topic string, msg *MessageContent, opts ...ProduceMessageContentOption) error {
	if p.disposed {
		return fmt.Errorf("the Producer has been disposed")
	}
	if !p.initialized {
		p.logger.Panic("the Producer haven't be initialized yet")
	}

	// apply options
	for _, opt := range opts {
		err := opt.apply(topic, msg)
		if err != nil {
			return err
		}
	}

	var (
		payload bytes.Buffer
		w       *bufio.Writer = bufio.NewWriter(&payload)
	)
	if _, err := msg.WriteTo(w); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}

	return p.Write(topic, payload.Bytes())
}

func (p *Producer) Write(topic string, body []byte) error {
	if p.disposed {
		return fmt.Errorf("the Producer has been disposed")
	}
	if !p.initialized {
		defaultLogger.Panic("the Producer haven't be initialized yet")
	}

	p.wg.Add(1)
	defer p.wg.Done()

	return p.pool.publish(topic, body)
}

func (p *Producer) DeferredWriteContent(topic string, delay time.Duration, msg *MessageContent, opts ...ProduceMessageContentOption) error {
	if p.disposed {
		return fmt.Errorf("the Producer has been disposed")
	}
	if !p.initialized {
		defaultLogger.Panic("the Producer haven't be initialized yet")
	}

	// apply options
	for _, opt := range opts {
		err := opt.apply(topic, msg)
		if err != nil {
			return err
		}
	}

	var (
		payload bytes.Buffer
		w       *bufio.Writer = bufio.NewWriter(&payload)
	)
	if _, err := msg.WriteTo(w); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}

	return p.DeferredWrite(topic, delay, payload.Bytes())
}

func (p *Producer) DeferredWrite(topic string, delay time.Duration, body []byte) error {
	if p.disposed {
		return fmt.Errorf("the Producer has been disposed")
	}
	if !p.initialized {
		defaultLogger.Panic("the Producer haven't be initialized yet")
	}

	p.wg.Add(1)
	defer p.wg.Done()

	return p.pool.deferredPublish(topic, delay, body)
}

func (p *Producer) Close() {
	if p.disposed {
		return
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.disposed = true

	p.wg.Wait()
	p.pool.dispose()
}

func (p *Producer) init(config *ProducerConfig) error {
	if p.initialized {
		return nil
	}

	if config == nil {
		config = &ProducerConfig{}
	}
	config.init()
	err := config.Validate()
	if err != nil {
		return err
	}

	// config logger
	p.configureLogger(config)

	// config Producer.pool
	{
		var (
			handles []*nsq.Producer
		)
		for _, addr := range config.Address {
			q, err := nsq.NewProducer(addr, config.Config)
			if err != nil {
				return err
			}
			if q != nil {
				q.SetLogger(p.logger, nsq.LogLevelInfo)
				// test the connection
				err = q.Ping()
				if err != nil {
					return fmt.Errorf("cannot establish connection to '%s'; %v", addr, err)
				}
				handles = append(handles, q)
			}
		}
		assert(len(handles) > 0, "assertion failed: Producer must own at least one nsq producer")

		pool := &ProducerPool{
			handles:           handles,
			replicationFactor: config.ReplicationFactor,
		}
		pool.init()

		p.pool = pool

		p.initialized = true
	}

	return nil
}

func (p *Producer) configureLogger(config *ProducerConfig) {
	if config.Logger != nil {
		p.logger = config.Logger
		return
	}
	p.logger = defaultLogger
}
